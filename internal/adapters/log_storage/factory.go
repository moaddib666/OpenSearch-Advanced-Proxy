package log_storage

import (
	"OpenSearchAdvancedProxy/internal/adapters/log_provider"
	"OpenSearchAdvancedProxy/internal/adapters/search"
	"OpenSearchAdvancedProxy/internal/adapters/websockets"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	log "github.com/sirupsen/logrus"
)

type BaseStorageFactory struct {
	ctx context.Context
}

// NewBaseStorageFactory - create new BaseStorageFactory
func NewBaseStorageFactory(ctx context.Context) *BaseStorageFactory {
	return &BaseStorageFactory{
		ctx: ctx,
	}
}

func (b *BaseStorageFactory) FromConfig(name string, config *models.SubConfig) (ports.Storage, error) {
	if config.Version != models.ConfigVersion1 {
		return nil, models.ErrUnsupportedVersion
	}
	if config.Fields == nil {
		return nil, models.ErrNoFields
	}
	aggregatorFactory := search.NewAggregatorFactory()
	fields := &models.Fields{}
	for fieldName, field := range config.Fields {
		fields.AddField(fieldName, field)
	}

	if config.Provider == models.JsonLogFileProvider {
		logFile, ok := config.ProviderConfig["logfile"]
		if !ok {
			log.Errorf("No logfile specified for %s", name)
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
	}

	if config.Provider == models.WebSocketProvider {
		proto := search.NewDistributedJsonSearchProtocol()
		eventProcessor := NewEventProcessor(proto)
		server := websockets.NewWebSocketServer(config.ProviderConfig["bindAddress"], eventProcessor)
		go server.Run(b.ctx)
		return NewWebsocketServerStorage(name, fields, server, eventProcessor, proto, aggregatorFactory), nil
	}
	return nil, models.ErrUnsupportedProvider
}
