package websocket

type Client interface {
	Connect(url string, token string) error
	Disconnect() error
	IsConnected() bool
	SendMessage(message []byte) error
}
