package jwttokens

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"zadanieweb/databases"
	"zadanieweb/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("mmy_sseccrrett_kkeyy")

// Генерирует access-токен
func GenerateAccessToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
		//2 часа - срок годности
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Генерирует refresh-токен
func GenerateRefreshToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		//Неделя - срок годности
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

// MiddleWare, который активирует переданный обработчик, если
// токены валидны. Токены берутся из данных пользователя в бд, а
// к данным пользователя мы получаем доступ через uuid в models/users/AuthMap[AuthKey]
func CheckTokenMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Наш пользователь
		cur_user := models.User{}

		//UUID пользователя в виде строки
		userUUIDstring := models.AuthMap[models.AuthKey]
		if userUUIDstring != "" {
			cur_user = models.User{}

			//Переводим строку в тип uuid.UUID
			userUUID, err := uuid.Parse(userUUIDstring)
			if err != nil {
				panic(err)
			}
			//Заполняем нашего пользователя данными из бд
			cur_user.FillUserWithUUID(userUUID)
		} else {
			http.Redirect(w, r, "/api/auth/getlogin", http.StatusSeeOther)
			return
		}

		accessToken_string := cur_user.AccessToken
		accessToken, err := ParseToken(accessToken_string)

		//Если токен не валиден
		if !accessToken.Valid || err != nil {
			refreshToken_string := cur_user.RefreshToken
			refreshToken, err := ParseToken(refreshToken_string)
			//Если refrest-токен тоже не валиден, тогда обрываем работу функции, не вызвав next(w, r)
			if !refreshToken.Valid || err != nil {
				fmt.Printf("Токен не валиден! ")
				http.Redirect(w, r, "/api/auth/getlogin", http.StatusSeeOther)
				return

				//Иначе создаем новый access-токен
			} else {
				new_access_token, err := GenerateAccessToken(cur_user.UserUUID)
				if err != nil {
					panic(err)
				}
				//Отправляем новый токен в бд, вызываем next(w, r)
				db, err := sql.Open("postgres", databases.ConnStr)
				if err != nil {
					fmt.Printf("Ошибка: jwttokens/jwttokens.go/CheckTokenMiddleWare - 1")
					panic(err)
				}
				defer db.Close()

				_, err = db.Exec("UPDATE users SET accesstoken = $1 WHERE userid = $2", new_access_token, cur_user.UserUUID)
				if err != nil {
					fmt.Printf("Ошибка: jwttokens/jwttokens.go/CheckTokenMiddleWare - 2")
					panic(err)
				}

				next(w, r)
			}
			//Если же access-токен валиден
		} else {
			next(w, r)
		}

	}

}
