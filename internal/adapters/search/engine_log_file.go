package search

import (
	"OpenSearchAdvancedProxy/internal/adapters/tracker"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	log "github.com/sirupsen/logrus"
)

type LogSearchEngine struct {
	provider      ports.SearchDataProvider
	filterFactory *FilterFactory
	convertor     ports.LogEntryConvertor
	aggregation   ports.SearchAggregatorFactory
	tracker       ports.TimeTracker
}

type LogSearchResult struct {
	matchCount        int
	filtersCount      int
	excluded          bool
	timeRangeExcluded bool
	outOfTimeRange    bool
}

// Found - set logSearchResult.found to true
func (lsr *LogSearchResult) Found() {
	lsr.matchCount = lsr.filtersCount
}

func (lsr *LogSearchResult) Match() {
	lsr.matchCount += 1
}

// IsFound - return logSearchResult.found
func (lsr *LogSearchResult) IsFound() bool {
	// found is valid found len == filters len

	return lsr.matchCount >= lsr.filtersCount && !lsr.excluded && !lsr.timeRangeExcluded
}

// Exclude - set logSearchResult.excluded to true
func (lsr *LogSearchResult) Exclude() {
	lsr.excluded = true
}

// TimeRangeExcluded - set logSearchResult.timeRangeExcluded to true
func (lsr *LogSearchResult) TimeRangeExcluded() {
	lsr.timeRangeExcluded = true
}

// IsTimeRangeExcluded - return logSearchResult.timeRangeExcluded
func (lsr *LogSearchResult) IsTimeRangeExcluded() bool {
	return lsr.timeRangeExcluded
}

// OutOfTimeRange - set logSearchResult.outOfTimeRange to true
func (lsr *LogSearchResult) OutOfTimeRange() {
	lsr.outOfTimeRange = true
}

// IsOutOfTimeRange - return true if logSearchResult.timeRangeExcluded is true
func (lsr *LogSearchResult) IsOutOfTimeRange() bool {
	return lsr.outOfTimeRange
}

// NewLogSearchResult - create a new logSearchResult
func NewLogSearchResult(filtersCount int) *LogSearchResult {
	return &LogSearchResult{
		matchCount:        0,
		excluded:          false,
		timeRangeExcluded: false,
		filtersCount:      filtersCount,
	}
}

// NewLogSearchEngine - create a new LogSearchEngine
func NewLogSearchEngine(provider ports.SearchDataProvider, conv ports.LogEntryConvertor, aggregation ports.SearchAggregatorFactory) *LogSearchEngine {
	return &LogSearchEngine{
		provider:      provider,
		filterFactory: NewFilterFactory(),
		convertor:     conv,
		tracker:       tracker.NewDefaultTimeTracker(),
		aggregation:   aggregation,
	}
}

func (s *LogSearchEngine) ProcessSearch(ctx context.Context, request *models.SearchRequest) (*models.SearchResult, error) {
	hits := models.NewHits()
	rg := request.GetRange()
	s.tracker.Start()
	s.provider.BeginScan(request)
	defer s.provider.EndScan()
	filter, err := s.filterFactory.FromQuery(request.Query)
	if err != nil {
		return nil, err
	}
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
		// Performance optimization: skip entries that are outside the time range before applying filters
		if entry.Timestamp().Before(rg.DateTime.GTE) {
			//log.Debugf("Entry %s is before range %s", entry.Timestamp(), rg.DateTime.GTE)
			continue
		}
		if entry.Timestamp().After(rg.DateTime.LTE) {
			//log.Debugf("Entry %s is after range %s", entry.Timestamp(), rg.DateTime.LTE)
			break
		}
		select {
		case <-ctx.Done():
			log.Warnf("Search request canceled as timeout reached")
			break SearchLoop
		default:
			// do nothing
		}
		if !filter.Match(entry) {
			continue
		}
	}
	s.tracker.Stop()
	successShardCount := 0
	failedShardCount := 0
	timeout := ctx.Err() != nil
	// TBD: not sure if we need this in future
	if err := s.provider.Err(); err != nil {
		return nil, err
	}
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
