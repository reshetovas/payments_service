package tests

import (
	"payments_service/services"
	"testing"
)

func TestGetPayments(t *testing.T) {
	mockStorage := &MockStorage{}
	service := services.NewPaymentService(mockStorage, &MockCurrencyAPI{})

	payments, err := service.GetPayments()
	if err != nil {
		t.Errorf("we have error: %v", err)
	}

	if len(payments) != 2 {
		t.Errorf("Expected 2 payments, but got: %v", err)
	}
}

func TestCreatePayment(t *testing.T) {
	mockStorage := &MockStorage{}
	service := services.NewPaymentService(mockStorage, &MockCurrencyAPI{})

	id, err := service.CreatePayment(services.Payment{ID: 3, Amount: 150, Description: "pending"})
	if err != nil {
		t.Errorf("Payment creation error: %v", err)
	}

	if id != 3 {
		t.Errorf("Expected ID=3 and got: %v", err)
	}
}

func TestUpdatePayment(t *testing.T) {
	mockStorage := &MockStorage{}
	service := services.NewPaymentService(mockStorage, &MockCurrencyAPI{})

	err := service.UpdatePayment(services.Payment{ID: 1, Amount: 150, Description: "mock"})
	if err != nil {
		t.Errorf("Payment update error: %v", err)
	}

	err = service.UpdatePayment(services.Payment{ID: 99, Amount: 300, Description: "mock"})
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDeletePayment(t *testing.T) {
	mockStorage := &MockStorage{}
	service := services.NewPaymentService(mockStorage, &MockCurrencyAPI{})

	err := service.DeletePayment(1)
	if err != nil {
		t.Errorf("Payment deletion error: %v", err)
	}

	err = service.DeletePayment(99)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestGetPaymentInCurrency(t *testing.T) {
	mockStorage := &MockStorage{}
	service := services.NewPaymentService(mockStorage, &MockCurrencyAPI{})

	payment, err := service.GetPaymentInCurrency(1, "EUR")
	if err != nil {
		t.Errorf("Error while getting payment in currency: %v", err)
	}

	// Проверяем, что сумма была конвертирована правильно
	expectedAmount := float32(100) * 0.85 // 100 USD * 0.85 (EUR)
	if payment.Amount != expectedAmount {
		t.Errorf("Expected amount %v, but got %v", expectedAmount, payment.Amount)
	}

	// Проверяем, что ID платежа корректный
	if payment.ID != 1 {
		t.Errorf("Expected ID=1, but got %v", payment.ID)
	}

}
