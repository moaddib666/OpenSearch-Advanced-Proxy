package main

import (
	"OpenSearchAdvancedProxy/internal/adapters/config"
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy"
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy/handlers"
	"OpenSearchAdvancedProxy/internal/adapters/log_storage"
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

var ProxyAddr = "0.0.0.0:6600"
var OpenSearchAddr = "http://localhost:9200"
var ConfigDir = "tmp/config"

func init() {
	log.SetLevel(log.DebugLevel)
	// Get ELASTICSEARCH_URL from environment
	if url := os.Getenv("ELASTICSEARCH_URL"); url != "" {
		log.Debugf("Using ELASTICSEARCH_URL from environment: %s", url)
		OpenSearchAddr = url
	}
}
func main() {
	ctx := context.Background()
	cfg := config.NewConfig(ConfigDir)
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
	proxy := http_proxy.NewHttpProxy(ProxyAddr, OpenSearchAddr, handlers.DefaultHandler(OpenSearchAddr))
	storageFactory := log_storage.NewBaseStorageFactory()
	// TODO add composite
	for indexName, logConfig := range cfg.AvailableIndexes() {
		storage, confError := storageFactory.FromConfig(indexName, logConfig)
		if confError != nil {
			log.Errorf("Error creating storage for %s: %s", indexName, confError)
			continue
		}
		proxy.AddStorage(storage)
	}
	//proxy.AddStorage(log_storage.NewMockStorage())
	log.Infof("Starting proxy server on %s", ProxyAddr)
	err = proxy.Start(ctx)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
