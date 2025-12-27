package jwttokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("mmy_sseccrrett_kkeyy")

func GenerateAccessToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		//24 часа - срок годности
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		//24 часа - срок годности
	} 

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
}
