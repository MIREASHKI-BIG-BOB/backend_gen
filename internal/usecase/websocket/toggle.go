package websocket

import (
	"backend_gen/internal/ports/websocket"
	"backend_gen/internal/usecase"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type WebSocketUseCase struct {
	client websocket.Client
	ticker *time.Ticker
	stopCh chan bool
}

func (uc *WebSocketUseCase) Connect(url string) error {
	if uc.client.IsConnected() {
		return fmt.Errorf("already connected")
	}

	err := uc.client.Connect(url)
	if err != nil {
		return err
	}

	return nil
}

func (uc *WebSocketUseCase) Disconnect() error {
	if !uc.client.IsConnected() {
		slog.Warn("WebSocket not connected")
		return fmt.Errorf("not connected")
	}

	err := uc.client.Disconnect()
	if err != nil {
		slog.Error("Failed to disconnect WebSocket", "error", err)
		return err
	}

	slog.Info("WebSocket disconnected successfully")
	return nil
}

func (uc *WebSocketUseCase) SendMessage(message any) error {
	if !uc.client.IsConnected() {
		return fmt.Errorf("not connected")
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		slog.Error("Failed to marshal message to JSON", "error", err)
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return uc.client.SendMessage(jsonData)
}

func (uc *WebSocketUseCase) StartSendingMessages() error {
	if !uc.client.IsConnected() {
		return fmt.Errorf("not connected")
	}
	if uc.ticker != nil {
		uc.StopSendingMessages()
	}

	slog.Info("Starting periodic message sending", "interval", "1s")

	uc.ticker = time.NewTicker(1 * time.Second)
	uc.stopCh = make(chan bool)

	go func() {
		for {
			select {
			case <-uc.ticker.C:
				message := websocket.MessageData{
					Sex: "yes",
				}
				if err := uc.SendMessage(message); err != nil {
					slog.Error("Failed to send periodic JSON message", "error", err)
				}
			case <-uc.stopCh:
				slog.Info("Stopping periodic message sending")
				return
			}
		}
	}()

	return nil
}

func (uc *WebSocketUseCase) StopSendingMessages() {
	if uc.ticker != nil {
		uc.ticker.Stop()
		uc.ticker = nil
	}
	if uc.stopCh != nil {
		close(uc.stopCh)
		uc.stopCh = nil
	}
}

func NewWebSocketUseCase(client websocket.Client) usecase.WebSocketUseCase {
	return &WebSocketUseCase{
		client: client,
	}
}
