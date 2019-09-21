package simplehealth

import (
	"fmt"
	"net/http"
)

func boolToFloat(val bool) float64 {
	if val {
		return 1
	}
	return 0
}

func boolToString(val bool) string {
	if val {
		return "healthy"
	}
	return "unhealthy"
}

func sendResponse(msg []byte, status int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(msg)
	})
}

func validateFormat(format string) error {
	switch format {
	case "prometheus":
		return nil
	case "json":
		return nil
	default:
		return fmt.Errorf("expostion format %s is invalid. choose one of prometheus or json", format)
	}
}
