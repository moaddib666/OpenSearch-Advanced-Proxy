package ports

type SearchDataProvider interface {
	Scan() bool
	Text() string
	Err() error
	LogEntry() LogEntry
}
