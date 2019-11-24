package solax_prometheus

import (
  "strconv"
  "go.uber.org/zap"
  "github.com/prometheus/client_golang/prometheus"
)

type Server struct {
  Version         float64
  ScrapeErrCount  int
  Registry        *prometheus.Registry
  Metrics         Metrics
  Endpoint        string
  Log             *zap.Logger
}

func (s *Server) Update( data *APIStructuredData ){
  s.registerIfNotAlready()
  s.updateMetrics(data)
}

func (s *Server) Clear() {
  s.unregisterIfAlready()
}

func Initialise() ( *Server ){
  s := &Server{
    Version: 1.1,
    Registry: prometheus.NewRegistry(),
  }
  s.initMetrics()
  s.registerCoreMetrics()
  s.registerIfNotAlready()
  return s
}

func FloatToString(input_num float64) string {
  return strconv.FormatFloat(input_num, 'f', 2, 64)
}
