package log_storage

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
	"time"
)

type FileStorage struct {
	name   string
	fields *models.Fields
	engine ports.SearchEngine
}

// NewFileStorage - create a new FileStorage
func NewFileStorage(name string, fields *models.Fields, engine ports.SearchEngine) *FileStorage {
	return &FileStorage{
		name:   name,
		fields: fields,
		engine: engine,
	}
}

func (f *FileStorage) Name() string {
	return f.name
}

func (f *FileStorage) Fields() *models.Fields {
	return f.fields
}

func (f *FileStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	start := time.Now()
	found, err := f.engine.ProcessSearch(r)
	timeTaken := int(time.Since(start).Milliseconds())
	if err != nil {
		log.Errorf("Error processing search: %s", err.Error())
		return nil, err
	}
	count := len(found)
	hits := make([]*models.Hit, count)
	for i, entry := range found {
		hits[i] = &models.Hit{
			Index:  f.name,
			Source: entry.Map(),
		}
	}

	return &models.SearchResult{
		Took:     timeTaken,
		TimedOut: false,
		Shards: &models.Shards{
			Total:      1,
			Successful: 1,
			Skipped:    0,
			Failed:     0,
		},
		Hits: &models.Hits{
			Total: &models.TotalValue{
				Value: count,
			},
			Hits: hits,
		},
	}, nil
}
