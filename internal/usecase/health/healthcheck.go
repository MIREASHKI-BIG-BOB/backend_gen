package health

import "backend_gen/internal/models/dto"

type healthUseCase struct{}

func NewHealthUseCase() HealthUseCase {
	return &healthUseCase{}
}

func (h *healthUseCase) CheckHealth() *dto.HealthResponse {
	return &dto.HealthResponse{
		Status: "OK",
	}
}
