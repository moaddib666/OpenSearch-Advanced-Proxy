package ports

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"time"
)

type SearchHitAggregator interface {
	Aggregate(hits []*models.Hit) *models.AggregationResult
}

type SearchHitAggregatorFactory interface {
	CreateAggregator(settings *models.SearchAggregation) SearchHitAggregator
}

type HitTimeParser func(hit *models.Hit, fieldName string) (time.Time, error)

type SearchAggregator interface {
	// AddResult adds a result to the aggregator
	AddResult(result *models.SearchResult)
	GetResult() *models.SearchResult
}

type SearchAggregatorFactory interface {
	// CreateAggregator creates a new aggregator
	CreateAggregator(request *models.SearchRequest) SearchAggregator
}