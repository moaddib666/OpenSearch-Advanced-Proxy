package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
	"strings"
)

type LogSearchEngine struct {
	provider ports.SearchDataProvider
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
		provider: provider,
	}
}

func (se *LogSearchEngine) matchMultiMatchQuery(multiMatch *models.MultiMatch, entry ports.LogEntry) bool {
	if multiMatch.Type == "phrase" {
		// FIXME: This is a very naive implementation
		log.Debugf("Matching phrase: %s", multiMatch.Query)
		return strings.Contains(entry.Raw(), multiMatch.Query)
	}
	if multiMatch.Type == "best_fields" {
		// FIXME: This is a very naive implementation
		log.Debugf("Matching best_fields: %s", multiMatch.Query)
		return strings.Contains(entry.Raw(), multiMatch.Query)
	}
	// Additional conditions for other MultiMatch types can be added here
	return false
}

func (se *LogSearchEngine) recursivelyProcessBoolQuery(boolQuery *models.BoolQuery, entry ports.LogEntry) (result *LogSearchResult) {
	result = NewLogSearchResult(len(boolQuery.Filter))
	if boolQuery == nil {
		result.Found()
		return
	}

	// Process 'filter' queries
	for _, filter := range boolQuery.Filter {
		if filter.Range != nil {
			if filter.Range.DateTime != nil && filter.Range.DateTime.Format == "strict_date_optional_time" {
				log.Debugf("Range query: %s", filter.Range.DateTime.GTE)
				log.Debugf("Range query: %s", filter.Range.DateTime.LTE)
				if entry.Timestamp().Before(filter.Range.DateTime.GTE) {
					log.Debugf("Entry %s is before range %s", entry.Timestamp(), filter.Range.DateTime.GTE)
					result.TimeRangeExcluded()
					return
				}
				if entry.Timestamp().After(filter.Range.DateTime.LTE) {
					log.Debugf("Entry %s is after range %s", entry.Timestamp(), filter.Range.DateTime.LTE)
					result.TimeRangeExcluded()
					return
				}
				result.Match()
			}
		}
		if filter.MatchAll != nil {
			result.Match()
			continue
		}
		if filter.Bool != nil {
			if se.recursivelyProcessBoolQuery(filter.Bool, entry).IsFound() {
				result.Match()
				continue
			}
		}
		if filter.MultiMatch != nil {
			if se.matchMultiMatchQuery(filter.MultiMatch, entry) {
				result.Match()
				continue
			}
		}
	}
	// Similar logic can be implemented for 'must', 'should', and 'must_not'
	// if they are required for your use case
	return
}

func (se *LogSearchEngine) ProcessSearch(request *models.SearchRequest) ([]ports.LogEntry, error) {
	var matchingLines []ports.LogEntry
	se.provider.BeginScan()
	for se.provider.Scan() {
		entry := se.provider.LogEntry()
		if request.Query != nil && request.Query.Bool != nil {
			log.Debugf("Processing entry: %s", entry.Raw())
			match := se.recursivelyProcessBoolQuery(request.Query.Bool, entry)
			if match.IsFound() {
				log.Debugf("Entry matches query: %s", entry.Raw())
				matchingLines = append(matchingLines, entry)
			} else if match.IsTimeRangeExcluded() {
				log.Debugf("Entry is out of time range stop searching: %s", entry.Raw())
				se.provider.EndScan()
				break
			}

		}
	}

	if err := se.provider.Err(); err != nil {
		return nil, err
	}

	return matchingLines, nil
}
