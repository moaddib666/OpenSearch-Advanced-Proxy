package ports

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"context"
)

type SearchEngine interface {
	ProcessSearch(ctx context.Context, request *models.SearchRequest) ([]LogEntry, error)
}

type SearchFilter interface {
	Match(entry LogEntry) bool
}

type SearchFilterFactory interface {
	NewFilter(filter *models.Filter) (SearchFilter, error)
}
