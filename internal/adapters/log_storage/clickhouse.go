package log_storage

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	_ "github.com/ClickHouse/clickhouse-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type ClickhouseStorage struct {
	name      string
	fields    *models.Fields
	engine    ports.SearchEngine
	searchTTL time.Duration
}

func (c *ClickhouseStorage) Name() string {
	return c.name
}

func (c *ClickhouseStorage) Fields() *models.Fields {
	return c.fields
}

func (c *ClickhouseStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	// FIXME: DRY looks like this abstraction is not needed lets reuse 1:1
	ctx, cancel := context.WithTimeout(context.Background(), c.searchTTL)
	defer cancel()
	start := time.Now()
	found, err := c.engine.ProcessSearch(ctx, r)
	timeTaken := int(time.Since(start).Milliseconds())
	if err != nil {
		log.Errorf("Error processing search: %s, took: %d nanoseconds", err.Error(), timeTaken)
	}
	count := len(found)
	hits := make([]*models.Hit, count)
	for i, entry := range found {
		hits[i] = &models.Hit{
			ID:     entry.ID(),
			Index:  c.name,
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

func NewClickhouseStorage(name string, fields *models.Fields, engine ports.SearchEngine) *ClickhouseStorage {
	return &ClickhouseStorage{
		name:      name,
		fields:    fields,
		engine:    engine,
		searchTTL: 60 * time.Second, // FIXME: Hardcoded
	}
}
