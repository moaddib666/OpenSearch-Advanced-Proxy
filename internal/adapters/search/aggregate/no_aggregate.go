package aggregate

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
)

type NoAggregation struct {
	result *models.SearchResult
}

func (n *NoAggregation) AddResult(result *models.SearchResult) {
	n.result = result
}

func (n *NoAggregation) GetResult() *models.SearchResult {
	return n.result
}

type NoAggregationFactory struct {
}

func (n *NoAggregationFactory) CreateAggregator(request *models.SearchRequest, provider ports.SearchMetadataProvider) ports.SearchAggregator {
	return &NoAggregation{}
}

func NewNoAggregationFactory() *NoAggregationFactory {
	return &NoAggregationFactory{}
}
