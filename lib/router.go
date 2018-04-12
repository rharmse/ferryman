// This is where the routing functions are mainly handled including any
// http Handlers.
package ferryman

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kellydunn/go-art"
)

// Primitive to indicate the Rule/Route types
type routeRule int

// This constant is used with the default handler in order to
// ensure any http code passed from the upstream is considered valid.
const (
	StatusALL int = 999
)

var StatusALLMap = map[int]int{StatusALL: StatusALL}

// Different routeRule types
const (
	_              = iota
	Drop routeRule = iota
	TempRedirect
	PermRedirect
	Forward
	Default
)

// Defines the base http Handlers types
type routeHandlerFunc func(rw http.ResponseWriter, r *http.Request, route *Route) (err error)
type fallbackHandlerFunc func(rw http.ResponseWriter, r *http.Request) (err error)
type errHandlerFunc func(rw http.ResponseWriter, r *http.Request, route *Route)
type dropHandlerFunc func(rw http.ResponseWriter, r *http.Request)

// This is the representation of a route, routes are bound to pools
// A routeHandlerFunc is 1..1, fallbackHandlerFunc and errHandlerFunc is 0..1.
// targetURI is the relative uriPath that the request is forarded to.
type Route struct {
	routeType            routeRule
	apply                routeHandlerFunc
	fallback             fallbackHandlerFunc
	err                  errHandlerFunc
	targetURI            string
	validRespStatusCodes map[int]int
	pool                 *Pool
}

// Houses the routes and specifies default pool
type Router struct {
	routes      *art.ArtTree
	defaultPool *Pool
}

var (
	DefaultHandler = func(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
		return defaultRouteHandler(rw, r, route)
	}

	RedirectingHandler = func(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
		return redirectingRouteHandler(rw, r, route)
	}

	ErrorHandler = func(rw http.ResponseWriter, r *http.Request, route *Route) {
		errorRouteHandler(rw, r, route)
	}

	DropHandler = func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusGone)
		rw.Write([]byte(http.StatusText(http.StatusGone)))
		defer r.Body.Close()
		r.Close = true
	}
)

// Build a new router load rules should be assigned here from config
func New(defaultPool *Pool) *Router {
	return &Router{routes: art.NewArtTree(), defaultPool: defaultPool}
}

// Load the routes into the router
func (router *Router) LoadRoutes(conf []*RuleConfig, t routeRule) {
	//implement
}

// Register a route with the the router
func (router *Router) AddRoute(routePath string, allowedRespStati map[int]int, routeType routeRule, defaultH routeHandlerFunc, errorH errHandlerFunc, fallbackH fallbackHandlerFunc) (route *Route, err error) {
	switch {
	case defaultH == nil && errorH == nil:
		route = &Route{
			apply:    DefaultHandler,
			err:      ErrorHandler,
			fallback: fallbackH}
	case defaultH != nil && errorH == nil:
		route = &Route{
			apply:    defaultH,
			err:      ErrorHandler,
			fallback: fallbackH}
	case defaultH != nil && errorH != nil:
		route = &Route{
			apply:    defaultH,
			err:      errorH,
			fallback: fallbackH}
	case defaultH == nil && errorH != nil:
		route = &Route{
			apply:    DefaultHandler,
			err:      errorH,
			fallback: fallbackH}
	}

	route.routeType = routeType
	route.pool = router.defaultPool
	route.targetURI = routePath
	route.validRespStatusCodes = allowedRespStati

	router.routes.Insert([]byte(routePath), route)
	return route, nil
}

// Implement the ServeHTTP function Sig, making the router a HTTP Handler
func (ro *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route := ro.GetRoute(r.RequestURI, true)
	switch route.routeType {
	case Default:
		//fmt.Printf("\nServing Default Route %v", r)
		route.apply(rw, r, route)
	case Drop:
		DropHandler(rw, r)
	default:
		err := errors.New("Undefined route type!")
		panic(err)
	}
}

// Fetch route for the current uri requested path
func (router *Router) GetRoute(uriPath string, createDefault bool) (route *Route) {
	potential := router.routes.Search([]byte(uriPath))
	if potential == nil {
		if createDefault {
			route := &Route{
				apply:                DefaultHandler,
				routeType:            Default,
				targetURI:            uriPath,
				validRespStatusCodes: map[int]int{StatusALL: StatusALL},
				pool:                 router.defaultPool,
			}
			router.routes.Insert([]byte(uriPath), route)
			return route
		} else {
			return nil
		}
	} else {
		return potential.(*Route)
	}
}

// Check if the response http.StausCode is valid for this route
// Default is StatusALL unless specific codes required
func (route *Route) ValidRouteResponseStatus(statusCode int) bool {
	code := route.validRespStatusCodes[statusCode]
	if route.validRespStatusCodes[StatusALL] == StatusALL || statusCode == code {
		return true
	}
	return false
}

// Build the proxy request to send to the upstream server
func buildProxyRequest(r *http.Request, baseURI, targetURI string) (pr *http.Request, err error) {
	//fmt.Printf("\nURI:%v", baseURI+targetURI)
	pr = r.WithContext(r.Context())
	url := new(url.URL)
	//fmt.Printf("\nURI:%v, ERR:%v", url, err)
	pr.URL, err = url.Parse(baseURI + targetURI)
	pr.RequestURI = ""
	pr.Close = false
	if err == nil {
		return pr, nil
	} else {
		return nil, err
	}
}

// Clone response headers from upstream response to client response
// TODO Filter allowed headers out, rewrite some where required
func copyResponseHeaders(rw http.ResponseWriter, res *http.Response) {
	for k, v := range res.Header {
		for _, vv := range v {
			rw.Header().Add(k, vv)
		}
	}
	rw.Header().Set("Served-By", "Ferryman")
}

// Clone the client request headers to the upsteam proxy request
// TODO Filter allowed headers in, rewrite some where required
func copyProxyRequestHeaders(req *http.Request, pReq *http.Request) {
	for k, v := range req.Header {
		for _, vv := range v {
			pReq.Header.Add(k, vv)
		}
	}
}

// Default ruleRoute HTTP handler.
func defaultRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
	st := time.Now()
	pool := route.pool
	var sessionId string

	member, session := pool.getUpstreamTarget(r)
	if session != nil {
		member = session.pinned
		sessionId = session.id
	} else {
		sessionId = "none"
	}
	pr, err := buildProxyRequest(r, member.nodeURI, route.targetURI)

	if err == nil {
		member.requestCnt++
		resp, err := member.httpClient.Do(pr)
		//fmt.Printf("Served with session:%v\n", session)

		if err == nil {
			if route.ValidRouteResponseStatus(resp.StatusCode) {
				if pool.sticky {
					if session == nil {
						pool.createSession(resp, member)
					} else {
						pool.maintainSession(session)
					}
				}
				defer r.Body.Close()
				defer resp.Body.Close()
				copyResponseHeaders(rw, resp)
				rw.Header().Set("X-Ferried-In", time.Since(st).String())
				//Triggers TTFB response write
				rw.WriteHeader(resp.StatusCode)

				//TODO potentially use bytecnt over time to adaptively size buffers according
				//to response size, can be tracked from route
				byteCnt, err := io.CopyBuffer(rw, resp.Body, make([]byte, 1024*64))

				if err == nil {
					fmt.Printf("\nRead :%v kb from [host:session][%v:%v] in %v", byteCnt/1024, member.hostname, sessionId, time.Since(st))
				}
				return err
			} else {
				return fallbackOrErr(rw, r, route)
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

// Fallback handler if there is a error from the upstream
func fallbackOrErr(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
	if route.fallback != nil {
		err := route.fallback(rw, r)
		if err != nil {
			route.err(rw, r, route)
			return err
		}
	} else {
		route.err(rw, r, route)
	}
	return nil
}

// This is the default routeRule handler for TempRedirect and PermRedirect
func redirectingRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
	return nil
}

// This is the default errorHandler implementation
func errorRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) {
	defer r.Body.Close()
	io.Copy(rw, r.Body)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Header().Set("Served-By", "Ferryman")
	return
}
