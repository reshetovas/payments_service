package routes

import (
	"payments_service/handlers"
	"payments_service/middleware"

	"github.com/gorilla/mux"
)

type UserRoutes struct {
	handler *handlers.UserHandler
}

func NewUserRoutes(handler *handlers.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

func (u *UserRoutes) UserRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/users/register", u.handler.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", u.handler.LoginUser).Methods("POST")

	return r
}
