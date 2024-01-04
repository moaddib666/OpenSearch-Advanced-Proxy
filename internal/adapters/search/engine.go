package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type LogSearchEngine struct {
	provider      ports.SearchDataProvider
	filterFactory *FilterFactory
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
func NewLogSearchEngine(provider ports.SearchDataProvider) *LogSearchEngine {
	return &LogSearchEngine{
		provider:      provider,
		filterFactory: NewFilterFactory(),
	}
}

func (se *LogSearchEngine) ProcessSearch(ctx context.Context, request *models.SearchRequest) ([]ports.LogEntry, error) {
	var matchingLines []ports.LogEntry
	// ---------------------------- begin ----------------------------
	// FIXME: currently support only one range and sort order
	if len(request.DocvalueFields) != 1 {
		return nil, models.ErrUnsupportedDocvalueFields
	}
	docValueField := request.DocvalueFields[0].Field
	srt := request.Sort[0][docValueField]
	var rg *models.Range
	for _, filter := range request.Query.Bool.Filter {
		if filter.Range != nil {
			rg = filter.Range
			break
		}
	}
	// ---------------------------- end  ----------------------------
	se.provider.BeginScan(request.Size, rg, srt)
	defer se.provider.EndScan()
	for se.provider.Scan() {
		// check if context is done canceled
		select {
		case <-ctx.Done():
			log.Warnf("Search request canceled as timeout reached")
			return matchingLines, fmt.Errorf("search request canceled as timeout reached")
		default:
			// do nothing
		}

		entry := se.provider.LogEntry()
		if entry == nil {
			continue
		}
		// Performance optimization: skip entries that are outside the time range before applying filters
		if entry.Timestamp().Before(rg.DateTime.GTE) {
			//log.Debugf("Entry %s is before range %s", entry.Timestamp(), rg.DateTime.GTE)
			continue
		}
		if entry.Timestamp().After(rg.DateTime.LTE) {
			//log.Debugf("Entry %s is after range %s", entry.Timestamp(), rg.DateTime.LTE)
			break
		}
		filter, err := se.filterFactory.FromQuery(request.Query)
		if err != nil {
			return nil, err
		}
		if !filter.Match(entry) {
			continue
		}
		matchingLines = append(matchingLines, entry)
	}

	if err := se.provider.Err(); err != nil {
		return nil, err
	}

	return matchingLines, nil
}
