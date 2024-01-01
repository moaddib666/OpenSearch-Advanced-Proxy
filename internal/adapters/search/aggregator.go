package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"sync"
)

type Aggregator struct {
	searchResults []*models.SearchResult
	mux           *sync.Mutex
}

// NewAggregator creates a new aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		searchResults: make([]*models.SearchResult, 0),
		mux:           &sync.Mutex{},
	}
}

func (a *Aggregator) AddResult(result *models.SearchResult) {
	a.mux.Lock()
	a.searchResults = append(a.searchResults, result)
	a.mux.Unlock()
}

// sortHits sorts hits by timestamp
func (a *Aggregator) sortHits(hits []*models.Hit) {

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

	result := &models.SearchResult{
		Took:         took,
		TimedOut:     false,
		Shards:       shards,
		Hits:         hits,
		Aggregations: nil,
	}
	return result
}

type AggregatorFactory struct {
}

func (a *AggregatorFactory) CreateAggregator() ports.SearchAggregator {
	return NewAggregator()
}

func NewAggregatorFactory() *AggregatorFactory {
	return &AggregatorFactory{}
}
