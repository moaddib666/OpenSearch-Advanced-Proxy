package monitoring

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// PrometheusMetricsExporter implements the MetricsExporter interface using Prometheus.
type PrometheusMetricsExporter struct{}

// Bind starts an HTTP server on the specified address and exposes Prometheus metrics.
func (e *PrometheusMetricsExporter) Bind(address string) {
	mux := http.NewServeMux()
	// Setup HTTP handler for Prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())

	// Start HTTP server
	go func() {
		log.Infof("Starting metrics server on %s", address)
		if err := http.ListenAndServe(address, mux); err != nil {
			log.Errorf("Error starting metrics server: %v", err)
		}
	}()
}

// NewMetrics returns a new PrometheusMetricsExporter instance.
func NewMetrics() *PrometheusMetricsExporter {
	return &PrometheusMetricsExporter{}
}
