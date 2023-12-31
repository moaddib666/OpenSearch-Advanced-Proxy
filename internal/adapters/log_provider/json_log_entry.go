package log_provider

import (
	"OpenSearchAdvancedProxy/internal/core/ports"
	"encoding/json"
	log "github.com/sirupsen/logrus"
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
		log.Debugf("Timestamp field `%s` not found", j.tsField)
		return time.Time{}
	}
	tsStr, ok := ts.(string)
	if !ok {
		log.Debugf("Timestamp field `%s` is not a string", j.tsField)
		return time.Time{}
	}
	tsTime, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		log.Debugf("Timestamp field `%s` is not a valid RFC3339 string", j.tsField)
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
