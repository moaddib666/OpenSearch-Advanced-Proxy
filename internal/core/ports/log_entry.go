package ports

import "time"

type LogEntry interface {
	RawString() string
	RawBytes() []byte
	Map() map[string]interface{}
	Timestamp() time.Time
	LoadString(string string) error
	LoadBytes([]byte) error
	LoadMap(map[string]interface{}) error
	ID() string
}

type EntryConstructor func() LogEntry
