package tests

import "fmt"

// MockCurrencyAPI — мок для тестов
type MockCurrencyAPI struct{}

// Получить курс валют из мока
func (m *MockCurrencyAPI) GetExchangeRate(from, to string) (float32, error) {
	if from == "USD" && to == "EUR" {
		return 0.85, nil
	} else if from == "USD" && to == "GBP" {
		return 0.75, nil
	}
	return 0, fmt.Errorf("currency not found")
}
