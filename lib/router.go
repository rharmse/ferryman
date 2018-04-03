//Router to implement adaptive radix tree https://pdfs.semanticscholar.org/6abf/5107efc723c655956f027b4a67565b048799.pdf
//Reference http://daslab.seas.harvard.edu/classes/cs265/files/presentations/CS265_presentation_Sinyagin.pdf
package ferryman

import (
	"net/http"
    "github.com/kellydunn/go-art"
)

type Route struct {
    defaultHandler      http.HandlerFunc
    errorHandler        http.HandlerFunc
    pool                *Pool
}

var (
    defaultHandler = func (rw http.ResponseWriter, r *http.Request) {
        defaultRouteHandler(rw, r)
    }
    
    redirectingHandler = func (rw http.ResponseWriter, r *http.Request) {
        redirectingRouteHandler(rw, r)
    }
    
    errorHandler = func (rw http.ResponseWriter, r *http.Request) {
        errorRouteHandler(rw, r)
    }
    
    dropHandler = func (rw http.ResponseWriter, r *http.Request) {
        dropRouteHandler(rw, r)
    }
)

type Router struct {
    routes *art.ArtTree
}

//Build a new router load rules should be assigned here from
//config
func New() *Router {
    return &Router{routes: art.NewArtTree()}
}

func (ro *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    route := ro.getRoute(r.RequestURI)
    route.defaultHandler(rw,r)
}

func (router *Router) getRoute(uriPath string) (route *Route) {
    route = router.routes.Search([]byte(uriPath)).(*Route)
    if route == nil {
        route = &Route{defaultHandler:defaultRouteHandler} 
    }
    return route
}  

func defaultRouteHandler(rw http.ResponseWriter, r *http.Request) {
    
}
    
func redirectingRouteHandler(rw http.ResponseWriter, r *http.Request) {
    
}
    
func errorRouteHandler(rw http.ResponseWriter, r *http.Request) {
    
}
    
func dropRouteHandler(rw http.ResponseWriter, r *http.Request) {
    
}