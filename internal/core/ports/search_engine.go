package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type SearchEngine interface {
	ProcessSearch(request *models.SearchRequest) ([]LogEntry, error)
}
