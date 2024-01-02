package websockets

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
)

type WebsocketProxy struct {
	storages map[string]ports.Storage
	protocol ports.DistributedSearchProtocol
	dsn      string
	conn     *websocket.Conn
	done     chan struct{}
	mu       *sync.Mutex
}

// NewWebsocketProxy creates a new WebsocketProxy
func NewWebsocketProxy(dsn string, protocol ports.DistributedSearchProtocol) *WebsocketProxy {
	return &WebsocketProxy{
		storages: make(map[string]ports.Storage),
		protocol: protocol,
		done:     make(chan struct{}),
		dsn:      dsn,
		mu:       &sync.Mutex{},
	}
}

func (w *WebsocketProxy) establishConnection() (err error) {
	log.Infof("Establishing connection to %s", w.dsn)
	w.conn, _, err = websocket.DefaultDialer.Dial(w.dsn, nil)
	if err != nil {
		return err
	}
	defer w.conn.Close()
	for {
		_, message, internalError := w.conn.ReadMessage()
		log.Debugf("Received message: %s", string(message))
		if internalError != nil {
			log.Errorf("Error reading message: %s", err.Error())
			return err
		}
		searchRequest, internalError := w.protocol.UnmarshallSearchRequest(message)
		if internalError != nil {
			log.Errorf("Error unmarshalling search request: %s", err.Error())
			continue
		}
		searchResult := &models.DistributedSearchResult{
			ID: searchRequest.ID,
		}
		w.mu.Lock()
		storage, ok := w.storages[searchRequest.Index]
		w.mu.Unlock()
		if !ok {
			log.Errorf("Storage %s not found", searchRequest.Index)
		} else {
			storageSearchResult, err := storage.Search(searchRequest.SearchRequest)
			if err != nil {
				log.Errorf("Error searching storage %s: %s", searchRequest.Index, err.Error())
			} else {
				searchResult = &models.DistributedSearchResult{
					ID:           searchRequest.ID,
					SearchResult: storageSearchResult,
				}
			}
		}
		raw := w.protocol.MarshallSearchResult(searchResult)
		log.Debugf("Sending message: %s", searchRequest.ID)
		internalError = w.conn.WriteMessage(websocket.TextMessage, raw)
		if internalError != nil {
			log.Errorf("Failed to send message: %s", err.Error())
		}
	}
}

func (w *WebsocketProxy) Start(ctx context.Context) (err error) {
	go func() {
		select {
		case <-ctx.Done():
			log.Debugf("Context is done, closing connection")
			err = w.conn.Close()
			return
		default:
			err = w.establishConnection()
		}
	}()
	return nil
}

func (w *WebsocketProxy) AddStorage(storage ports.Storage) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.storages[storage.Name()] = storage
}
