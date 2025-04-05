package services

import (
	"errors"
	"fmt"
	"time"

	"payments_service/storage"
)

type Payment struct {
	ID          int
	Amount      float32
	Description string
	CreatedAt   time.Time
	Currency    string
	ShopID      int
	Address     string
	State       string
	Items       []Item
}

type CurrencyClient interface {
	GetExchangeRate(from, to string) (float32, error)
}

type PaymentService struct {
	storage         storage.StorageActions
	currencyService CurrencyClient
}

// function for creating object (ekzemplyar)
/*StorageActions - interface => to use NewPaymentService struct shoul include interface function
here NewPaymentService must have input struct with inetrface StorageActions*/
func NewPaymentService(storage storage.StorageActions, currencyService CurrencyClient) *PaymentService {
	return &PaymentService{
		storage:         storage,
		currencyService: currencyService,
	}
}

// method create
func (s *PaymentService) CreatePayment(payment Payment) (int, error) {
	fmt.Println("CreatePayment called in service")
	if payment.ID == 0 {
		return 0, errors.New("id is required")
	}

	//fill in atribete a
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = time.Now()
	}

	storagePayment := storage.Payment{
		ID:          payment.ID,
		Amount:      payment.Amount,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
	}

	return s.storage.CreatePayment(storagePayment)
}

// method get
func (s *PaymentService) GetPayments() ([]storage.Payment, error) {
	fmt.Println("GetPayments called in service")
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
func (s *PaymentService) UpdatePayment(payment Payment) error {
	fmt.Println("UpdatePayment called in service")
	if payment.ID == 0 {
		return errors.New("id is required")
	}

	storagePayment := storage.Payment{
		ID:          payment.ID,
		Amount:      payment.Amount,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
	}

	return s.storage.UpdatePayment(storagePayment)
}

func (s *PaymentService) PatchPayment(id int, updates map[string]interface{}) error {
	fmt.Println("UpdatePayment called in service")
	if id == 0 {
		return errors.New("id is required")
	}

	return s.storage.PartialUpdatePayment(id, updates)
}

// method delete
func (s *PaymentService) DeletePayment(id int) error {
	fmt.Println("DeletePayment called in service")
	if id == 0 {
		return errors.New("id is required")
	}

	return s.storage.DeletePayment(id)
}

func (s *PaymentService) GetPaymentInCurrency(id int, currency string) (Payment, error) {
	payment, err := s.storage.GetPaymentByID(id)
	if err != nil {
		return Payment{}, err
	}

	rate, err := s.currencyService.GetExchangeRate("USD", currency)
	if err != nil {
		return Payment{}, err
	}

	convertedAmount := payment.Amount * rate
	return Payment{
		ID:          payment.ID,
		Amount:      convertedAmount,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
	}, nil
}

func (s *PaymentService) PaymentClose(id int) error {
	fmt.Println("UpdatePayment called in service")
	if id == 0 {
		return errors.New("id is required")
	}

	updates := map[string]interface{}{
		"state": "close",
	}

	return s.storage.PartialUpdatePayment(id, updates)
}
