package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type SQLDBSearchEngine struct {
	provider ports.SearchDataProvider
}

// NewSQLDBSearchEngine - create a new SQLDBSearchEngine
func NewSQLDBSearchEngine(provider ports.SearchDataProvider) *SQLDBSearchEngine {
	return &SQLDBSearchEngine{
		provider: provider,
	}
}

func (s *SQLDBSearchEngine) ProcessSearch(ctx context.Context, request *models.SearchRequest) ([]ports.LogEntry, error) {
	var matchingLines []ports.LogEntry
	s.provider.BeginScan(request)
	defer s.provider.EndScan()
	for s.provider.Scan() {
		// check if context is done canceled
		select {
		case <-ctx.Done():
			log.Warnf("Search request canceled as timeout reached")
			return matchingLines, fmt.Errorf("search request canceled as timeout reached")
		default:
			// do nothing
		}
		entry := s.provider.LogEntry()
		if entry == nil {
			continue
		}
		matchingLines = append(matchingLines, entry)
	}

	if err := s.provider.Err(); err != nil {
		return nil, err
	}
	metadata := s.provider.SearchMetadata()
	if metadata != nil {
		jsonMetadata, _ := json.Marshal(metadata)
		log.Infof("Search metadata: %s", string(jsonMetadata))
	}
	return matchingLines, nil
}
