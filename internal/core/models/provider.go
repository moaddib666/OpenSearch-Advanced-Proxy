package models

type ProviderType string

const (
	JsonLogFileProvider ProviderType = "jsonLogFile"
	WebSocketProvider   ProviderType = "webSocketServer"
)
