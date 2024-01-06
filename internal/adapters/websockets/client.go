package websockets

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type WebsocketProxy struct {
	storages map[string]ports.Storage
	protocol ports.DistributedSearchProtocol
	dsn      string
	conn     *websocket.Conn
	done     chan struct{}
	mu       *sync.Mutex

	reconnectAfter time.Duration
}

// NewWebsocketProxy creates a new WebsocketProxy
func NewWebsocketProxy(dsn string, protocol ports.DistributedSearchProtocol) *WebsocketProxy {
	return &WebsocketProxy{
		storages:       make(map[string]ports.Storage),
		protocol:       protocol,
		done:           make(chan struct{}),
		dsn:            dsn,
		mu:             &sync.Mutex{},
		reconnectAfter: 3 * time.Second,
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
		if internalError != nil {
			log.Errorf("Error reading message: %s", internalError.Error())
			return err
		}
		searchRequest, internalError := w.protocol.UnmarshallSearchRequest(message)
		if internalError != nil {
			log.Errorf("Error unmarshalling search request: %s", internalError.Error())
			continue
		}
		var searchResult *models.DistributedSearchResult
		w.mu.Lock()
		storage, ok := w.storages[searchRequest.Index]
		w.mu.Unlock()
		if !ok {
			log.Errorf("Storage %s not found", searchRequest.Index)
			searchResult = models.DistributedSearchResultFailed(searchRequest.ID, false)
		} else {
			storageSearchResult, internalError := storage.Search(searchRequest.SearchRequest)
			if err != nil {
				log.Errorf("Error searching storage %s: %s", searchRequest.Index, internalError.Error())
				searchResult = models.DistributedSearchResultFailed(searchRequest.ID, false)
			} else {
				searchResult = models.DistributedSearchResultSuccess(searchRequest.ID, storageSearchResult)
			}
		}
		raw := w.protocol.MarshallSearchResult(searchResult)
		log.Debugf("Sending message: %s", searchRequest.ID)
		internalError = w.conn.WriteMessage(websocket.TextMessage, raw)
		if internalError != nil {
			log.Errorf("Failed to send message: %s", internalError.Error())
		}
	}
}

func (w *WebsocketProxy) Start(ctx context.Context) (err error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugf("Context is done, closing connection")
				err = w.conn.Close()
				return
			default:
				log.Debugf("Context is not done, establishing new connection")
				err = w.establishConnection()
				if err != nil {
					log.Errorf("Error on connection: %s", err.Error())
				}
			}
			log.Infof("reconnecting in %f seconds", w.reconnectAfter.Seconds())
			<-time.After(w.reconnectAfter)
		}
	}()
	return nil
}

func (w *WebsocketProxy) AddStorage(storage ports.Storage) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.storages[storage.Name()] = storage
}
