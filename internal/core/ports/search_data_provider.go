package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type SearchDataProvider interface {
	BeginScan(r *models.SearchRequest)
	Scan() bool
	Text() string
	Err() error
	LogEntry() LogEntry
	EndScan()
	SearchMetadata() *models.OngoingSearchMetadata
}
