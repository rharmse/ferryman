package ferryman

import (
	"net/http"
	"testing"
)

func TestStatusValid(t *testing.T) {
	testMap := map[int]int{StatusALL: StatusALL,
		http.StatusOK: http.StatusOK}
	isValid := statusValid(-1, testMap)
	if !isValid {
		t.Errorf("Expected %v, but got %v", true, isValid)
	}
}
