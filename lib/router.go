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
	rType                routeRule
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
	defaultHandler = func(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
		return defaultRouteHandler(rw, r, route)
	}

	redirectingHandler = func(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
		return redirectingRouteHandler(rw, r, route)
	}

	errorHandler = func(rw http.ResponseWriter, r *http.Request, route *Route) {
		errorRouteHandler(rw, r, route)
	}

	dropHandler = func(rw http.ResponseWriter, r *http.Request) {
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

// Implement the ServeHTTP function Sig, making the router a HTTP Handler
func (ro *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	route := ro.getRoute(r.RequestURI)
	switch route.rType {
	case Default:
		//fmt.Printf("\nServing Default Route %v", r)
		route.apply(rw, r, route)
	case Drop:
		dropHandler(rw, r)
	default:
		err := errors.New("Undefined route type!")
		panic(err)
	}
}

// Fetch route for the current uri requested path
func (router *Router) getRoute(uriPath string) (route *Route) {
	potential := router.routes.Search([]byte(uriPath))
	if potential == nil {
		route = &Route{
			apply:                defaultRouteHandler,
			rType:                Default,
			targetURI:            uriPath,
			validRespStatusCodes: map[int]int{StatusALL: StatusALL},
			pool:                 router.defaultPool,
		}
	} else {
		route = potential.(*Route)
	}
	return route
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

// Check if the response http.StausCode is valid for this route
// Default is StatusALL unless specific codes required
func (route *Route) validRouteResponseStatus(statusCode int) bool {
	code := route.validRespStatusCodes[statusCode]
	if route.validRespStatusCodes[StatusALL] == StatusALL || code > 0 {
		return true
	} else {
		return statusCode > 0
	}
}

// Fallback handler if there is a error from the upstream
func fallbackOrErr(rw http.ResponseWriter, r *http.Request, route *Route) {
	if route.fallback != nil {
		if route.fallback(rw, r) != nil {
			route.err(rw, r, route)
		}
	} else {
		route.err(rw, r, route)
	}
}

// Clone response headers from upstream response to client response
// TODO Filter allowed headers out
func copyResponseHeaders(rw http.ResponseWriter, res *http.Response) {
	for k, v := range res.Header {
		for _, vv := range v {
			rw.Header().Add(k, vv)
		}
	}
}

// Clone the client request headers to the upsteam proxy request
// TODO Filter allowed headers out
func copyProxyRequestHeaders(req *http.Request, pReq *http.Request) {
	for k, v := range req.Header {
		for _, vv := range v {
			pReq.Header.Add(k, vv)
		}
	}
}

// Default routeRule handler.
func defaultRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
	st := time.Now()
	member := route.pool.getLeastBusy()
	pr, err := buildProxyRequest(r, member.nodeURI, route.targetURI)
	//fmt.Printf("\nDoing %v", pr.URL.String())
	member.requestCnt++
	resp, err := member.httpClient.Do(pr)
	//fmt.Printf("\nDid %v, Err:%v", pr.URL.String(), err)
	if err == nil {
		if route.validRouteResponseStatus(resp.StatusCode) {
			defer r.Body.Close()
			defer resp.Body.Close()
			copyResponseHeaders(rw, resp)
			rw.Header().Set("Served-By", "Ferryman")
			rw.WriteHeader(resp.StatusCode)
			byteCnt, err := io.CopyBuffer(rw, resp.Body, make([]byte, 1024*64))
			if err != nil {
				fmt.Printf("\nError reading response:%v", err)
			} else {
				fmt.Printf("\nRead :%v kb, from %v in %v", byteCnt/1024, pr.URL.String(), time.Since(st))
			}

		} else {
			fallbackOrErr(rw, r, route)
		}
	} else {
		//fmt.Printf("\nErr:%v", err)
		return err
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
