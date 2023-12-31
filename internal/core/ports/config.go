package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type AppConfig interface {
	AvailableIndexes() []*models.SubConfig
}
