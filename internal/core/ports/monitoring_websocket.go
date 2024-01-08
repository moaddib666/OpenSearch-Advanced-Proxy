package ports

type WebsocketServerMonitor interface {
	RegisterClient(client WebsocketServerClient)
	UnregisterClient(client WebsocketServerClient)
	Init()
}
