package ferryman

import (
	"net/http"
	"testing"
)

// Tests if http.StatusInternalServerError is a valid response for a route
func TestValidRouteResponseStatus(t *testing.T) {
	route := &Route{
		validRespStatusCodes: map[int]int{StatusALL: StatusALL},
	}

	if route.validRouteResponseStatus(http.StatusInternalServerError) {
		t.Errorf("Invalid http.StatusCode for route")
	}
}
