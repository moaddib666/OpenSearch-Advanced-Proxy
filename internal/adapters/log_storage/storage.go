package log_storage

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type GenericStorage struct {
	name      string
	fields    *models.Fields
	engine    ports.SearchEngine
	searchTTL time.Duration
}

// NewGenericStorage - create a new GenericStorage
func NewGenericStorage(name string, fields *models.Fields, engine ports.SearchEngine) *GenericStorage {
	return &GenericStorage{
		name:      name,
		fields:    fields,
		engine:    engine,
		searchTTL: 60 * time.Second,
	}
}

func (f *GenericStorage) Name() string {
	return f.name
}

func (f *GenericStorage) Fields() *models.Fields {
	return f.fields
}

func (f *GenericStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.searchTTL)
	defer cancel()
	start := time.Now()
	found, err := f.engine.ProcessSearch(ctx, r)
	timeTaken := int(time.Since(start).Milliseconds())
	if err != nil {
		log.Errorf("Error processing search: %s, took: %d nanoseconds", err.Error(), timeTaken)
	}
	count := len(found)
	hits := make([]*models.Hit, count)
	for i, entry := range found {
		hits[i] = &models.Hit{
			ID:     entry.ID(),
			Index:  f.name,
			Source: entry.Map(),
		}
	}

	successShardCount := 0
	failedShardCount := 0
	timeout := ctx.Err() != nil
	if !timeout {
		successShardCount = 1
	} else {
		failedShardCount = 1
	}
	return &models.SearchResult{
		Took:     timeTaken,
		TimedOut: timeout,
		Shards: &models.Shards{
			Total:      1,
			Successful: successShardCount,
			Skipped:    0,
			Failed:     failedShardCount,
		},
		Hits: &models.Hits{
			Total: &models.TotalValue{
				Value: count,
			},
			Hits: hits,
		},
	}, nil
}
