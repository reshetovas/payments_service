package routes

import (
	"payments_service/handlers"
	"payments_service/middleware"

	"github.com/gorilla/mux"
)

type PaymentRoutes struct {
	handler       *handlers.PaymentHandler
	authorization *middleware.AuthMiddleware
}

// function for creating object (ekzemplyar)
func NewPaymentRoutes(handler *handlers.PaymentHandler, auth *middleware.AuthMiddleware) *PaymentRoutes {
	return &PaymentRoutes{
		handler:       handler,
		authorization: auth,
	}
}

func (p *PaymentRoutes) PaymentRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(p.authorization.JWTAuthMiddleware)

	r.HandleFunc("/payment", p.handler.CreatePayment).Methods("POST")
	r.HandleFunc("/payments", p.handler.GetPayments).Methods("GET")
	r.HandleFunc("/payment", p.handler.UpdatePayment).Methods("PUT") // Исправлено с UPDATE на PUT
	r.HandleFunc("/payment", p.handler.PatchPayment).Methods("PATCH")
	r.HandleFunc("/payment/{id}", p.handler.DeletePayment).Methods("DELETE")
	r.HandleFunc("/payment/{id}/inCurrency", p.handler.GetPaymentInCurrency).Methods("GET")
	r.HandleFunc("/payment/{id}/close", p.handler.PaymentClose).Methods("POST")

	return r
}
