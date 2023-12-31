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

func (se *LogSearchEngine) recursivelyProcessBoolQuery(boolQuery *models.BoolQuery, entry ports.LogEntry) bool {
	if boolQuery == nil {
		return true
	}

	// Process 'filter' queries
	for _, filter := range boolQuery.Filter {
		// TODO Implement range
		if filter.MatchAll != nil {
			return true
		}
		if filter.Bool != nil && !se.recursivelyProcessBoolQuery(filter.Bool, entry) {
			return false
		}
		if filter.MultiMatch != nil && !se.matchMultiMatchQuery(filter.MultiMatch, entry) {
			return false
		}
	}
	// Similar logic can be implemented for 'must', 'should', and 'must_not'
	// if they are required for your use case
	return true
}

func (se *LogSearchEngine) ProcessSearch(request *models.SearchRequest) ([]ports.LogEntry, error) {
	var matchingLines []ports.LogEntry
	for se.provider.Scan() {
		entry := se.provider.LogEntry()
		if request.Query != nil && request.Query.Bool != nil {
			log.Debugf("Processing entry: %s", entry.Raw())
			if match := se.recursivelyProcessBoolQuery(request.Query.Bool, entry); match {
				log.Debugf("Entry matches query: %s", entry.Raw())
				matchingLines = append(matchingLines, entry)
			}
		}
	}

	if err := se.provider.Err(); err != nil {
		return nil, err
	}

	return matchingLines, nil
}
