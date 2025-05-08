package ws

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type Hub struct {
	mu          sync.RWMutex
	connections map[int]*websocket.Conn // UserID -> list of connections
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int]*websocket.Conn),
	}
}

func (h *Hub) Register(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[userID] = conn

	log.Info().Msgf("User %d connected from WebSocket", userID)
}

func (h *Hub) Unregister(userID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Закрываем соединение и удаляем его из мапы
	if conn, ok := h.connections[userID]; ok {
		err := conn.Close()
		if err != nil {
			log.Error().Msgf("Error closing connection for user %d: %v", userID, err)
		}
		delete(h.connections, userID)
		log.Info().Msgf("User %d unregistered from WebSocket", userID)
	}
}

func (h *Hub) SendToUser(userID int, message []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Проверяем, существует ли соединение для пользователя
	if conn, ok := h.connections[userID]; ok {
		return conn.WriteMessage(websocket.TextMessage, message)
	}
	return nil
}
