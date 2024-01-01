package main

import (
	"OpenSearchAdvancedProxy/internal/adapters/config"
	"OpenSearchAdvancedProxy/internal/adapters/log_storage"
	"OpenSearchAdvancedProxy/internal/adapters/search"
	"OpenSearchAdvancedProxy/internal/adapters/websockets"
	"context"
	log "github.com/sirupsen/logrus"
)

var ConfigDir = "tmp/client_config"

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	ctx := context.Background()
	dsn := "ws://localhost:8080/"
	cfg := config.NewConfig(ConfigDir)
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}

	protocol := search.NewDistributedJsonSearchProtocol()
	proxy := websockets.NewWebsocketProxy(dsn, protocol)
	storageFactory := log_storage.NewBaseStorageFactory(ctx)
	for indexName, logConfig := range cfg.AvailableIndexes() {
		storage, confError := storageFactory.FromConfig(indexName, logConfig)
		if confError != nil {
			log.Errorf("Error creating storage for %s: %s", indexName, confError)
			continue
		}
		proxy.AddStorage(storage)
	}
	err = proxy.Start(ctx)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	// wait forever
	select {}
}
