package ports

import "OpenSearchAdvancedProxy/internal/core/models"

// Storage is the interface that represents the external log storage.
type Storage interface {
	Name() string
	Fields() *models.Fields
	Search(r *models.SearchRequest) (*models.SearchResult, error)
}

// SearchFunc is a function that searches the storage.
type SearchFunc func(r *models.SearchRequest) (*models.SearchResult, error)
