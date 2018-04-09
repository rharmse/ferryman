package ferryman_test

import (
	"net/http"
	"testing"

	"github.com/rharmse/ferryman/lib"
)

//Tests Router Initialisation
func TestNew(t *testing.T) {
	router := ferryman.New(&ferryman.Pool{})
	if router == nil {
		t.Errorf("No route store created")
	}
}

// Tests Getting a Default Route and if a http 500 is a valid response
func TestGetRoute(t *testing.T) {
	routePath := "/you/have/no/routes/here/haha"

	defaultPool := &ferryman.Pool{}
	router := ferryman.New(defaultPool)
	route := router.GetRoute(routePath, true)

	//Route must be present
	if route == nil {
		t.Errorf("No default route created")
	}

	//Any http response must be valid
	if !route.ValidRouteResponseStatus(http.StatusInternalServerError) {
		t.Errorf("Invalid http.StatusCode for route")
	}
}

// Tests if http.StatusInternalServerError is a valid response for a route
// when it is expecting http.StatusOK, no handlers specified
func TestInvValidRouteResponseStatus(t *testing.T) {
	routePath := "/you/have/no/routes/here/haha"
	defaultPool := &ferryman.Pool{}

	router := ferryman.New(defaultPool)
	gotRoute := router.GetRoute(routePath, false)

	if gotRoute == nil {
		route, err := router.AddRoute(routePath, map[int]int{http.StatusOK: http.StatusOK}, ferryman.Default, nil, nil, nil)

		if err == nil && route != nil {
			route = router.GetRoute(routePath, false)
			if !(route.ValidRouteResponseStatus(http.StatusInternalServerError) == false && route.ValidRouteResponseStatus(http.StatusOK)) {
				t.Errorf("Route must fail given http response http.StatusInternalServerError and route expects http.StatusOK.")
			}
		} else {
			t.Errorf("No route %v or err %v", route, err)
		}
	} else {
		t.Errorf("Route should be nil, no default asked for.")
	}
}

// Test Adding Routes with different permutations of handlers
func TestAddRoute(t *testing.T) {
	routePathDef := "/route/with/default/handler"
	routePathDefErr := "/route/with/default/and/error/handler"
	routePathAll := "/route/path/with/all/hanlers"

	defaultPool := &ferryman.Pool{}
	router := ferryman.New(defaultPool)

	var defaultHandler = func(rw http.ResponseWriter, r *http.Request, route *ferryman.Route) (err error) {
		return nil
	}

	var fallbackHandler = func(rw http.ResponseWriter, r *http.Request) (err error) {
		return nil
	}

	var errHandler = func(rw http.ResponseWriter, r *http.Request, route *ferryman.Route) {

	}

	route1, err1 := router.AddRoute(routePathDef, ferryman.StatusALLMap, ferryman.Default, defaultHandler, nil, fallbackHandler)
	route2, err2 := router.AddRoute(routePathDefErr, ferryman.StatusALLMap, ferryman.Default, defaultHandler, errHandler, nil)
	route3, err3 := router.AddRoute(routePathAll, ferryman.StatusALLMap, ferryman.Default, nil, errHandler, nil)
	route4, err4 := router.AddRoute(routePathAll, ferryman.StatusALLMap, ferryman.Default, nil, nil, nil)

	if err1 != nil && route1 == nil {
		t.Errorf("Route creation with specified default handler failed. %v", err1)
	}

	if err2 != nil && route2 == nil {
		t.Errorf("Route creation with specified default and error handler failed. %v", err2)
	}

	if err3 != nil && route3 == nil {
		t.Errorf("Route creation with specified error handler failed. %v", err3)
	}

	if err4 != nil && route4 == nil {
		t.Errorf("Route creation with ferryman default and error handler failed. %v", err3)
	}
}
