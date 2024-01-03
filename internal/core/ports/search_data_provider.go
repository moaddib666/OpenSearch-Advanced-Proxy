package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type SearchDataProvider interface {
	BeginScan(size int, r *models.Range, s *models.SortOrder)
	Scan() bool
	Text() string
	Err() error
	LogEntry() LogEntry
	EndScan()
}
