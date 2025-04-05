package tests

import (
	"errors"
	"time"

	"payments_service/storage"
)

type MockStorage struct{}

func (m *MockStorage) GetPayments() ([]storage.Payment, error) {
	return []storage.Payment{
		{ID: 1, Amount: 100, Description: "test_1", CreatedAt: time.Now()},
		{ID: 2, Amount: 200, Description: "test_2", CreatedAt: time.Now()},
	}, nil
}

func (m *MockStorage) CreatePayment(payment storage.Payment) (int, error) {
	return 3, nil
}

func (m *MockStorage) UpdatePayment(payment storage.Payment) error {
	if payment.ID == 1 {
		return nil
	}
	return errors.New("payment not found")
}

func (m *MockStorage) DeletePayment(id int) error {
	if id == 1 {
		return nil
	}
	return errors.New("payment not found")
}

func (m *MockStorage) GetPaymentByID(id int) (storage.Payment, error) {
	if id == 1 {
		return storage.Payment{ID: 1, Amount: 100, Description: "test_1", CreatedAt: time.Now()}, nil
	}
	return storage.Payment{}, errors.New("payment not found")
}

// bad
func (m *MockStorage) CreateItem(item storage.Item) error {
	return nil
}
