package currency_service

// CurrencyClient — интерфейс для работы с внешними API
type CurrencyClient interface {
	GetExchangeRate(from, to string) (float32, error)
}
