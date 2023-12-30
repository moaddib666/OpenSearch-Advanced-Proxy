package main

import (
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy"
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy/handlers"
	"OpenSearchAdvancedProxy/internal/adapters/mock_storage"
	"context"
	log "github.com/sirupsen/logrus"
)

var ProxyAddr = "localhost:6600"
var OpensearchAdde = "localhost:9200"

func init() {
	log.SetLevel(log.DebugLevel)

}

func main() {
	ctx := context.Background()
	proxy := http_proxy.NewHttpProxy(ProxyAddr, OpensearchAdde, handlers.DefaultHandler(OpensearchAdde))
	log.Infof("Starting proxy server on %s", ProxyAddr)
	proxy.AddStorage(mock_storage.NewMockStorage())
	err := proxy.Start(ctx)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
