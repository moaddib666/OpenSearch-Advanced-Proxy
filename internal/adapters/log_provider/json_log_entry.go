package log_provider

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
	"time"
)

type JsonLogEntry struct {
	raw            []byte
	TimeStampField string
	_map           map[string]interface{}
	_ts            time.Time
}

func (j *JsonLogEntry) RawString() string {
	return string(j.raw)
}

func (j *JsonLogEntry) RawBytes() []byte {
	return j.raw
}

func (j *JsonLogEntry) Map() map[string]interface{} { // TODO cache this
	if j._map != nil {
		return j._map
	}
	j._map = make(map[string]interface{})
	err := json.Unmarshal(j.raw, &j._map) // FIXME bytes
	if err != nil {
		log.Warnf("Error unmarshalling json log entry: %s", err.Error())
	}
	return j._map
}

// As json log does not have an id generate new uuid
func (j *JsonLogEntry) ID() string {
	return uuid.New().String()
}

func (j *JsonLogEntry) Timestamp() time.Time {
	if !j._ts.IsZero() {
		return j._ts
	}
	logMap := j.Map()
	ts, ok := logMap[j.TimeStampField]
	if !ok {
		log.Debugf("Timestamp field `%s` not found", j.TimeStampField)
		return time.Time{}
	}
	tsStr, ok := ts.(string)
	if !ok {
		log.Debugf("Timestamp field `%s` is not a string", j.TimeStampField)
		return time.Time{}
	}
	tsTime, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		log.Debugf("Timestamp field `%s` is not a valid RFC3339 string - %s", j.TimeStampField, logMap[j.TimeStampField])
		return time.Time{}
	}
	j._ts = tsTime
	return j._ts
}

func (j *JsonLogEntry) LoadString(raw string) error {
	j.raw = []byte(raw)
	return nil
}

func (j *JsonLogEntry) LoadBytes(raw []byte) error {
	j.raw = raw
	return nil
}

func (j *JsonLogEntry) LoadMap(raw map[string]interface{}) error {
	var err error
	j.raw, err = json.Marshal(raw)
	if err != nil {
		return err
	}
	return err
}

// JsonLogEntryConstructor creates a new JsonLogEntry struct
func JsonLogEntryConstructor() ports.LogEntry {
	return &JsonLogEntry{
		TimeStampField: "datetime",
	}
}
