package usecase

import "backend_gen/internal/models/dto"

type WebSocketUseCase interface {
	Connect(url string) error
	Disconnect() error
	// TODO: remove this pizdes
	SendMessage(message any) error
	StartSendingMessages() error
	StopSendingMessages()
}

type HealthUseCase interface {
	CheckHealth() *dto.HealthResponse
}
