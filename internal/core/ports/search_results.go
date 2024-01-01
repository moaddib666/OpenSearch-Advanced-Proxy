package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type SearchAggregator interface {
	// AddResult adds a result to the aggregator
	AddResult(result *models.SearchResult)
	GetResult() *models.SearchResult
}

type SearchAggregatorFactory interface {
	// CreateAggregator creates a new aggregator
	CreateAggregator() SearchAggregator
}
