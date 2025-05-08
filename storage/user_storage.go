package storage

import (
	"database/sql"
	"errors"

	"payments_service/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type UserStorage struct {
	dbProperty *sql.DB
}

// function for creating object (ekzemplyar)
// output PaymentStorage which implements the interface
func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		dbProperty: db,
	}
}

type UserStorageActions interface {
	GetUserByName(user models.User) (models.User, error)
	CreateUser(user models.User) (int, error)
}

// method get
func (s *UserStorage) GetUserByName(user models.User) (models.User, error) {
	//query to db
	log.Info().Msg("GetUser called in storage")
	row := s.dbProperty.QueryRow("SELECT id, username, password, created_at FROM users where username = ?", user.Username)

	u := models.User{}

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return u, nil
}

// method create
func (s *UserStorage) CreateUser(user models.User) (int, error) {
	log.Info().Msg("CreateUser called in storage")

	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	result, err := s.dbProperty.Exec(query, user.Username, user.Password)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username " {
			return 0, errors.New("user with this username already exists")
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("User created, id: %d", id)
	return int(id), nil
}
