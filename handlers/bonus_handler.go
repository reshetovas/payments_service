package handlers

import (
	"encoding/json"
	"net/http"
	"payments_service/models"
	"payments_service/services"

	"github.com/rs/zerolog/log"
)

type BonusHandler struct {
	service *services.BonusService
}

func NewBonusHandler(service *services.BonusService) *BonusHandler {
	return &BonusHandler{
		service: service,
	}
}

// method create
func (bh *BonusHandler) CreateBonus(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Handler CreateBonus called")

	var bonus models.Bonus
	err := json.NewDecoder(r.Body).Decode(&bonus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := bh.service.CreateBonus(bonus)
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

// method get
func (bh *BonusHandler) GetBonus(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Handler GetBonus called")

	bonuses, err := bh.service.GetBonus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(bonuses)
}

// method update
func (bh *BonusHandler) UpdateBonus(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Handler UpdateBonus called")
	var bonus models.Bonus
	err := json.NewDecoder(r.Body).Decode(&bonus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = bh.service.UpdateBonus(bonus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
