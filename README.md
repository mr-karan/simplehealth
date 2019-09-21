# healthexporter
> Expose health check endpoint in Prometheus format in seconds 💫

## Overview
`healthexporter` makes is ridiculously easy to add health checks in your Go service. It is a tiny abstraction on [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics) which is a lightweight alternative to the official Prometheus [client](https://github.com/prometheus/client_golang) library.

## Features
- Accepts a map of service name and a callback func to determine if the service is up or not.
- Exposition format is configurable, can be in JSON or Prometheus
- Plug and play in any `http` router

## Installation
Install `healthexporter` using

```
go get -u github.com/mr-karan/healthexporter
```

## Usage

### Examples
Check the [_examples](/_examples) directory for complete examples.

### Registering Metrics

Metrics need to be registered with a `manager` using the `NewManager` method.
`NewManager` accepts a map of metric label value and callback func which return `bool` to indicate whether the service is up or not.

```go
// {"api": true} indicates api service is up and running...
callbacks := map[string]func() bool{"api": func() bool{
	return true
}}
// will construct a metric like `namespace{service="api"} 1`
manager := healthexporter.NewManager("namespace", callbacks)
```

### Exposing Metrics

`manager` comes with a `Collect` method which returns a `http.HandlerFunc` which can be used to expose the metrics on an HTTP endpoint

```
router := http.NewServeMux()
// Expose the registered metrics at `/metrics` path.
router.Handle("/metrics", m.Collect())
```

