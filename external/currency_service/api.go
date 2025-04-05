package currency_service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CurrencyAPI — реальная реализация работы с внешним API
type CurrencyAPI struct {
	APIURL string
}

// Новый клиент для работы с API
func NewCurrencyAPI(apiURL string) *CurrencyAPI {
	return &CurrencyAPI{APIURL: apiURL}
}

// Получить курс валют с реального API
func (c *CurrencyAPI) GetExchangeRate(from, to string) (float32, error) {
	apiURL := fmt.Sprintf("%s/latest/%s", c.APIURL, from)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	rates, ok := result["rates"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to get exchange rates")
	}

	rate, ok := rates[to].(float64)
	if !ok {
		return 0, fmt.Errorf("currency %s not found", to)
	}

	return float32(rate), nil
}
