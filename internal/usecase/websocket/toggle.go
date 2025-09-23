package websocket

import (
	"backend_gen/internal/ports/websocket"
	"backend_gen/internal/usecase"
	"fmt"
)

type websocketUseCase struct {
	client websocket.Client
}

func (uc *websocketUseCase) Connect(url string) error {
	if uc.client.IsConnected() {
		return fmt.Errorf("already connected")
	}

	return uc.client.Connect(url)
}

func (uc *websocketUseCase) Disconnect() error {
	if !uc.client.IsConnected() {
		return fmt.Errorf("not connected")
	}

	return uc.client.Disconnect()
}

func (uc *websocketUseCase) SendMessage(message any) error {
	if !uc.client.IsConnected() {
		return fmt.Errorf("not connected")
	}

	// TODO: implement message sending when needed
	return fmt.Errorf("send message not implemented yet")
}

func NewWebSocketUseCase(client websocket.Client) usecase.WebSocketUseCase {
	return &websocketUseCase{
		client: client,
	}
}
