package usecase

import "backend_gen/internal/models/dto"

type WebSocketUseCase interface {
	Connect(url string) error
	Disconnect() error
	SendMessage(message any) error
}

type HealthUseCase interface {
	CheckHealth() *dto.HealthResponse
}
