package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // разрешаем все origin для простоты
	},
}

// Глобальный менеджер клиентов
var clientManager = &ClientManager{
	clients:    make(map[*websocket.Conn]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
}

type ClientManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
}

func (cm *ClientManager) run() {
	for {
		select {
		case client := <-cm.register:
			cm.mutex.Lock()
			cm.clients[client] = true
			cm.mutex.Unlock()
			log.Printf("Новый клиент подключен. Всего клиентов: %d", len(cm.clients))

		case client := <-cm.unregister:
			cm.mutex.Lock()
			if _, ok := cm.clients[client]; ok {
				delete(cm.clients, client)
				client.Close()
			}
			cm.mutex.Unlock()
			log.Printf("Клиент отключен. Всего клиентов: %d", len(cm.clients))

		case message := <-cm.broadcast:
			cm.mutex.RLock()
			for client := range cm.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Ошибка отправки сообщения клиенту: %v", err)
					client.Close()
					delete(cm.clients, client)
				}
			}
			cm.mutex.RUnlock()
		}
	}
}

func main() {
	r := chi.NewRouter()

	//для лог
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ws", handleWebSocket)

	// Запуск менеджера клиентов в отдельной горутине
	go clientManager.run()

	log.Println("WebSocket сервер запущен на :8081")
	log.Println("Готов принимать сообщения от генератора и транслировать их клиентам на ws://localhost:8081/ws")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении до WebSocket: %v", err)
		return
	}

	// Регистрируем клиента
	clientManager.register <- conn

	defer func() {
		// Отменяем регистрацию при отключении
		clientManager.unregister <- conn
	}()

	// Обработка входящих сообщений
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket ошибка: %v", err)
			}
			break
		}

		// Выводим сообщение в консоль (от генератора)
		log.Printf("📊 Получено от генератора: %s", string(message))

		// Транслируем сообщение всем подключенным клиентам (фронтенду)
		clientManager.broadcast <- message
	}
}
