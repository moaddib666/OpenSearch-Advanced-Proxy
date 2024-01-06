package aggragation

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
)

type SQLDatabaseAggregator struct {
	provider ports.SearchMetadataProvider
	request  *models.SearchRequest
	result   *models.SearchResult
}

func (d *SQLDatabaseAggregator) AddResult(result *models.SearchResult) {
	d.result = result
}

func (d *SQLDatabaseAggregator) aggregate() {
	d.result.Aggregations = make(map[string]*models.AggregationResult)
	for name, agr := range d.request.Aggregations {
		d.result.Aggregations[name] = d.provider.AggregateResult(agr)
	}
}
func (d *SQLDatabaseAggregator) GetResult() *models.SearchResult {
	d.aggregate()
	return d.result
}

type SQLDatabaseAggregatorFactory struct{}

func (d *SQLDatabaseAggregatorFactory) CreateAggregator(request *models.SearchRequest, provider ports.SearchMetadataProvider) ports.SearchAggregator {
	return &SQLDatabaseAggregator{
		provider: provider,
		request:  request,
	}
}

func NewSQLDatabaseAggregatorFactory() *SQLDatabaseAggregatorFactory {
	return &SQLDatabaseAggregatorFactory{}
}
