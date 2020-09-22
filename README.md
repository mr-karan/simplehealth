<a href="https://zerodha.tech"><img src="https://zerodha.tech/static/images/github-badge.svg" align="right" /></a>

# simplehealth
> Expose health check endpoints as Prometheus/JSON format in seconds 💫

[![GoDoc](https://godoc.org/github.com/mr-karan/simplehealth?status.svg)](http://godoc.org/github.com/mr-karan/simplehealth)

## Overview
`simplehealth` makes it ridiculously easy to expose health checks in your Go service. It is a tiny abstraction over [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics) which is a lightweight alternative to the official Prometheus [client](https://github.com/prometheus/client_golang) library.

## Features
- Accepts a map of service name and a callback func to determine if the service is up or not.
- Exposition format is configurable, can be in `JSON` or [Prometheus](https://prometheus.io/docs/instrumenting/writing_exporters/) format.
- Extract metrics as a plug and play on any HTTP handler which implements `HandlerFunc`.

## Installation
Install `simplehealth` using

```
go get -u github.com/mr-karan/simplehealth
```

## Usage

### Examples

Check the [_examples](/_examples) directory for a complete working example.

### Registering Metrics

Metrics need to be registered with a `manager` using the `NewManager` method.
`NewManager` accepts a map of service name with it's callback function for executing the health check.

For example:

```go
// {"api": true} indicates api service is up and running...
callbacks := map[string]func() bool{"api": func() bool{
	return true
}}
// will construct a metric like `namespace{service="api"} 1`
manager := simplehealth.NewManager(callbacks, simplehealth.Options{})
```

### Exposing Metrics

`manager` comes with a `Collect` method which returns a `http.HandlerFunc`.

```go
router := http.NewServeMux()
// Expose the registered metrics at `/metrics` path.
router.Handle("/metrics", m.Collect())
```

## Exposition Format

You can configure to export metrics either in Prometheus (_default_) or JSON format.

- Prometheus example

```bash
curl localhost:8888/metrics

app{service="db"} 1
app{service="redis"} 1
```

- JSON example

```bash
curl localhost:8888/metrics

{ 
   "db":"healthy",
   "redis":"healthy"
}
```

## Configuring Manager

While creating a new manager, you can pass in additional options with `simplehealth.Options{}`:

| Option       | Type                           | Description                                                                                                                                               |
| ------------ | -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------|
| Namespace        | `string`                  | Global namespace for each metric exposed as Prometheus format. _Optional. (Default: "app")_        |                     
| ExposeDefaultMetrics | `bool`                  |  Whether to expose default metrics like `go_*` and `process_*` _Optional. (Default: false)_  |
| ExpositionFormat | `string`                  | Format to expose the metrics. Can be one of `prometheus` or `json` _Optional. (Default: "prometheus")_ |
