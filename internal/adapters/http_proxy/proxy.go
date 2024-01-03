package http_proxy

import (
	"OpenSearchAdvancedProxy/internal/adapters/http_proxy/handlers"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

type HttpProxy struct {
	*http.ServeMux
	srcAddr        string
	destAddr       string
	defaultHandler http.HandlerFunc
}

func (h *HttpProxy) Start(ctx context.Context) error {
	h.ServeMux.HandleFunc("/", h.defaultHandler)

	srv := &http.Server{
		Addr:    h.srcAddr,
		Handler: h.ServeMux,
	}

	go func() {
		<-ctx.Done()          // wait for context cancellation
		_ = srv.Shutdown(ctx) // shutdown the server when context is done
	}()

	return srv.ListenAndServe() // start the server
}

func (h *HttpProxy) AddStorage(storage ports.Storage) {
	name := storage.Name()
	log.Infof("Adding storage %s::%s", reflect.TypeOf(storage).Elem().Name(), name)
	log.Debugf("Registering handlers for storage: `%s`", name)
	indexHandlerName := "/_resolve/index/" + name
	log.Debugf("Registering handler: `%s`", indexHandlerName)
	h.ServeMux.HandleFunc(indexHandlerName, handlers.IndexHandler(name))

	fieldCapsHandlerName := "/" + name + "/_field_caps"
	log.Debugf("Registering handler: `%s`", fieldCapsHandlerName)
	h.ServeMux.HandleFunc(fieldCapsHandlerName, handlers.FieldCapsHandler(storage.Fields()))
	searchHandlerName := "/" + name + "/_search"
	log.Debugf("Registering handler: `%s`", searchHandlerName)
	h.ServeMux.HandleFunc(searchHandlerName, handlers.SearchHandler(storage))
}

// NewHttpProxy creates a new HttpProxy.
func NewHttpProxy(src, dest string, defaultHandler http.HandlerFunc) *HttpProxy {
	return &HttpProxy{
		ServeMux:       http.NewServeMux(),
		srcAddr:        src,
		destAddr:       dest,
		defaultHandler: defaultHandler,
	}
}
