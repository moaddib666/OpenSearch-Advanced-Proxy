package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type AppConfig interface {
	AvailableIndexes() []*models.SubConfig
}

type ProviderConfig interface {
	GetProviderConfig(config interface{}) error
	GetProvider() models.ProviderType
}
