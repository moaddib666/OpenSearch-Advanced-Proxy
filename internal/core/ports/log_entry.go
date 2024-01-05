package ports

import "time"

type LogEntry interface {
	Raw() string
	Map() map[string]interface{}
	Timestamp() time.Time
	LoadString(string string) error
	LoadBytes([]byte) error
	LoadMap(map[string]interface{}) error
	ID() string
}

type EntryConstructor func() LogEntry
