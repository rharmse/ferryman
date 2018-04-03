//Router to implement adaptive radix tree https://pdfs.semanticscholar.org/6abf/5107efc723c655956f027b4a67565b048799.pdf
//Reference http://daslab.seas.harvard.edu/classes/cs265/files/presentations/CS265_presentation_Sinyagin.pdf
package ferryman

import (
	"fmt"
	"io"
	"net/http"
	"github.com/kellydunn/go-art"
)

//Primitive to indicate the rule types
type ruleType int

const (
	StatusALL int = -1
)

//Route Types
const (
	_                  = iota
	RouteDrop ruleType = iota
	RouteTempRedirect
	RoutePermRedirect
	RouteForward
	Default
)

type routeHandlerFunc func(rw http.ResponseWriter, r *http.Request, route *Route) (err error)
type fallbackHandlerFunc func(rw http.ResponseWriter, r *http.Request) (err error)
type errHandlerFunc func(rw http.ResponseWriter, r *http.Request, route *Route)
type dropHandlerFunc func(rw http.ResponseWriter, r *http.Request)

type Route struct {
	rType           ruleType
	apply           routeHandlerFunc
	fallback        fallbackHandlerFunc
	err             errHandlerFunc
	targetURI       string
	expResStatCodes map[int]int
	pool            *Pool
}

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

//Build a new router load rules should be assigned here from
//config
func New(defaultPool *Pool) *Router {
	return &Router{routes: art.NewArtTree(), defaultPool: defaultPool}
}

//Need to discern types of routes
func (router *Router) LoadRoutes(conf []*RuleConfig, t ruleType) {
	//implement
}

//Implement the ServeHTTP function Sig, making the router a HTTP Handler
func (ro *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	/*ctx := r.Context()
	if cn, ok := rw.(http.CloseNotifier); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		notifyChan := cn.CloseNotify()
		go func() {
			select {
			case <-notifyChan:
				cancel()
			case <-ctx.Done():
			}
		}()
	}*/

	route := ro.getRoute(r.RequestURI)
	switch route.rType {
	case Default:
		fmt.Printf("\nServing Default Route %v", r)
		route.apply(rw, r, route)
	default:
		fmt.Printf("\nServing %v", r)
		route.apply(rw, r, route)
	case RouteDrop:
		dropHandler(rw, r)
	}
}

func (router *Router) getRoute(uriPath string) (route *Route) {
	potential := router.routes.Search([]byte(uriPath))
	if potential == nil {
		route = &Route{
			apply:           defaultRouteHandler,
			rType:           Default,
			targetURI:       uriPath,
			expResStatCodes: map[int]int{StatusALL: StatusALL},
			pool:            router.defaultPool,
		}
	} else {
		route = potential.(*Route)
	}
	return route
}

func buildProxyRequest(r *http.Request, baseURI, targetURI string) (pr *http.Request, err error) {
	fmt.Printf("\nURI:%v", baseURI+targetURI)
	pr = r.WithContext(r.Context())
	url, err := r.URL.Parse(baseURI + targetURI)
	fmt.Printf("\nURI:%v, ERR:%v", url, err)
	if err == nil {
		pr.URL = url
		return pr, nil
	} else {
		return nil, err
	}
}

func statusValid(statusCode int, validResponses map[int]int) bool {
	code := validResponses[statusCode]
	fmt.Printf("\nResponse Code:%v:", code)
	if code == StatusALL || code > 0 {
		return true
	} else {
		return validResponses[statusCode] > 0
	}
}

func fallbackOrErr(rw http.ResponseWriter, r *http.Request, route *Route) {
	if route.fallback != nil {
		if route.fallback(rw, r) != nil {
			route.err(rw, r, route)
		}
	} else {
		route.err(rw, r, route)
	}
}

func copyResponseHeaders(rw http.ResponseWriter, res *http.Response) {
	for k, v := range res.Header {
		for _, vv := range v {
			rw.Header().Add(k, vv)
		}
	}
}

func copyProxyRequestHeaders(req *http.Request, pReq *http.Request) {
	for k, v := range req.Header {
		for _, vv := range v {
			pReq.Header.Add(k, vv)
		}
	}
}

//Represents the majority of traffic being handled by the proxy
func defaultRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
	for _, member := range route.pool.members {
		req, err := buildProxyRequest(r, member.nodeURI, route.targetURI)
		fmt.Printf("\nDoing %v", req.URL.String())
		resp, err := member.httpClient.Do(req)
		if err == nil {
			if statusValid(resp.StatusCode, route.expResStatCodes) {
				copyResponseHeaders(rw, resp)
				io.Copy(rw, resp.Body)
				defer resp.Body.Close()
			} else {
				fallbackOrErr(rw, r, route)
			}
		} else {
			fmt.Printf("\nErr:%v", err)
			return err
		}
	}
	return nil
}

func redirectingRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) (err error) {
	return nil
}

func errorRouteHandler(rw http.ResponseWriter, r *http.Request, route *Route) {

}
