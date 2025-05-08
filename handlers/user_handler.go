package handlers

import (
	"encoding/json"
	"net/http"
	"payments_service/models"
	"payments_service/services"

	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Handler RegisterUser called")
	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	id, err := u.service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type Response struct {
		Success bool
		ID      int
	}

	successful := Response{Success: true, ID: id}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(successful)
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Handler LoginUser called")

	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := u.service.LoginUser(user)
	if err != nil {
		http.Error(w, token, http.StatusInternalServerError)
		return
	}

	type Response struct {
		Success bool
		Token   string
	}

	successful := Response{Success: true, Token: token}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(successful)
}
