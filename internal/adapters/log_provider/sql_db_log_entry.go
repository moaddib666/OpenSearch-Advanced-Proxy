package log_provider

import (
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/ports"
	"time"
)

type SQLDBLogEntry struct {
	fieldsMap map[string]interface{}
}

func (S *SQLDBLogEntry) RawString() string {
	return ""
}

func (S *SQLDBLogEntry) RawBytes() []byte {
	return []byte{}
}

func (S *SQLDBLogEntry) Map() map[string]interface{} {
	return S.fieldsMap
}

func (S *SQLDBLogEntry) Timestamp() time.Time {
	return time.Time{}
}

func (S *SQLDBLogEntry) LoadString(string string) error {
	return nil
}

func (S *SQLDBLogEntry) LoadBytes(bytes []byte) error {
	return nil
}

func (S *SQLDBLogEntry) LoadMap(raw map[string]interface{}) error {
	S.fieldsMap = raw
	return nil
}

func (S *SQLDBLogEntry) ID() string {
	return ""
}

// NewSQLDBLogEntry creates a new SQLDBLogEntry
func NewSQLDBLogEntry() *SQLDBLogEntry {
	return &SQLDBLogEntry{}
}

// SqlLogEntryConstructor is a constructor for SQLDBLogEntry
func SqlLogEntryConstructor() ports.LogEntry {
	return NewSQLDBLogEntry()
}
