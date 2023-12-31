package log_provider

import (
	"OpenSearchAdvancedProxy/internal/core/ports"
	"encoding/json"
	"time"
)

type JsonLogEntry struct {
	raw     string
	tsField string
}

func (j *JsonLogEntry) Raw() string {
	return j.raw
}

func (j *JsonLogEntry) Map() map[string]interface{} { // TODO cache this
	logMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(j.raw), &logMap) // FIXME bytes
	return logMap
}

func (j *JsonLogEntry) Timestamp() time.Time {
	logMap := j.Map()
	ts, ok := logMap[j.tsField]
	if !ok {
		return time.Time{}
	}
	tsStr, ok := ts.(string)
	if !ok {
		return time.Time{}
	}
	tsTime, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		return time.Time{}
	}
	return tsTime
}

func (j *JsonLogEntry) Load(raw string) error {
	j.raw = raw
	return nil
}

func (j *JsonLogEntry) LoadBytes(raw []byte) error {
	j.raw = string(raw)
	return nil
}

// JsonLogEntryConstructor creates a new JsonLogEntry struct
func JsonLogEntryConstructor() ports.LogEntry {
	return &JsonLogEntry{
		tsField: "datetime",
	}
}
