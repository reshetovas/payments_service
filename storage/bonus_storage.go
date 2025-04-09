package storage

import (
	"database/sql"
	"errors"

	"payments_service/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type BonusStorage struct {
	dbProperty *sql.DB
}

func NewBonusStorage(db *sql.DB) *BonusStorage {
	return &BonusStorage{
		dbProperty: db,
	}
}

type BonusStorageActions interface {
	GetBonuses() ([]models.Bonus, error)
	CreateBonus(payment models.Bonus) (int, error)
	UpdateBonus(payment models.Bonus) error
}

func (bs *BonusStorage) GetBonuses() ([]models.Bonus, error) {
	//query to db
	log.Info().Msg("GetBonuses called in storage")
	rows, err := bs.dbProperty.Query("SELECT id, payment_id, amount FROM bonuses")
	if err != nil {
		return nil, err
	}
	defer rows.Close() //to distonnect the connection to the db

	//read each database entry and fills the object
	var bonuses []models.Bonus
	for rows.Next() {
		var b models.Bonus
		err := rows.Scan(&b.ID, &b.PaymentID, &b.Amount)
		if err != nil {
			return nil, err
		}

		bonuses = append(bonuses, b)
	}
	return bonuses, nil
}

// method create
func (bs *BonusStorage) CreateBonus(bonus models.Bonus) (int, error) {
	log.Info().Msg("CreateBonus called in storage")

	query := `INSERT INTO bonuses (id, payment_id, amount) VALUES (?, ?, ?)`
	result, err := bs.dbProperty.Exec(query, bonus.ID, bonus.PaymentID, bonus.Amount)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("Bonus created, id: %d", id)
	return int(id), nil
}

// method update put
func (bs *BonusStorage) UpdateBonus(bonus models.Bonus) error {
	log.Info().Msg("UpdateBonus called in storage")
	result, err := bs.dbProperty.Exec(
		"UPDATE payments SET payment_id = ?, amount = ? WHERE id = ?",
		bonus.PaymentID, bonus.Amount, bonus.ID,
	)
	if err != nil {
		return err
	}

	// checking the update result (for cases where id not found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("bonus not found")
	}
	return nil
}
