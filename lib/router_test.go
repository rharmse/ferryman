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

//Tests Getting a Default Route (no route found in ARTree)
func TestGetRoute(t *testing.T) {
	routePath := "/you/have/no/routes/here/haha"
	defaultPool := &ferryman.Pool{}

	//
	router := ferryman.New(defaultPool)
	if router.routes == nil {
		t.Errorf("No route store created")
	}
	route := router.getRoute(routePath)

	if route == nil &&
		//has StatusALL as valid response type
		route.validRouteResponseStatus(ferryman.StatusALL) &&
		//has a defaultRouteHandler
		route.apply != nil &&
		//Has a default pool, same as constructed on router
		(route.pool != nil && route.pool == defaultPool) &&
		//Has a default type route
		route.rType == ferryman.Default {

		t.Errorf("No default route created")
	}
}

// Tests Default Rout
func TestValidRouteResponseStatus(t *testing.T) {
	route := &ferryman.Route{
		validRespStatusCodes: map[int]int{StatusALL: StatusALL},
	}

	if !route.validRouteResponseStatus(http.StatusInternalServerError) {
		t.Errorf("Invalid http.StatusCode for route")
	}
}

// Tests if http.StatusInternalServerError is a valid response for a route
func TestInvValidRouteResponseStatus(t *testing.T) {
	route := &ferryman.Route{
		validRespStatusCodes: map[int]int{http.StatusOK: http.StatusOK},
	}

	if route.validRouteResponseStatus(http.StatusInternalServerError) {
		t.Errorf("Invalid http.StatusCode for route")
	}
}
