// Package micrometric is a tiny metrics exporter.
// As of now it only exports to Prometheus.
package micrometric

// Metric is a structure to contain a metric name/value pair as well
// as a list of associated labels and values
type Metric struct {
	Labels map[string]string
	Name   string
	Value  float64
}

// Exporter is the interface all micrometrics exporters must implement.
type Exporter interface {
	// Export sets up the list of current metrics to be exported.
	Export([]Metric) error

	// Serve sets up the exporter server.
	Serve() error
}
