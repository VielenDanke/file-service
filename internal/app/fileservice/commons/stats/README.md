# statscheck

A library for easy publishing of k8s liveness and readiness probes, AppVersion and BuildDate, Prometheus metrics.

### usage

`NewStatsCheckMux` returns ServeMux, which can be merged with another http Mux/Handler/Server
`RunStatsCheckServer` runs StatsCheck Server in a goroutine