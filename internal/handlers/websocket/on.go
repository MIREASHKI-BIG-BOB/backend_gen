package websocket

import (
	"backend_gen/internal/usecase"
	httpErr "backend_gen/pkg/http/error"
	"fmt"
	"net/http"
)

func OnSocket(uc usecase.WebSocketUseCase, sensorID string, sensorToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("ws://localhost:8080/ws/sensor?sensor_id=%s", sensorID)
		err := uc.Connect(url, sensorToken)
		if err != nil {
			httpErr.InternalError(w, fmt.Errorf("failed to connect: %w", err))
			return
		}

		// Запускаем отправку сообщений каждую секунду
		err = uc.StartSendingMessages()
		if err != nil {
			httpErr.InternalError(w, fmt.Errorf("failed to start sending messages: %w", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
