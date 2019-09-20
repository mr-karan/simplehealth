package healthexporter

import (
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

// Options represents options for registering a new manager.
type Options struct {
	namespace            string                 // global namespace for the app
	metrics              map[string]func() bool // metric names with their callbacks to extract the value
	exposeDefaultMetrics bool                   // whether to expose go_* and process_* metrics
}

// Manager represents the set of methods for collecting and exposing metrics.
type Manager interface {
	Collect() http.HandlerFunc
}

// NewManager instantiates an object of Manager.
func NewManager(ns string, ms map[string]func() bool, expose bool) Manager {
	return &Options{
		namespace:            ns,
		metrics:              ms,
		exposeDefaultMetrics: expose,
	}
}

// Collect creates metrics on the fly and fetches value of the metric by executing the callbacks.s
func (opt *Options) Collect() http.HandlerFunc {
	// Register gauge with two labels.
	for k, v := range opt.metrics {
		// Yes, the label is hardcoded. This library is mostly suited for health-checks, where this label suits most
		// of the use cases. I wanted to abstract away any kind of Prometheus label definitions from the manager and chose to
		// keep the implementation really slimmed. If you don't agree with this you can always use the original library
		// which gives more features and extensibility.
		metrics.GetOrCreateGauge(fmt.Sprintf(`%s{service="%s"}`, opt.namespace, k), func() float64 {
			return boolToFloat(v())
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, opt.exposeDefaultMetrics)
	})

}

func boolToFloat(val bool) float64 {
	if val {
		return 1
	}
	return 0
}
