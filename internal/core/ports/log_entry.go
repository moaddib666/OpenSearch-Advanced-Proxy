package ports

import "time"

type LogEntry interface {
	Raw() string
	Map() map[string]interface{}
	Timestamp() time.Time
	Load(string string) error
	LoadBytes([]byte) error
	ID() string
}

type EntryConstructor func() LogEntry
