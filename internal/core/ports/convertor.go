package ports

import (
	"OpenSearchAdvancedProxy/internal/core/models"
)

type LogEntryConvertor interface {
	Convert(entry LogEntry) (*models.Hit, error)
	ConvertBatch(entries []LogEntry) (*models.Hits, error)
}
