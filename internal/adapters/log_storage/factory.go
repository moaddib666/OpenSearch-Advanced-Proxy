package log_storage

import (
	"OpenSearchAdvancedProxy/internal/adapters/indexer"
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
	fields := &models.Fields{}
	for fieldName, field := range config.Fields {
		fields.AddField(fieldName, field)
	}

	return b.createStorage(name, config.Provider, fields, config.Timestamp.Field)
}

// createStorage - create a storage from a config
func (b *BaseStorageFactory) createStorage(name string, cfg ports.ProviderConfig, fields *models.Fields, timestampField string) (ports.Storage, error) {

	aggregatorFactory := search.NewAggregatorFactory()

	switch cfg.GetProvider() {
	case models.JsonLogFileProvider:
		config := &models.JsonLogFileProviderConfig{}
		err := cfg.GetProviderConfig(config)
		if err != nil {
			return nil, err
		}
		if config.LogFile == "" {
			return nil, models.ErrNoLogFile
		}
		var idx ports.Indexer
		if config.Index != nil {
			idx = indexer.NewJsonFileIndexer(config.LogFile, timestampField, config.Index.Resolution)
			err := idx.LoadOrCreateIndex()
			if err != nil {
				log.Fatalf("Error loading index: %s", err.Error())
			}
		}
		provider := log_provider.NewLogFileProvider(config.LogFile,
			func() ports.LogEntry {
				return &log_provider.JsonLogEntry{
					TimeStampField: timestampField,
				}
			},
			idx,
		)
		engine := search.NewLogSearchEngine(provider)
		return NewGenericStorage(name, fields, engine), nil
	case models.WebSocketProvider:
		config := &models.WebSocketProviderConfig{}
		err := cfg.GetProviderConfig(config)
		if err != nil {
			return nil, err
		}
		if config.BindAddress == "" {
			return nil, models.ErrNoBindAddress
		}
		proto := search.NewDistributedJsonSearchProtocol()
		eventProcessor := NewEventProcessor(proto)
		server := websockets.NewWebSocketServer(config.BindAddress, eventProcessor)
		go server.Run(b.ctx)
		return NewWebsocketServerStorage(name, fields, server, eventProcessor, proto, aggregatorFactory), nil
	case models.ClickhouseProvider:
		config := &models.ClickhouseProviderConfig{}
		err := cfg.GetProviderConfig(config)
		if err != nil {
			return nil, err
		}
		if config.DSN == "" {
			return nil, models.ErrNoClickhouseDSN
		}
		var searchableFields []string
		for fieldName := range fields.Fields {
			// Note: currently support of flag isSearchable
			searchableFields = append(searchableFields, fieldName)
		}
		factory := search.NewSQLQueryBuilderFactory(searchableFields, timestampField)
		provider := log_provider.NewClickhouseProvider(config.DSN, config.Table, factory, func() ports.LogEntry {
			return log_provider.SqlLogEntryConstructor()
		})
		engine := search.NewSQLDBSearchEngine(provider)
		return NewGenericStorage(name, fields, engine), nil
	case models.AggregateProvider:
		config := &models.AggregateProviderConfig{}
		err := cfg.GetProviderConfig(config)
		if err != nil {
			return nil, err
		}
		storages := make([]ports.Storage, len(config.SubProviders))
		for i, subConfig := range config.SubProviders {
			storage, err := b.createStorage(name, subConfig, fields, timestampField)
			if err != nil {
				return nil, err
			}
			storages[i] = storage
		}
		return NewAggregateStorage(name, storages, fields, aggregatorFactory), nil
	}
	return nil, models.ErrUnsupportedProvider
}
