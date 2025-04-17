package services

import (
	"errors"
	"payments_service/models"
	"payments_service/storage"

	"github.com/rs/zerolog/log"
)

type BonusService struct {
	storage storage.BonusStorageActions
}

func NewBonusService(storage storage.BonusStorageActions) *BonusService {
	return &BonusService{
		storage: storage,
	}
}

func (bs *BonusService) CreateBonus(bonus models.Bonus) (int, error) {
	log.Info().Msg("CreateBonus called in service")
	if bonus.ID == 0 {
		return 0, errors.New("id is required")
	}

	storageBonus := models.Bonus{
		ID:        bonus.ID,
		PaymentID: bonus.PaymentID,
		Amount:    bonus.Amount,
	}

	return bs.storage.CreateBonus(storageBonus)
}

// method get
func (bs *BonusService) GetBonus() ([]models.Bonus, error) {
	log.Info().Msg("GetBonus called in service")
	storagePayments, err := bs.storage.GetBonuses()
	if err != nil {
		return nil, err
	}

	return storagePayments, nil
}

// method update
func (bs *BonusService) UpdateBonus(bonus models.Bonus) error {
	log.Info().Msg("UpdatePayment called in service")
	if bonus.ID == 0 {
		return errors.New("id is required")
	}

	storageBonus := models.Bonus{
		ID:        bonus.ID,
		PaymentID: bonus.PaymentID,
		Amount:    bonus.Amount,
	}

	return bs.storage.UpdateBonus(storageBonus)
}

func (bs *BonusService) GetBonusByID(id int) (models.Bonus, error) {
	bonus, err := bs.storage.GetBonusByID(id)
	if err != nil {
		return models.Bonus{}, err
	}

	return bonus, nil
}
