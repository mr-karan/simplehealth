package healthexporter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	kv = map[string]func() bool{"dummy": func() bool {
		return true
	}}
)

func TestManagerDefaultConfiguration(t *testing.T) {
	// initialize manager with default options
	j, err := NewManager(kv, Options{})
	assert.Nil(t, err)
	assert.Implements(t, (*Manager)(nil), j)
}

func TestManagerValidConfiguration(t *testing.T) {
	// initialize manager with valid options
	j, err := NewManager(kv, Options{ExpositionFormat: "json", ExposeDefaultMetrics: true})
	assert.Nil(t, err)
	assert.Implements(t, (*Manager)(nil), j)

	p, err := NewManager(kv, Options{Namespace: "hello", ExpositionFormat: "prometheus"})
	assert.Nil(t, err)
	assert.Implements(t, (*Manager)(nil), p)
}

func TestManagerInvalidConfiguration(t *testing.T) {
	// initialize manager with invalid options
	kvs := map[string]func() bool{"dummy": func() bool {
		return true
	}}
	var m, err = NewManager(kvs, Options{ExpositionFormat: "abcdlol"})
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestHandlerPrometheusFormat(t *testing.T) {
	// initialize manager
	p, err := NewManager(kv, Options{Namespace: "hello", ExpositionFormat: "prometheus"})
	assert.Nil(t, err)
	// initialize request object
	req, err := http.NewRequest("GET", "/health-check", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	// call the collect method
	handler := http.HandlerFunc(p.Collect())
	handler.ServeHTTP(rr, req)
	// Assert whether status code is 200 OK.
	assert.Equal(t, rr.Code, http.StatusOK, "Status code differs")
	// Assert whether response body is what we expect.
	expected := `hello{service="dummy"} 1`
	assert.Contains(t, rr.Body.String(), expected, "Response body differs")
}

func TestHandlerJSONFormat(t *testing.T) {
	// initialize manager
	p, err := NewManager(kv, Options{ExpositionFormat: "json"})
	assert.Nil(t, err)
	// initialize request object
	req, err := http.NewRequest("GET", "/health-check", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	// call the collect method
	handler := http.HandlerFunc(p.Collect())
	handler.ServeHTTP(rr, req)
	// Assert whether status code is 200 OK.
	assert.Equal(t, rr.Code, http.StatusOK, "Status code differs")
	// Assert whether response body is what we expect.
	expected := `{"dummy": "healthy"}`
	assert.JSONEqf(t, expected, rr.Body.String(), "Response body differs")
}
