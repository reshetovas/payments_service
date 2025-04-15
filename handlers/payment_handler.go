package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"payments_service/models"
	"payments_service/services"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type PaymentHandler struct {
	service *services.PaymentService
	parse   *services.ParseService
}

// function for creating object (ekzemplyar)
func NewPaymentHandler(service *services.PaymentService, parse *services.ParseService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		parse:   parse,
	}
}

// method create
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("method:", r.Method).
		Str("endpoint", "/payment").
		Msg("Processing payment request")

	var payment models.Payment
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := h.service.CreatePayment(payment)
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
func (h *PaymentHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("method:", r.Method).
		Str("endpoint", "/payment").
		Msg("Processing payment request")

	payments, err := h.service.GetPayments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// method update
func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("UpdatePayment called")
	var payment models.Payment
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.UpdatePayment(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *PaymentHandler) PatchPayment(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("PatchPayment called")
	// ID from path
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Error().Err(err).Msg("Error decoding JSON r.Body")
		return
	}

	err = h.service.PatchPayment(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// method delete
func (h *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("DeletePayment called")
	idStr := r.URL.Path[len("/payments/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.DeletePayment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// method
func (h *PaymentHandler) GetPaymentInCurrency(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GetPaymentInCurrency")
	// ID from path
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Currency from query
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		http.Error(w, "currency is required", http.StatusBadRequest)
		return
	}

	// Получаем платеж в указанной валюте
	payment, err := h.service.GetPaymentInCurrency(id, currency)
	if err != nil {
		log.Error().Err(err).Msg("Error GetPaymentInCurrency")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) ProcessExistingFiles(ctx context.Context, dir string) {

	log.Info().Msg("start procces existing files")

	processedFiles := make(map[string]struct{})

	// GAVE NAME loop "LOOP"
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			log.Error().Err(err).Msg("Reading dir error:")
			return
		}
		filesCount := len(files)
		if filesCount == 0 {
			continue
		}
		//log.Println("files count", filesCount)

		for _, file := range files {
			if !file.IsDir() {
				_, ok := processedFiles[file.Name()]
				if ok {
					continue
				}
				log.Info().Str("file detected:", file.Name())
				go h.parse.ParsePaymentsFile(dir + "/" + file.Name())
				processedFiles[file.Name()] = struct{}{}

			}
		}
	}

	log.Info().Msg("finish procces existing files")
}

func (h *PaymentHandler) PaymentClose(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("PatchPayment called")
	// ID from path
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.PaymentClose(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// func (h *PaymentHandler) DirListeningParser(ctx context.Context, dirName string) {
// 	//create file system whatcher
// 	whatcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer whatcher.Close() //at the end of the program, we are guaranteed to close whatcher

// 	//Add(path string) error - add directory or file in whatching list
// 	err = whatcher.Add(dirName)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Whatching for changes in", dirName)

// 	h.processExistingFiles(ctx, dirName)

// }

// //start a go routine (bachgound process) that handles events
// go func() {
// 	log.Println("start scanning")
// 	for {
// 		//select - for async lisentening the same channels
// 		select {
// 		//read channel with events
// 		case event, ok := <-whatcher.Events: //whatcher.Events - events channel
// 			if !ok {
// 				fmt.Println("channel is closed")
// 				return //if the channel is closed, exit the go routine
// 			}
// 			fmt.Println("got event")
// 			if event.Has(fsnotify.Create) {
// 				fmt.Println("New file detected:", event.Name)

// 				go h.parse.ParsePaymentsFile(event.Name)
// 			}

// 		//read channel with errors
// 		case err, ok := <-whatcher.Errors: //whatcher.Errors - events channel
// 			if !ok {
// 				return
// 			}
// 			fmt.Println("Error", err)
// 		}
// 	}
// }()
