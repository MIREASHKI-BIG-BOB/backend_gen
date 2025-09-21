package health

import "backend_gen/internal/models/dto"

type HealthUseCase interface {
	CheckHealth() *dto.HealthResponse
}
