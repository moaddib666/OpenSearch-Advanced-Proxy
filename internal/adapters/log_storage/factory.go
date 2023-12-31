package log_storage

import (
	"OpenSearchAdvancedProxy/internal/adapters/log_provider"
	"OpenSearchAdvancedProxy/internal/adapters/search"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
)

type BaseStorageFactory struct {
}

// NewBaseStorageFactory - create new BaseStorageFactory
func NewBaseStorageFactory() *BaseStorageFactory {
	return &BaseStorageFactory{}
}

func (b *BaseStorageFactory) FromConfig(name string, config *models.SubConfig) (ports.Storage, error) {
	if config.Version != models.ConfigVersion1 {
		return nil, models.ErrUnsupportedVersion
	}
	if config.Fields == nil {
		return nil, models.ErrNoFields
	}
	fields := &models.Fields{}
	for fieldName, field := range config.Fields {
		fields.AddField(fieldName, field)
	}

	if config.Provider == models.JsonLogFileProvider {
		logFile, ok := config.ProviderConfig["logfile"]
		if !ok {
			return nil, models.ErrNoLogFile
		}
		provider := log_provider.NewLogFileProvider(logFile,
			func() ports.LogEntry {
				return &log_provider.JsonLogEntry{
					TimeStampField: config.Timestamp.Field,
				}
			})
		engine := search.NewLogSearchEngine(provider)
		return NewFileStorage(name, fields, engine), nil
	} else {
		return nil, models.ErrUnsupportedProvider
	}
}
