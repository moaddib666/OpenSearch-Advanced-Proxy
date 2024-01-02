package ports

type MetricsExporter interface {
	Bind(address string)
}
