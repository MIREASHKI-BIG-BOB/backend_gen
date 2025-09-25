package health

import (
	"backend_gen/internal/usecase"
	"backend_gen/pkg/http/writer"
	"net/http"
)

func NewHealthHandler(healthUC usecase.HealthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := healthUC.CheckHealth()

		if response.Status != "OK" {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		writer.WriteStatusOK(w)
		writer.WriteJson(w, response)
	}
}
