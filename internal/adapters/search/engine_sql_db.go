package search

import (
	"OpenSearchAdvancedProxy/internal/adapters/tracker"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	log "github.com/sirupsen/logrus"
)

type SQLDBSearchEngine struct {
	provider    ports.SearchDataProvider
	convertor   ports.LogEntryConvertor
	aggregation ports.SearchAggregatorFactory
	tracker     ports.TimeTracker
}

// NewSQLDBSearchEngine - create a new SQLDBSearchEngine
func NewSQLDBSearchEngine(provider ports.SearchDataProvider, conv ports.LogEntryConvertor, aggregation ports.SearchAggregatorFactory) *SQLDBSearchEngine {
	return &SQLDBSearchEngine{
		provider:    provider,
		convertor:   conv,
		tracker:     tracker.NewDefaultTimeTracker(),
		aggregation: aggregation,
	}
}

func (s *SQLDBSearchEngine) ProcessSearch(ctx context.Context, request *models.SearchRequest) (*models.SearchResult, error) {
	hits := models.NewHits()
	s.tracker.Start()
	s.provider.BeginScan(request)
	defer s.provider.EndScan()
SearchLoop:
	for s.provider.Scan() {
		entry := s.provider.LogEntry()
		if entry == nil {
			continue
		}
		hit, err := s.convertor.Convert(entry)
		if err != nil {
			log.Errorf("Error converting entry: %s", err.Error())
			continue // FIXME: raise condition if all next entries are not converted deadline exceeded will not be raised
		}
		hits.AddHit(hit)
		select {
		case <-ctx.Done():
			log.Warnf("Search request canceled as timeout reached")
			break SearchLoop
		default:
			// do nothing
		}
	}
	s.tracker.Stop()
	timeout := ctx.Err() != nil
	// TBD: not sure if we need this in future
	if err := s.provider.Err(); err != nil {
		return nil, err
	}
	successShardCount := 0
	failedShardCount := 0

	if !timeout {
		successShardCount = 1
	} else {
		failedShardCount = 1
	}

	shards := &models.Shards{
		Total:      1,
		Successful: successShardCount,
		Skipped:    0,
		Failed:     failedShardCount,
	}
	aggregate := s.aggregation.CreateAggregator(request, s.provider)
	aggregate.AddResult(models.NewSearchResult(int(s.tracker.GetDuration().Seconds()), timeout, shards, hits))
	return aggregate.GetResult(), nil
}
