package ws

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"time"
)

type Hub struct {
	clients    map[int]*Client // userID -> client
	register   chan *Client
	unregister chan *Client
	sendTo     chan Message
}

type Message struct {
	UserID  int
	Payload []byte
}

const (
	writeWait  = 10 * time.Second
	pingPeriod = (writeWait * 9) / 10
)

type Client struct {
	userID int
	conn   *websocket.Conn
	send   chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sendTo:     make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
			log.Info().Msgf("Client %d registered", client.userID)
			log.Info().Msgf("Client %+v connected", h.clients[client.userID])
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				log.Info().Msgf("Client %d unregistered", client.userID)
			}
		case message := <-h.sendTo:
			if client, ok := h.clients[message.UserID]; ok {
				select {
				case client.send <- message.Payload:
					log.Info().Msgf("SendToUser %d", message.UserID)
				default:
					close(client.send)
					delete(h.clients, client.userID)
					log.Info().Msgf("Client %d unregistered", client.userID)
				}
			}
		}
	}
}

func (h *Hub) Register(userID int, conn *websocket.Conn) *Client {
	client := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 1024),
	}
	h.register <- client
	log.Info().Msgf("Client %d registered in chanel", userID)

	return client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
	log.Info().Msgf("Client %d unregistered in chanel", client.userID)
}

func (h *Hub) SendToUser(userID int, payload []byte) {
	h.sendTo <- Message{
		UserID:  userID,
		Payload: payload,
	}
	log.Info().Msgf("SendToUser in chanel %d", userID)
	log.Info().Msgf("SendToUser payload: %+v", h.clients[userID])
}

// Читает входящие (если нужно) — можно оставить пустым
func (c *Client) ReadPump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			log.Warn().Err(err).Msg("ReadPump error")
			break
		}
	}
}

// Пишет в сокет всё, что приходит в канал Send
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Error().Err(err).Msg("WriteMessage failed")
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
