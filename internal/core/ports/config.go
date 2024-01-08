package ports

import "github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"

type AppConfig interface {
	AvailableIndexes() []*models.SubConfig
}

type ProviderConfig interface {
	GetProviderConfig(config interface{}) error
	GetProvider() models.ProviderType
}
