package ports

import (
	"context"
	"net/http"
)

type WebsocketServerClient interface {
	SendMessage(message []byte) error
	Listen()
}

type WebsocketServer interface {
	Run(ctx context.Context)
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
	Broadcast(message []byte) int
}

type WebsocketServerEventProcessor interface {
	OnIncomingMessage(client WebsocketServerClient, message []byte) error
}
