package services

import (
	"errors"
	"fmt"
	"time"

	"payments_service/internal/ws"
	"payments_service/storage"

	"payments_service/models"

	"github.com/rs/zerolog/log"
)

type CurrencyClient interface {
	GetExchangeRate(from, to string) (float32, error)
}

type PaymentService struct {
	storage         storage.PaymentStorageActions
	currencyService CurrencyClient
	hub             *ws.Hub
}

// function for creating object (ekzemplyar)
/*StorageActions - interface => to use NewPaymentService struct shoul include interface function
here NewPaymentService must have input struct with inetrface StorageActions*/
func NewPaymentService(storage storage.PaymentStorageActions, currencyService CurrencyClient, hub *ws.Hub) *PaymentService {
	return &PaymentService{
		storage:         storage,
		currencyService: currencyService,
		hub:             hub,
	}
}

// method create
func (s *PaymentService) CreatePayment(payment models.Payment) (int, error) {
	log.Info().Msg("CreatePayment called in service")
	// if payment.ID == 0 {
	// 	return 0, errors.New("id is required")
	// }

	//fill in atribete a
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = time.Now()
	}

	storagePayment := models.Payment{
		ID:          payment.ID,
		UserID:      payment.UserID,
		Amount:      payment.Amount,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
	}

	id, err := s.storage.CreatePayment(storagePayment)
	if err != nil {
		return 0, err
	}

	s.hub.SendToUser(storagePayment.UserID, []byte(fmt.Sprintf("New payment with id: %d", storagePayment.ID)))
	return id, err
}

// method get
func (s *PaymentService) GetPayments() ([]models.Payment, error) {
	log.Info().Msg("GetPayments called in service")
	storagePayments, err := s.storage.GetPayments()
	if err != nil {
		return nil, err
	}

	// var payments []Payment
	// for _, p := range storagePayments {
	// 	payments = append(payments, Payment{
	// 		ID:          p.ID,
	// 		Amount:      p.Amount,
	// 		Description: p.Description,
	// 		CreatedAt:   p.CreatedAt,
	// 		State:       p.State,
	// 		Items:       p.Items,
	// 	})
	// }

	return storagePayments, nil
}

// method update
func (s *PaymentService) UpdatePayment(payment models.Payment) error {
	log.Info().Msg("UpdatePayment called in service")
	if payment.ID == 0 {
		return errors.New("id is required")
	}

	storagePayment := models.Payment{
		ID:          payment.ID,
		Amount:      payment.Amount,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
	}

	return s.storage.UpdatePayment(storagePayment)
}

func (s *PaymentService) PatchPayment(id int, updates map[string]interface{}) error {
	log.Info().Msg("UpdatePayment called in service")
	if id == 0 {
		return errors.New("id is required")
	}

	return s.storage.PartialUpdatePayment(id, updates)
}

// method delete
func (s *PaymentService) DeletePayment(id int) error {
	log.Info().Msg("DeletePayment called in service")
	if id == 0 {
		return errors.New("id is required")
	}

	return s.storage.DeletePayment(id)
}

func (s *PaymentService) GetPaymentInCurrency(id int, currency string) (models.Payment, error) {
	payment, err := s.storage.GetPaymentByID(id)
	if err != nil {
		return models.Payment{}, err
	}

	rate, err := s.currencyService.GetExchangeRate("USD", currency)
	if err != nil {
		return models.Payment{}, err
	}

	convertedAmount := payment.Amount * rate
	return models.Payment{
		ID:          payment.ID,
		Amount:      convertedAmount,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
	}, nil
}

func (s *PaymentService) PaymentClose(id int) error {
	log.Info().Msg("UpdatePayment called in service")
	if id == 0 {
		return errors.New("id is required")
	}

	updates := map[string]interface{}{
		"state": "close",
	}

	return s.storage.PartialUpdatePayment(id, updates)
}
