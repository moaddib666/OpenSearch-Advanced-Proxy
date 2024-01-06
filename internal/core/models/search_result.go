package models

// SearchResult is a struct that represents the result of a search query
type SearchResult struct {
	Took         int                           `json:"took"`
	TimedOut     bool                          `json:"timed_out"`
	Shards       *Shards                       `json:"_shards"`
	Hits         *Hits                         `json:"hits"`
	Aggregations map[string]*AggregationResult `json:"aggregations"`
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
	Index  string                 `json:"_index,omitempty"`
	Type   string                 `json:"_type,omitempty"`
	ID     string                 `json:"_id,omitempty"`
	Score  float64                `json:"_score,omitempty"`
	Source map[string]interface{} `json:"_source,omitempty"`
	Fields interface{}            `json:"fields,omitempty"`
	Sort   HitSort                `json:"sort,omitempty"`
}
type HitSort []int

type AggregationResult struct {
	Buckets []*Bucket `json:"buckets"`
}

func (ar *AggregationResult) AddBucket(bucket *Bucket) {
	ar.Buckets = append(ar.Buckets, bucket)
}

func (ar *AggregationResult) DocsCount() int {
	count := 0
	for _, bucket := range ar.Buckets {
		count += bucket.DocCount
	}
	return count
}

type Bucket struct {
	KeyAsString string `json:"key_as_string"`
	Key         int64  `json:"key"`
	DocCount    int    `json:"doc_count"`
}
