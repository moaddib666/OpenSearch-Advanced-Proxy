package aggregate

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
)

type NoAggregation struct {
	result *models.SearchResult
}

func (n *NoAggregation) AddResult(result *models.SearchResult) {
	n.result = result
}

func (n *NoAggregation) GetResult() *models.SearchResult {
	if n.result.Aggregations == nil {
		n.result.Aggregations = make(map[string]*models.AggregationResult)
	}
	for name, agr := range n.result.Aggregations {
		log.Debugf("%T: GetResult got %d bukets for %s", n, len(agr.Buckets), name)
	}
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
