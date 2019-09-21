package simplehealth

import (
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

const (
	// Default namespace for collecting metrics.
	defaultNamespace = "app"
	// Default format for collecting metrics.
	defaultFormat = "prometheus"
)

// Target represents options for registering a new manager.
type Target struct {
	metrics map[string]func() bool // metric names with their callbacks to extract the value
	opts    *Options               // additional options for exporter
}

// Options represents configuration option for exporter.
type Options struct {
	Namespace            string // global namespace for the app
	ExposeDefaultMetrics bool   // whether to expose go_* and process_* metrics
	ExpositionFormat     string // prometheus or json
}

// Manager represents the set of methods for collecting and exposing metrics.
type Manager interface {
	Collect() http.HandlerFunc
}

// NewManager instantiates an object of Manager.
func NewManager(ms map[string]func() bool, opts Options) (Manager, error) {
	t := &Target{
		metrics: ms,
		opts:    &opts,
	}
	// Set default namespace
	if t.opts.Namespace == "" {
		t.opts.Namespace = defaultNamespace
	}
	// Set default format
	if t.opts.ExpositionFormat == "" {
		t.opts.ExpositionFormat = defaultFormat
	}
	// Validate format value
	err := validateFormat(t.opts.ExpositionFormat)
	if err != nil {
		return nil, err
	}
	return t, err
}

// Collect creates metrics on the fly and fetches value of the metric by executing the callbacks.s
func (target *Target) Collect() http.HandlerFunc {
	if target.opts.ExpositionFormat == "json" {
		data, err := collectJSONFormat(target.metrics)
		if err != nil {
			fmt.Errorf("error collecting data in json: %s", err)
			return sendResponse([]byte("Error collecting targets and exposing in json format. Please check logs"), http.StatusInternalServerError)
		}
		return sendResponse(data, http.StatusOK)
	}
	collectPrometheusFormat(target.metrics, target.opts.Namespace)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, target.opts.ExposeDefaultMetrics)
	})
}
