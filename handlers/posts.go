// posts
package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"zadanieweb/databases"
	myerrors "zadanieweb/errors"
	"zadanieweb/jwttokens"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Обрабатывается, когда пользователь пытается войти в учётную запись (нажимает "войти")
func Handle_postLogin(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				myerrors.Cur_error.ErrMsg = "Неверный логин или пароль! :("
				http.Redirect(w, r, "/api/auth/getlogin", http.StatusSeeOther)

				fmt.Println("Перехвачено исключение:", r)
			}
		}()

		email := r.FormValue("email")
		password := r.FormValue("password")

		db, err := sql.Open("postgres", databases.ConnStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		var userId uuid.UUID
		err = db.QueryRow("SELECT FROM users WHERE email = $1 AND passwordhash = $2", email, password).Scan(&userId)
		if err != nil {
			panic(err)
		}

		//Генерим токен и отправляем в бд
		access_token, err := jwttokens.GenerateAccessToken(userId)
		if err != nil {
			panic(err)
		}
		refresh_token, err := jwttokens.GenerateRefreshToken(userId)
		if err != nil {
			panic(err)
		}
		insert, err := db.Query("INSERT INTO tokens (userid, access_token, refresh_token) VALUES ($1, $2, $3)", userId, access_token, refresh_token)
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		//Передаем в контекст из main.go айди аутентифицированного пользователя
		newCtx := context.WithValue(context.Background(), "AuthorizedUserId", userId)
		ctx = newCtx

		//Возвращаем на главную страницу
		http.Redirect(w, r, "/api", http.StatusSeeOther)
	}
}
