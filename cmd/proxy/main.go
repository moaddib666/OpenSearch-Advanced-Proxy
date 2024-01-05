package main

import (
	"OpenSearchAdvancedProxy/internal/adapters/config"
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy"
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy/handlers"
	"OpenSearchAdvancedProxy/internal/adapters/log_storage"
	"OpenSearchAdvancedProxy/internal/adapters/monitoring"
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

var ProxyAddr = "0.0.0.0:6600"
var MetricsAddr = "0.0.0.0:9002"
var OpenSearchAddr = "http://localhost:9200"
var ConfigDir = ".local/config"

func init() {
	log.SetLevel(log.InfoLevel)
	if url := os.Getenv("ELASTICSEARCH_URL"); url != "" {
		log.Debugf("Using ELASTICSEARCH_URL from environment: %s", url)
		OpenSearchAddr = url
	}
	metrics := monitoring.NewMetrics()
	metrics.Bind(MetricsAddr)
}

func main() {
	ctx := context.Background()
	cfg := config.NewConfig(ConfigDir)
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
	proxy := http_proxy.NewHttpProxy(ProxyAddr, OpenSearchAddr, handlers.DefaultHandler(OpenSearchAddr))
	storageFactory := log_storage.NewBaseStorageFactory(ctx)
	// TODO add composite
	for indexName, logConfig := range cfg.AvailableIndexes() {
		storage, confError := storageFactory.FromConfig(indexName, logConfig)
		if confError != nil {
			log.Errorf("Error creating storage for %s: %s", indexName, confError)
			continue
		}
		proxy.AddStorage(storage)
	}
	log.Infof("Starting proxy server on %s", ProxyAddr)
	err = proxy.Start(ctx)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
