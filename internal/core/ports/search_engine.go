package ports

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"context"
)

type SearchEngine interface {
	ProcessSearch(ctx context.Context, request *models.SearchRequest) ([]LogEntry, error)
}
