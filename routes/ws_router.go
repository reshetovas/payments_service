package routes

import (
	"payments_service/handlers"
	"payments_service/middleware"

	"github.com/gorilla/mux"
)

type WSRoutes struct {
	WSHandler     *handlers.WSHandler
	authorization *middleware.AuthMiddleware
}

func NewWSRoutes(WSHandler *handlers.WSHandler, authorization *middleware.AuthMiddleware) *WSRoutes {
	return &WSRoutes{
		WSHandler:     WSHandler,
		authorization: authorization,
	}
}

func (w *WSRoutes) WSRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(w.authorization.JWTAuthMiddleware)

	r.HandleFunc("/ws", w.WSHandler.WebSocketHandler())

	return r
}
