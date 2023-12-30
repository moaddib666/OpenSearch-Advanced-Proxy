package models

// SearchRequest represents the overall structure of an OpenSearch request
type SearchRequest struct {
	Sort           []map[string]*SortOrder `json:"sort"`
	Size           int                     `json:"size"`
	Version        bool                    `json:"version"`
	StoredFields   []string                `json:"stored_fields"`
	ScriptFields   map[string]interface{}  `json:"script_fields"`
	DocvalueFields []*DocvalueField        `json:"docvalue_fields"`
	Source         *SourceSetting          `json:"_source"`
	Query          *Query                  `json:"query"`
	Highlight      *Highlight              `json:"highlight"`
}

// SortOrder represents sorting options
type SortOrder struct {
	Order string `json:"order"`
}

// DocvalueField represents docvalue fields settings
type DocvalueField struct {
	Field  string `json:"field"`
	Format string `json:"format"`
}

// SourceSetting represents settings for _source field
type SourceSetting struct {
	Excludes []string `json:"excludes"`
}

// Query represents the query structure
type Query struct {
	Bool *BoolQuery `json:"bool,omitempty"`
	//MultiMatch *MultiMatch `json:"multi_match"`
}

type Filter struct {
	Bool       *BoolQuery  `json:"bool,omitempty"`
	MultiMatch *MultiMatch `json:"multi_match"`
}

// BoolQuery represents a boolean query
type BoolQuery struct {
	//Must    *MatchQuery `json:"must,omitempty"`
	Filter []*Filter `json:"filter,omitempty"`
	//Should  *MatchQuery `json:"should,omitempty"`
	//MustNot *MatchQuery `json:"must_not,omitempty"`
}

// MatchQuery represents a match query
type MatchQuery struct {
	MultiMatch *MultiMatch `json:"multi_match"`
}

// MultiMatch represents a multi-match query
type MultiMatch struct {
	Type    string `json:"type"`
	Query   string `json:"query"`
	Lenient bool   `json:"lenient"`
}

// Highlight represents highlight settings
type Highlight struct {
	PreTags      []string               `json:"pre_tags"`
	PostTags     []string               `json:"post_tags"`
	Fields       map[string]interface{} `json:"fields"`
	FragmentSize int                    `json:"fragment_size"`
}
