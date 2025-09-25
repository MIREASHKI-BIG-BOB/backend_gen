package websocket

type Client interface {
	Connect(url string) error
	Disconnect() error
	IsConnected() bool
	SendMessage(message []byte) error
}
