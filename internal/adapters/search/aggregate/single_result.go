package aggregate

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
)

type SingleResultAggregate struct {
	provider ports.SearchMetadataProvider
	request  *models.SearchRequest
	result   *models.SearchResult
}

func (d *SingleResultAggregate) AddResult(result *models.SearchResult) {
	d.result = result
}

func (d *SingleResultAggregate) aggregate() {
	d.result.Aggregations = make(map[string]*models.AggregationResult)
	for name, agr := range d.request.Aggregations {
		agrResult := d.provider.AggregateResult(agr)
		d.result.Aggregations[name] = agrResult
		for _, bucket := range agrResult.Buckets {
			d.result.Hits.Total.Value += bucket.DocCount
		}
	}
}

func (d *SingleResultAggregate) GetResult() *models.SearchResult {
	d.aggregate()
	for name, agr := range d.result.Aggregations {
		log.Debugf("%T: GetResult got %d bukets for %s", d, len(agr.Buckets), name)
	}
	return d.result
}

type SingleResultAggregateFactory struct{}

func (d *SingleResultAggregateFactory) CreateAggregator(request *models.SearchRequest, provider ports.SearchMetadataProvider) ports.SearchAggregator {
	return &SingleResultAggregate{
		provider: provider,
		request:  request,
	}
}

func NewSingleResultAggregateFactory() *SingleResultAggregateFactory {
	return &SingleResultAggregateFactory{}
}
