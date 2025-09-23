package health

import (
	"backend_gen/internal/models/dto"
	"backend_gen/internal/usecase"
)

type healthUseCase struct{}

func NewHealthUseCase() usecase.HealthUseCase {
	return &healthUseCase{}
}

func (h *healthUseCase) CheckHealth() *dto.HealthResponse {
	return &dto.HealthResponse{
		Status: "OK",
	}
}
