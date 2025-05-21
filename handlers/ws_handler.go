package handlers

import (
	"net/http"
	"payments_service/ctxutils"
	"payments_service/internal/ws"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WSHandler) WebSocketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Извлечение userID из контекста
		userID, ok := ctxutils.GetUserID(r.Context())
		log.Info().Msgf("Client %d identified", userID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Апгрейд соединения в WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		/*
			1. Check headers: Upgrade, connection, Sec-WebSocket
			2. Send http-response 101
			3. w.(http.Hidjacker).Hijack()
				- return net.Conn (TCP connection)
				- disables the standart operation of http.ResponseWriter
				- make you owner socket
		*/
		if err != nil {
			log.Error().Err(err).Msg("failed to upgrade connection")
			http.Error(w, "Failed to upgrade", http.StatusInternalServerError)
			return
		}

		// Регистрация пользователя в хабе
		client := ws.hub.Register(userID, conn)

		go client.ReadPump()
		go client.WritePump()
	}
}

// websocat -H="Authorization: Bearer " ws://localhost:8080/ws
