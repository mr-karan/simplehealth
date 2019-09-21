package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/mr-karan/simplehealth"
	"github.com/prometheus/common/log"
)

var (
	kvs = map[string]func() bool{"redis": dumbHealthCheck, "db": dumbHealthCheck}
)

func dumbHealthCheck() bool {
	// flaky service. Sometimes returns false, sometimes true.
	if rand.Intn(2) == 0 {
		return false
	}
	return true
}

func main() {
	// initialize manager
	var m, err = simplehealth.NewManager(kvs, simplehealth.Options{ExpositionFormat: "prometheus"})
	if err != nil {
		panic(fmt.Errorf("error registering manager %s", err))
	}
	// Default index handler.
	handleIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to simplehealth. Visit /metrics to scrape prometheus metrics."))
	})
	// Initialize router and define all endpoints.
	router := http.NewServeMux()
	router.Handle("/", handleIndex)
	// Expose the registered metrics at `/metrics` path.
	router.Handle("/metrics", m.Collect())
	server := &http.Server{
		Addr:         "localhost:8888",
		Handler:      router,
		ReadTimeout:  6000 * time.Millisecond,
		WriteTimeout: 6000 * time.Millisecond,
	}
	// Start the server. Blocks the main thread.
	log.Infof("starting server listening on %s", "127.0.0.1:8888")
	if err := server.ListenAndServe(); err != nil {
		log.Errorf("error starting server: %s", err)
	}
}
