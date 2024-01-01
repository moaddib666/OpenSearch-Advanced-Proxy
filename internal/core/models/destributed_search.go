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
