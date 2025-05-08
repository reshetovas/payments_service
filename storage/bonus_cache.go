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

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type BonusCache struct {
	storage   BonusStorageActions
	cacheType string
	redis     *redis.Client
	data      map[string]string
	mutex     sync.RWMutex //for blocking all struct
}

func NewBonusCache(storage BonusStorageActions, cacheType string, redis *redis.Client) *BonusCache {
	return &BonusCache{
		storage:   storage,
		cacheType: cacheType,
		redis:     redis,
		data:      make(map[string]string),
	}
}

func (bc *BonusCache) GetBonuses() ([]models.Bonus, error) {
	return bc.storage.GetBonuses()
}

// method create
func (bc *BonusCache) CreateBonus(bonus models.Bonus) (int, error) {
	return bc.storage.CreateBonus(bonus)
}

// method update put
func (bc *BonusCache) UpdateBonus(bonus models.Bonus) error {
	key := fmt.Sprintf("payment:%d", bonus.ID)

	switch strings.ToLower(bc.cacheType) {
	case "redis":
		ctx := context.Background()
		bc.setInRedis(ctx, bonus, key)
	case "memory":
		bc.setInMemory(bonus, key)
	}

	return bc.storage.UpdateBonus(bonus)
}

func (bc *BonusCache) GetBonusByID(id int) (models.Bonus, error) {
	log.Info().Msg("Cache GetBonusByID called")
	key := fmt.Sprintf("payment:%d", id)
	ctx := context.Background()

	//serch in cache
	switch strings.ToLower(bc.cacheType) {
	case "redis":
		bonus, err := bc.getInRedis(ctx, key)
		if err == nil {
			return bonus, err
		}
	case "memory":
		payment, err := bc.getInMemory(key)
		if err == nil {
			return payment, err
		}
	}

	//search in DB
	p, err := bc.storage.GetBonusByID(id)
	if err != nil {
		return models.Bonus{}, err
	}

	switch strings.ToLower(bc.cacheType) {
	case "redis":
		bc.setInRedis(ctx, p, key)
	case "memory":
		log.Debug().Msg("Set Case = memory, is done")
		bc.setInMemory(p, key)
	}
	log.Info().Msgf("memory: %v", bc.data)
	return p, nil
}

func (bc *BonusCache) getInMemory(key string) (models.Bonus, error) {
	log.Info().Msg("getInMemory called")
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if bc.data == nil {
		return models.Bonus{}, errors.New("not found")
	}

	val, ok := bc.data[key]
	if !ok {
		return models.Bonus{}, errors.New("not found")
	}

	var p models.Bonus
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		log.Debug().Msg("Get payment InMemory: unmarshal error")
		return models.Bonus{}, errors.New("unmarshal error")
	}

	return p, nil
}

func (bc *BonusCache) setInMemory(payment models.Bonus, key string) {
	log.Info().Msg("setInMemory called")
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	value, _ := json.Marshal(payment)

	bc.data[key] = string(value)
}

func (bc *BonusCache) getInRedis(ctx context.Context, key string) (models.Bonus, error) {
	log.Info().Msg("getInRedis called")
	cached, err := bc.redis.Get(ctx, key).Result()
	if err == nil {
		var p models.Bonus
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			return p, nil
		}
		return models.Bonus{}, err
	}
	return models.Bonus{}, err
}

func (bc *BonusCache) setInRedis(ctx context.Context, payment models.Bonus, key string) {
	log.Info().Msg("getInRedis called")
	data, _ := json.Marshal(payment)
	if err := bc.redis.Set(ctx, key, data, 10*24*time.Hour).Err(); err != nil {
		log.Error().Err(err).Msg("setInRedis error")
	}
}
