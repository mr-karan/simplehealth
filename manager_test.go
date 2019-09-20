package healthexporter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}

func initDummyManager() Manager {
	kvs := map[string]func() bool{"dummy": func() bool {
		return true
	}}
	return NewManager("test_app", kvs, false)
}

func TestMetricsHandler(t *testing.T) {
	// initialize manager
	var m = initDummyManager()
	// initialize request object
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// call the collect method
	handler := http.HandlerFunc(m.Collect())
	handler.ServeHTTP(rr, req)

	// Assert whether status code is 200 OK.
	assert.Equal(t, rr.Code, http.StatusOK, "Status code should be 200")
	// Assert whether response body is what we expect.
	expected := `test_app{service="dummy"} 1`
	assert.Contains(t, rr.Body.String(), expected, "Response body differs")
}
