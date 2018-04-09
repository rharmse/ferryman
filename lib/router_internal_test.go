package ferryman

import (
	"net/http"
	"testing"
)

//Tests Router Initialisation
func TestNew(t *testing.T) {
	router := New(&Pool{})
	if router.routes == nil {
		t.Errorf("No route store created")
	}
}

//Tests Getting a Default Route (no route found in ARTree)
func TestGetRoute(t *testing.T) {
	routePath := "/you/have/no/routes/here/haha"
	defaultPool := &Pool{}

	//
	router := New(defaultPool)
	if router.routes == nil {
		t.Errorf("No route store created")
	}
	route := router.GetRoute(routePath, true)

	if route == nil &&
		//has StatusALL as valid response type
		route.ValidRouteResponseStatus(StatusALL) &&
		//has a defaultRouteHandler
		route.apply != nil &&
		//Has a default pool, same as constructed on router
		(route.pool != nil && route.pool == defaultPool) &&
		//Has a default type route
		route.routeType == Default {

		t.Errorf("No default route created")
	}
}

// Tests Default Rout
func TestValidRouteResponseStatus(t *testing.T) {
	route := &Route{
		validRespStatusCodes: map[int]int{StatusALL: StatusALL},
	}

	if !route.ValidRouteResponseStatus(http.StatusInternalServerError) {
		t.Errorf("Invalid http.StatusCode for route")
	}
}

// Tests if http.StatusInternalServerError is a valid response for a route
func TestInvValidRouteResponseStatus(t *testing.T) {
	route := &Route{
		validRespStatusCodes: map[int]int{http.StatusOK: http.StatusOK},
	}

	if route.ValidRouteResponseStatus(http.StatusInternalServerError) {
		t.Errorf("Invalid http.StatusCode for route")
	}
}
