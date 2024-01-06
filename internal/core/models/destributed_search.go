package models

type DistributedSearchRequest struct {
	ID            string         `json:"id"`
	Index         string         `json:"index"`
	SearchRequest *SearchRequest `json:"search_request"`
}

type DistributedSearchResult struct {
	ID           string        `json:"id"`
	SearchResult *SearchResult `json:"search_result"`
}

// DistributedSearchResultFailed is a failure response
func DistributedSearchResultFailed(id string, timeout bool) *DistributedSearchResult {
	return &DistributedSearchResult{
		ID: id,
		SearchResult: &SearchResult{
			Took:     0,
			TimedOut: timeout,
			Shards: &Shards{
				Total:      1,
				Successful: 0,
				Skipped:    0,
				Failed:     1,
			},
			Hits:         &Hits{Total: &TotalValue{Value: 0}, Hits: make([]*Hit, 0)},
			Aggregations: make(map[string]*AggregationResult),
		},
	}
}

// DistributedSearchResultTimeout is a timeout response
func DistributedSearchResultTimeout(id string) *DistributedSearchResult {
	return DistributedSearchResultFailed(id, true)
}

// DistributedSearchResultSuccess is a success response
func DistributedSearchResultSuccess(id string, result *SearchResult) *DistributedSearchResult {
	return &DistributedSearchResult{
		ID:           id,
		SearchResult: result,
	}
}
