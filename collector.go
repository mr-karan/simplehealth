package simplehealth

import (
	"encoding/json"
	"fmt"

	"github.com/VictoriaMetrics/metrics"
)

func collectPrometheusFormat(m map[string]func() bool, ns string) {
	for k, v := range m {
		// Yes, the label is hardcoded. This library is mostly suited for health-checks, where this label suits most
		// of the use cases. I wanted to abstract away any kind of Prometheus label definitions from the manager and chose to
		// keep the implementation really slimmed. If you don't agree with this you can always use the original library
		// which gives more features and extensibility.
		metrics.GetOrCreateGauge(fmt.Sprintf(`%s{service="%s"}`, ns, k), func() float64 {
			return boolToFloat(v())
		})
	}
}

func collectJSONFormat(m map[string]func() bool) ([]byte, error) {
	var data = make(map[string]string)
	for k, v := range m {
		data[k] = boolToString(v())
	}
	return json.Marshal(data)
}
