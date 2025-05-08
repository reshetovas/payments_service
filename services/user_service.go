package services

import (
	"payments_service/models"
	"payments_service/storage"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	storage storage.UserStorageActions
	token   *TokenStruct
}

func NewUserService(storage storage.UserStorageActions, token *TokenStruct) *UserService {
	return &UserService{
		storage: storage,
		token:   token,
	}
}

func (s *UserService) CreateUser(user models.User) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.Password = string(hashedPassword)

	return s.storage.CreateUser(user)
}

func (s *UserService) LoginUser(user_request models.User) (string, error) {

	user_storage, err := s.storage.GetUserByName(user_request)
	if err != nil {
		return "user not found", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user_storage.Password), []byte(user_request.Password)); err != nil {
		return "invalid credentials", err
	}

	token, err := s.token.GenerateJWT(user_storage.ID, user_storage.Username)
	if err != nil {
		return "token generation error", err
	}

	return token, nil
}
