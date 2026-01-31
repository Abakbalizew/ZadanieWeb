package jwttokens

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		//Обычный access-токен (2 часа жизни)
		{"token is valid", generateTokenWithTime(uuid.New(), time.Hour*2), "Valid"},
		//Просроченный на час токен
		{"token is not valid", generateTokenWithTime(uuid.New(), -time.Hour), "Not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseToken(tt.token)

			var ans string

			if err != nil {
				ans = "Not valid"
			} else {
				ans = "Valid"
			}

			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

// Создаёт токен с переданным сроком жизни.
// Это нужно, чтобы создавать просроченные токены
func generateTokenWithTime(userId uuid.UUID, life_time time.Duration) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(life_time).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_str, err := token.SignedString(secretKey)
	if err != nil {
		return ""
	}
	return token_str
}
