package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // разрешаем все origin для простоты
	},
}

func main() {
	r := chi.NewRouter()

	//для лог
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ws", handleWebSocket)

	log.Println("WebSocket сервер запущен на :8081")
	log.Println("Готов принимать сообщения от других сервисов на ws://localhost:8081/ws")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении до WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Новый WebSocket клиент подключен")

	// зулеплено
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			break
		}

		// Выводим сообщение в консоль
		log.Printf("� Получено сообщение: %s", string(message))
	}

	log.Println("WebSocket клиент отключен")
}
