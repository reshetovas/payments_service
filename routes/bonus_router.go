package routes

import (
	"payments_service/handlers"
	"payments_service/middleware"

	"github.com/gorilla/mux"
)

type BonusRoutes struct {
	handler *handlers.BonusHandler
}

func NewBonusRoutes(handler *handlers.BonusHandler) *BonusRoutes {
	return &BonusRoutes{
		handler: handler,
	}
}

func (bp *BonusRoutes) BonusRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/bonus", bp.handler.CreateBonus).Methods("POST")
	r.HandleFunc("/bonuses", bp.handler.GetBonus).Methods("GET")
	r.HandleFunc("/bonus", bp.handler.UpdateBonus).Methods("PUT") // Исправлено с UPDATE на PUT
	r.HandleFunc("/bonus/{id}", bp.handler.GetBonusByID).Methods("GET")

	return r
}
