package ports

import (
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
)

type LogEntryConvertor interface {
	Convert(entry LogEntry) (*models.Hit, error)
	ConvertBatch(entries []LogEntry) (*models.Hits, error)
}
