package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	request       *models.SearchRequest
	result        *models.SearchResult
	searchResults []*models.SearchResult
	mux           *sync.Mutex
}

// NewAggregator creates a new aggregator
func NewAggregator(request *models.SearchRequest) *Aggregator {
	return &Aggregator{
		searchResults: make([]*models.SearchResult, 0),
		mux:           &sync.Mutex{},
		request:       request,
		result: &models.SearchResult{
			Took:         0,
			TimedOut:     false,
			Shards:       nil,
			Hits:         nil,
			Aggregations: nil,
		},
	}
}

func (a *Aggregator) AddResult(result *models.SearchResult) {
	a.mux.Lock()
	a.searchResults = append(a.searchResults, result)
	a.mux.Unlock()
}

func parseHitTime(hit *models.Hit, fieldName string, sortIndex int) (time.Time, error) {
	// check if iter already in hit.Sort
	//if hit.Sort == nil {
	//	hit.Sort = make([]int, sortIndex+1)
	//}
	//if len(hit.Sort) > sortIndex {
	//	return time.Unix(int64(hit.Sort[sortIndex]), 0), nil
	//}

	value, ok := hit.Source[fieldName]
	if !ok {
		return time.Time{}, nil
	}
	switch value.(type) {
	case string:
		result, err := time.Parse(time.RFC3339, value.(string))
		//if err == nil {
		//	// insert to hit.Sort
		//	hit.Sort[sortIndex] = int(result.Unix())
		//}
		return result, err
	case time.Time:
		return value.(time.Time), nil
	default:
		return time.Time{}, nil
	}
}

// Sort sorts the hits by request parameters
func (a *Aggregator) Sort() {
	// TODO add sort abstraction
	if a.request.Sort == nil || len(a.request.Sort) == 0 {
		return
	}

	for ruleId, sortRule := range a.request.Sort {
		for fieldName, sortOrder := range sortRule {
			// Sorting logic
			sort.SliceStable(a.result.Hits.Hits, func(i, j int) bool {
				hitI := a.result.Hits.Hits[i]
				hitJ := a.result.Hits.Hits[j]

				timeI, errI := parseHitTime(hitI, fieldName, ruleId)
				timeJ, errJ := parseHitTime(hitJ, fieldName, ruleId)

				if errI != nil || errJ != nil {
					return false
				}
				hitI.Sort = append(hitI.Sort, int(timeI.Unix()))
				hitJ.Sort = append(hitJ.Sort, int(timeI.Unix()))

				if sortOrder.Order == "desc" {
					return timeI.After(timeJ)
				}
				return timeI.Before(timeJ)
			})
		}
	}
}

// Aggregate fill in the aggregations
func (a *Aggregator) Aggregate() {
	log.Warnf("Aggregation is not implemented yet")
}

func (a *Aggregator) GetResult() *models.SearchResult {
	a.mux.Lock()
	defer a.mux.Unlock()
	took := 0
	shards := &models.Shards{
		Total:      0,
		Successful: 0,
		Skipped:    0,
	}

	hits := &models.Hits{
		Total: &models.TotalValue{
			Value: 0,
		},
		Hits: make([]*models.Hit, 0),
	}

	for _, result := range a.searchResults {
		took += result.Took
		shards.Total += result.Shards.Total
		shards.Failed += result.Shards.Failed
		shards.Skipped += result.Shards.Skipped
		shards.Successful += result.Shards.Successful
		hits.Total.Value += result.Hits.Total.Value
		hits.Hits = append(hits.Hits, result.Hits.Hits...) // Not ordered yet
	}

	a.result.Took = took
	a.result.Shards = shards
	a.result.Hits = hits

	a.Sort()
	a.Aggregate()

	return a.result
}

type AggregatorFactory struct {
}

func (a *AggregatorFactory) CreateAggregator(request *models.SearchRequest) ports.SearchAggregator {
	return NewAggregator(request)
}

func NewAggregatorFactory() *AggregatorFactory {
	return &AggregatorFactory{}
}
