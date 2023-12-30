package models

// SearchRequest is a struct that represents a search query
type SearchRequest struct {
	// FIXME: add the rest of the fields
}

// SearchResult is a struct that represents the result of a search query
type SearchResult struct {
	Took         int           `json:"took"`
	TimedOut     bool          `json:"timed_out"`
	Shards       *Shards       `json:"_shards"`
	Hits         *Hits         `json:"hits"`
	Aggregations *Aggregations `json:"aggregations"` // FIXME: does not work as expected
}

// NewSearchResult creates a new SearchResult struct
func NewSearchResult(took int, timedOut bool, shards *Shards, hits *Hits) *SearchResult {
	return &SearchResult{
		Took:     took,
		TimedOut: timedOut,
		Shards:   shards,
		Hits:     hits,
	}

}

// Shards is a struct that represents a shard
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total    *TotalValue `json:"total"`
	MaxScore float64     `json:"max_score"`
	Hits     []*Hit      `json:"hits"`
}

type TotalValue struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type Hit struct {
	Index  string      `json:"_index,omitempty"`
	Type   string      `json:"_type,omitempty"`
	ID     string      `json:"_id,omitempty"`
	Score  float64     `json:"_score,omitempty"`
	Source interface{} `json:"_source,omitempty"`
}

type Aggregations struct {
	BucketsAggregate map[string]*Buckets `json:"2"` // Use the appropriate key that matches your JSON structure
}

type Buckets struct {
	Buckets []*Bucket `json:"buckets"`
}

type Bucket struct {
	KeyAsString string `json:"key_as_string"`
	Key         int64  `json:"key"`
	DocCount    int    `json:"doc_count"`
}
