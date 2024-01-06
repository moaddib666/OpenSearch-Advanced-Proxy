package ports

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"context"
)

type SearchEngine interface {
	ProcessSearch(ctx context.Context, request *models.SearchRequest) (*models.SearchResult, error)
}

type SearchFilter interface {
	Match(entry LogEntry) bool
}

type SearchFilterFactory interface {
	CreateFilter(filter *models.Filter) (SearchFilter, error)
	CreateBoolFilter(filter *models.BoolFilter) (SearchFilter, error)
	FromQuery(query *models.Query) (SearchFilter, error)
}

type QueryBuilder interface {
	BuildQuery() (string, error)
}

type QueryBuilderFactory interface {
	CreateQueryBuilder(filter *models.Filter) (QueryBuilder, error)
	CreateBoolConditionBuilder(filter *models.BoolFilter) (QueryBuilder, error)
	FromQuery(query *models.Query) (QueryBuilder, error)
}
