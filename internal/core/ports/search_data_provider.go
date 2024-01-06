package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type SearchMetadataProvider interface {
	AggregateResult(request *models.SearchAggregation) *models.AggregationResult
}

type SearchDataProvider interface {
	SearchMetadataProvider
	BeginScan(r *models.SearchRequest)
	Scan() bool
	Text() string
	Err() error
	LogEntry() LogEntry
	EndScan()
}
