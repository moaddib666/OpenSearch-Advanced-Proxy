package log_storage

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Request struct {
	id             string
	searchRequest  *models.SearchRequest
	callback       chan *models.SearchResult
	answerCount    int
	answerRequires int
}

// IsReady returns true if the request is ready to be answered
func (r *Request) IsReady() bool {
	return r.answerCount >= r.answerRequires
}

type EventProcessor struct {
	protocol ports.DistributedSearchProtocol
	requests map[string]*Request
	mux      *sync.Mutex
}

func NewEventProcessor(proto ports.DistributedSearchProtocol) *EventProcessor {
	return &EventProcessor{
		protocol: proto,
		requests: make(map[string]*Request),
		mux:      &sync.Mutex{},
	}
}

// MakeRequest creates a new request and returns its ID
func (e *EventProcessor) MakeRequest(searchRequest *models.SearchRequest) (string, <-chan *models.SearchResult) {
	id := uuid.NewString()
	callback := make(chan *models.SearchResult)
	e.mux.Lock()
	e.requests[id] = &Request{
		id:             id,
		searchRequest:  searchRequest,
		callback:       callback,
		answerCount:    0,
		answerRequires: 9999, // FIXME: raise condition if answer were already received and it already equals to answerRequires
	}
	e.mux.Unlock()
	return id, callback
}

// deleteRequest deletes request from the map
func (e *EventProcessor) deleteRequest(request *Request) {
	e.mux.Lock()
	close(request.callback)
	delete(e.requests, request.id)
	e.mux.Unlock()
}

// AnswerRequest sends the search result to the client
func (e *EventProcessor) AnswerRequest(id string, result *models.SearchResult) {
	request, ok := e.requests[id]
	if !ok {
		return
	}
	request.callback <- result
	request.answerCount++
	if request.IsReady() {
		e.deleteRequest(request)
	}
}

func (e *EventProcessor) ResponseExpected(requestId string, answerRequires int) {
	e.mux.Lock()
	request, ok := e.requests[requestId]
	e.mux.Unlock()
	if !ok {
		return
	}
	request.answerRequires = answerRequires
	if request.IsReady() {
		e.deleteRequest(request)
	}
	// FIXME: raise condition if answer were already received and it already equals to answerRequires
}

func (e *EventProcessor) OnIncomingMessage(client ports.WebsocketServerClient, message []byte) error {
	result, err := e.protocol.UnmarshallSearchResult(message)
	if err != nil {
		return err
	}
	e.AnswerRequest(result.ID, result.SearchResult)
	return nil
}

type WebsocketServerStorage struct {
	name              string
	fields            *models.Fields
	server            ports.WebsocketServer
	processor         ports.DistributedRequestsProcessor
	protocol          ports.DistributedSearchProtocol
	aggregatorFactory ports.SearchAggregatorFactory
}

// NewWebsocketServerStorage creates a new WebsocketServerStorage struct
func NewWebsocketServerStorage(name string, fields *models.Fields, server ports.WebsocketServer, processor ports.DistributedRequestsProcessor, proto ports.DistributedSearchProtocol, aggregator ports.SearchAggregatorFactory) *WebsocketServerStorage {
	// TODO create composite for args
	return &WebsocketServerStorage{
		name:              name,
		fields:            fields,
		server:            server,
		processor:         processor,
		protocol:          proto,
		aggregatorFactory: aggregator,
	}
}

func (w *WebsocketServerStorage) Name() string {
	return w.name
}

func (w *WebsocketServerStorage) Fields() *models.Fields {
	return w.fields
}

func (w *WebsocketServerStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	id, callback := w.processor.MakeRequest(r)
	rawRequest := w.protocol.MarshallSearchRequest(&models.DistributedSearchRequest{
		ID:            id,
		SearchRequest: r,
		Index:         w.name,
	})
	awaitCount := w.server.Broadcast(rawRequest)
	w.processor.ResponseExpected(id, awaitCount)
	aggregate := w.aggregatorFactory.CreateAggregator(r)
	for result := range callback {
		log.Debugf("Got result from %+v", result)
		aggregate.AddResult(result)
	}
	// TODO context with timeout for waiting responses 5 minutes
	return aggregate.GetResult(), nil
}
