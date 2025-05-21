package services

import (
	"errors"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenStruct struct {
	secret []byte
}

func NewTokenStruct(secret []byte) *TokenStruct {
	return &TokenStruct{
		secret: secret,
	}
}

type Claims struct {
	User_id int
	jwt.RegisteredClaims
}

func (t TokenStruct) GenerateJWT(userID int, userName string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"exp":       time.Now().Add(24 * time.Hour).Unix(), //TTL - 1 day
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //set the signature type
	return token.SignedString(t.secret)
}

func (t TokenStruct) ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	})
	log.Debug().Msgf("token: %+v", token)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		log.Debug().Msgf("claims: %+v", claims)
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
