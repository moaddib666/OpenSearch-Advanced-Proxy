package indexer

import (
	"OpenSearchAdvancedProxy/internal/adapters/lock"
	"OpenSearchAdvancedProxy/internal/core/models"
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type JsonFileIndexer struct {
	file    string
	tsField string
	Index   *models.Index
	locker  *lock.TryLocker // FIXME move to pkg/lock
}

func (j *JsonFileIndexer) SearchStartPos(ts time.Time) (int64, error) {
	var index int64
	for i, entry := range j.Index.Entries {
		if ts.After(entry.Timestamp) || ts.Equal(entry.Timestamp) {
			index = j.Index.Entries[i].Position
		}
	}
	return index, nil
}

func (j *JsonFileIndexer) SearchEndPos(ts time.Time) (int64, error) {
	// FIXME implement me
	for _, entry := range j.Index.Entries {
		if entry.Timestamp.Equal(ts.Truncate(j.Index.Step)) || entry.Timestamp.After(ts) {
			return entry.Position, nil
		}
	}
	return -1, fmt.Errorf("end position not found")
}

// NewJsonFileIndexer creates a new JsonFileIndexer struct
func NewJsonFileIndexer(file string, tsField string, resolution time.Duration) *JsonFileIndexer {
	return &JsonFileIndexer{
		file:    file,
		tsField: tsField,
		Index: &models.Index{
			Entries: make([]*models.IndexEntry, 0),
			Step:    resolution,
		},
		locker: lock.NewTryLocker(),
	}
}

func (j *JsonFileIndexer) GetIndexName() string {
	//TODO implement me
	panic("implement me")
}

func (j *JsonFileIndexer) CreateIndex() error {
	log.Infof("Creating index %s", j.file+".index")
	file, err := os.Open(j.file)
	if err != nil {
		return err
	}
	defer file.Close()

	var lastTimestamp time.Time
	var position int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var logEntry map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &logEntry); err != nil {
			continue // Skip lines that can't be parsed
		}

		// Extract and parse the timestamp
		if ts, ok := logEntry[j.tsField].(string); ok {
			timestamp, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				continue // Skip lines with invalid timestamp
			}

			// Check if the timestamp is in a new interval
			if lastTimestamp.IsZero() || !isInSameInterval(lastTimestamp, timestamp, j.Index.Step) {
				j.Index.Entries = append(j.Index.Entries, &models.IndexEntry{
					Timestamp: timestamp,
					Position:  position,
				})
				lastTimestamp = timestamp
			}
		}
		position += int64(len(scanner.Bytes())) + 1 // +1 for newline character
	}
	if scanner.Err() != nil {
		log.Errorf("Error creating index %s: %s", j.file, scanner.Err())
		return scanner.Err()
	}
	log.Infof("Indexing complete, saving index %s", j.file+".index")
	return j.SaveIndex()
}

// Helper function to check if two timestamps are in the same interval
func isInSameInterval(a, b time.Time, interval time.Duration) bool {
	return a.Truncate(interval).Equal(b.Truncate(interval))
}

func (j *JsonFileIndexer) LoadIndex() error {
	log.Infof("Loading index %s", j.file+".index")
	bytes, err := os.ReadFile(j.file + ".index")
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, j.Index)
}

func (j *JsonFileIndexer) SaveIndex() error {
	bytes, err := json.Marshal(j.Index)
	if err != nil {
		return err
	}
	return os.WriteFile(j.file+".index", bytes, 0644)
}

func (j *JsonFileIndexer) LoadOrCreateIndex() error {
	// FIXME lock file to prevent concurrent access
	if !j.locker.TryLock() {
		return fmt.Errorf("indexing already in progress")
	}
	defer j.locker.Unlock()
	if fInfo, err := os.Stat(j.file + ".index"); os.IsNotExist(err) {
		log.Debugf("Creating index for %s", j.file+".index")
		return j.CreateIndex()
	} else {
		if fInfo.ModTime().Before(time.Now().Add(-1 * j.Index.Step)) {
			log.Debugf("Updating index for %s", j.file+".index")
			// TODO: Continue index creation from last position
			// 	- If last position is not found in file, create index from scratch
			//return j.CreateIndex()
		}
	}
	return j.LoadIndex()
}
