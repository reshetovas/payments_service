package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"payments_service/models"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type PaymentCache struct {
	storage       PaymentStorageActions
	cacheType     string
	redis         *redis.Client
	InMemoryCashe InMemoryCashe
}

type InMemoryCashe struct {
	data  map[string]string
	mutex sync.RWMutex //for blocking all struct
}

func NewPaymentCache(storage PaymentStorageActions, cacheType string, redis *redis.Client) *PaymentCache {
	return &PaymentCache{
		storage:   storage,
		cacheType: cacheType,
		redis:     redis,
	}
}

func (c *PaymentCache) GetPayments() ([]models.Payment, error) {
	return c.storage.GetPayments()
}

// method create
func (c *PaymentCache) CreatePayment(payment models.Payment) (int, error) {
	return c.storage.CreatePayment(payment)
}

// method update put
func (c *PaymentCache) UpdatePayment(payment models.Payment) error {
	key := fmt.Sprintf("payment:%d", payment.ID)

	switch strings.ToLower(c.cacheType) {
	case "redis":
		ctx := context.Background()
		c.setInRedis(ctx, payment, key)
	case "memory":
		c.InMemoryCashe.setInMemory(payment, key)
	}

	return c.storage.UpdatePayment(payment)
}

func (c *PaymentCache) PartialUpdatePayment(id int, updates map[string]interface{}) error {
	key := fmt.Sprintf("payment:%d", id)

	switch strings.ToLower(c.cacheType) {
	case "redis":
		ctx := context.Background()
		err := c.redis.Del(ctx, key).Err()
		if err != nil {
			log.Error().Err(err).Msgf("Error redis.Del(ctx, %s)", key)
		}
	case "memory":
		delete(c.InMemoryCashe.data, key)
	}

	return c.storage.PartialUpdatePayment(id, updates)
}

// method delete
func (c *PaymentCache) DeletePayment(id int) error {
	key := fmt.Sprintf("payment:%d", id)

	switch strings.ToLower(c.cacheType) {
	case "redis":
		ctx := context.Background()
		err := c.redis.Del(ctx, key).Err()
		if err != nil {
			log.Error().Err(err).Msgf("Error redis.Del(ctx, %s)", key)
		}
	case "memory":
		delete(c.InMemoryCashe.data, key)
	}

	return c.storage.DeletePayment(id)
}

// method get by id
func (c *PaymentCache) GetPaymentByID(id int) (models.Payment, error) {
	log.Info().Msg("Cache GetPaymentByID called")
	key := fmt.Sprintf("payment:%d", id)
	ctx := context.Background()

	//serch in cache
	switch strings.ToLower(c.cacheType) {
	case "redis":
		payment, err := c.getInRedis(ctx, key)
		if err == nil {
			return payment, err
		}
	case "memory":
		log.Debug().Msg("Get Case = memory, is done")
		payment, err := c.InMemoryCashe.getInMemory(key)
		if err == nil {
			return payment, err
		}
	}

	//search in DB
	p, err := c.storage.GetPaymentByID(id)
	if err != nil {
		return models.Payment{}, err
	}

	switch strings.ToLower(c.cacheType) {
	case "redis":
		c.setInRedis(ctx, p, key)
	case "memory":
		log.Debug().Msg("Set Case = memory, is done")
		c.InMemoryCashe.setInMemory(p, key)
	}
	fmt.Println(c.InMemoryCashe.data)
	return p, nil
}

func (c *PaymentCache) CreateItem(item models.Item) error {
	return c.storage.CreateItem(item)
}

func (c *PaymentCache) GetItemsByPaymentID(paymentID int) ([]models.Item, error) {
	return c.storage.GetItemsByPaymentID(paymentID)
}

// pyments NOT IN ('CLOSED', 'FAILED')
func (c *PaymentCache) GetPendingPayments() ([]models.Payment, error) {
	return c.storage.GetPendingPayments()
}

// /
func (i *InMemoryCashe) getInMemory(key string) (models.Payment, error) {
	log.Info().Msg("getInMemory called")
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	val, ok := i.data[key]
	if !ok {
		return models.Payment{}, errors.New("not found")
	}

	var p models.Payment
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		log.Debug().Msg("Get payment InMemory: unmarshal error")
		return models.Payment{}, errors.New("unmarshal error")
	}

	return p, nil
}

func (i *InMemoryCashe) setInMemory(payment models.Payment, key string) {
	log.Info().Msg("getInMemory called")
	i.mutex.Lock()
	defer i.mutex.RUnlock()
	value, _ := json.Marshal(payment)

	i.data[key] = string(value)
}

func (c *PaymentCache) getInRedis(ctx context.Context, key string) (models.Payment, error) {
	log.Info().Msg("getInRedis called")
	cached, err := c.redis.Get(ctx, key).Result()
	if err == nil {
		var p models.Payment
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			return p, nil
		}
		return models.Payment{}, err
	}
	return models.Payment{}, err
}

func (c *PaymentCache) setInRedis(ctx context.Context, payment models.Payment, key string) {
	log.Info().Msg("getInRedis called")
	data, _ := json.Marshal(payment)
	if err := c.redis.Set(ctx, key, data, 10*24*time.Hour).Err(); err != nil {
		log.Error().Err(err).Msg("setInRedis error")
	}
}
