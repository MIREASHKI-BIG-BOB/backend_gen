package websocket

import (
	"backend_gen/internal/usecase"
	httpErr "backend_gen/pkg/http/error"
	"fmt"
	"net/http"
)

func OffSocket(uc usecase.WebSocketUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := uc.Disconnect()
		if err != nil {
			httpErr.InternalError(w, fmt.Errorf("failed to disconnect: %w", err))
			return
		}
		
		w.WriteHeader(http.StatusOK)
	}
}
