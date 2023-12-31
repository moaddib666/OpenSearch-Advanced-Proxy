package ports

type SearchDataProvider interface {
	BeginScan()
	Scan() bool
	Text() string
	Err() error
	LogEntry() LogEntry
	EndScan()
}
