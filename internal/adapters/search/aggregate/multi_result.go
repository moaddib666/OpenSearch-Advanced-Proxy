package aggregate

import (
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
	"sort"
)

type MultiResultAggregate struct {
	provider ports.SearchMetadataProvider
	request  *models.SearchRequest

	rawResults      []*models.SearchResult
	rawHits         []*models.Hit
	rawAggregations map[string]*models.AggregationResult
	result          *models.SearchResult
}

func NewMultiResultAggregate(request *models.SearchRequest, provider ports.SearchMetadataProvider) *MultiResultAggregate {
	return &MultiResultAggregate{
		provider:        provider,
		request:         request,
		rawResults:      make([]*models.SearchResult, 0),
		rawHits:         make([]*models.Hit, 0),
		rawAggregations: make(map[string]*models.AggregationResult),
		result: &models.SearchResult{
			Took:     0,
			TimedOut: false,
			Shards: &models.Shards{
				Total:      0,
				Successful: 0,
				Skipped:    0,
				Failed:     0,
			},
			Hits: &models.Hits{
				Total: &models.TotalValue{
					Value: 0,
				},
			},
		},
	}
}

func (d *MultiResultAggregate) AddResult(result *models.SearchResult) {
	d.rawResults = append(d.rawResults, result)
}

func (d *MultiResultAggregate) aggregate() {
	for _, r := range d.rawResults {
		d.result.Took += r.Took
		d.result.Shards.Total += r.Shards.Total
		d.result.Shards.Failed += r.Shards.Failed
		d.result.Shards.Skipped += r.Shards.Skipped
		d.result.Shards.Successful += r.Shards.Successful
		d.result.Hits.Total.Value += r.Hits.Total.Value
		d.rawHits = append(d.rawHits, r.Hits.Hits...) // Not ordered yet

		for name, agr := range r.Aggregations {
			if d.rawAggregations[name] == nil {
				d.rawAggregations[name] = agr
			} else {
				d.rawAggregations[name].Buckets = append(d.rawAggregations[name].Buckets, agr.Buckets...)
			}
		}
	}
	d.processRawHits()
	d.aggregateRawAggregations()
}

func (d *MultiResultAggregate) processRawHits() {
	sort.Slice(d.rawHits, func(i, j int) bool {
		hit1 := d.rawHits[i]
		hit2 := d.rawHits[j]
		return hit1.IsBeforeHit(hit2)
	})
	// shape raw hits to result.Size
	if d.request.Size > 0 && len(d.rawHits) > d.request.Size {
		d.rawHits = d.rawHits[:d.request.Size]
	}
	d.result.Hits.Hits = d.rawHits
}

func (d *MultiResultAggregate) aggregateRawAggregations() {
	d.result.Aggregations = make(map[string]*models.AggregationResult)
	// TODO Optimize this, to avoit multieple iterations and sourtin on each iteration
	for name, agr := range d.rawAggregations {

		// Convert existing buckets into a map for quick lookup
		tmpBucketMap := make(map[int64]*models.Bucket) // Replace KeyType with the actual type of the bucket key
		// Merge or add buckets
		for _, bucket := range agr.Buckets {
			if existingBucket, exists := tmpBucketMap[bucket.Key]; exists {
				existingBucket.DocCount += bucket.DocCount
			} else {
				tmpBucketMap[bucket.Key] = bucket
			}
		}

		// Convert map back to slice
		updatedBuckets := make([]*models.Bucket, 0, len(tmpBucketMap))
		for _, bucket := range tmpBucketMap {
			updatedBuckets = append(updatedBuckets, bucket)
		}

		// Sort buckets by key
		sort.Slice(updatedBuckets, func(i, j int) bool {
			return updatedBuckets[i].Key < updatedBuckets[j].Key
		})

		d.result.Aggregations[name] = &models.AggregationResult{
			Buckets: updatedBuckets,
		}
	}
}

func (d *MultiResultAggregate) GetResult() *models.SearchResult {
	d.aggregate()
	for name, agr := range d.result.Aggregations {
		log.Debugf("%T: GetResult got %d bukets for %s", d, len(agr.Buckets), name)
	}
	return d.result
}

type MultiResultAggregateFactory struct{}

func NewMultiResultAggregateFactory() *MultiResultAggregateFactory {
	return &MultiResultAggregateFactory{}
}
func (m *MultiResultAggregateFactory) CreateAggregator(request *models.SearchRequest, provider ports.SearchMetadataProvider) ports.SearchAggregator {
	return NewMultiResultAggregate(request, provider)
}
