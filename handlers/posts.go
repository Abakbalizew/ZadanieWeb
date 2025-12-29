// posts
package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"zadanieweb/databases"
	myerrors "zadanieweb/errors"
	"zadanieweb/jwttokens"
	"zadanieweb/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Обрабатывается, когда пользователь пытается войти в учётную запись (нажимает "Войти")
func Handle_postLogin(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			//Если ошибка, то выводим её на экран.
			myerrors.Cur_error.ErrMsg = "Неверный логин или пароль! :("
			http.Redirect(w, r, "/api/auth/getlogin", http.StatusSeeOther)

			fmt.Println("Перехвачено исключение:", r)
		}
	}()

	email := r.FormValue("email")
	enteredPassword := r.FormValue("password")

	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка подключения бд :( ")
		panic(err)
	}
	defer db.Close()

	var userUUID uuid.UUID
	var hashedPassword []byte
	err = db.QueryRow("SELECT userid, passwordhash FROM users WHERE email = $1",
		email).Scan(&userUUID, &hashedPassword)
	if err != nil {
		fmt.Printf("Такого пользователя нет! :( ")
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(enteredPassword))
	if err != nil {
		myerrors.Cur_error.ErrMsg = "Неверный пароль"
		http.Redirect(w, r, "/api/auth/getlogin", http.StatusSeeOther)
		return
	}
	//Генерируем токены и отправляем в бд
	access_token, err := jwttokens.GenerateAccessToken(userUUID)
	if err != nil {
		fmt.Printf("Ошибка генерации токена :( ")
		panic(err)
	}
	refresh_token, err := jwttokens.GenerateRefreshToken(userUUID)
	if err != nil {
		fmt.Printf("Ошибка генерации токена :( ")
		panic(err)
	}
	_, err = db.Exec("UPDATE users SET accesstoken = $1, refreshtoken = $2 WHERE userid = $3", access_token, refresh_token, userUUID)
	if err != nil {
		fmt.Printf("Не удалось обновить токены :( ")
		panic(err)
	}

	//Передаем в состояние из main.go айди аутентифицированного пользователя
	models.AuthMap[models.AuthKey] = userUUID.String()

	//Возвращаем на главную страницу
	http.Redirect(w, r, "/api", http.StatusSeeOther)
}

// Обрабатывается, когда пользователь пытается зарегистрироваться (нажимает "Зарегистрироваться")
func Handle_postRegister(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			//Если ошибка, то выводим её на экран.

			http.Redirect(w, r, "/api/auth/getregister", http.StatusSeeOther)

			fmt.Println("Перехвачено исключение:", r)
		}
	}()

	//Дефолтная ошибка
	myerrors.Cur_error.ErrMsg = "Что-то пошло не так! :("

	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка подключения бд :( ")
		panic(err)
	}
	defer db.Close()

	users, err := db.Exec("SELECT FROM users WHERE email = $1", email)
	if err != nil {
		fmt.Printf("Ошибка поиска совпадений по бд :( ")
		panic(err)
	}
	//Проверка, что почта не занята
	count, err := users.RowsAffected()
	if err != nil {
		fmt.Printf("Ошибка c bd.RowsAffected() :( ")
		panic(err)
	}
	if count != 0 {
		myerrors.Cur_error.ErrMsg = "Почта уже зарегистрирована"
		http.Redirect(w, r, "/api/auth/getregister", http.StatusSeeOther)
		return
	}

	//Проверка, что пароль содержит символы
	if len(strings.TrimSpace(password)) == 0 {
		myerrors.Cur_error.ErrMsg = "Пароль должен содержать хотя бы 1 символ"
		http.Redirect(w, r, "/api/auth/getregister", http.StatusSeeOther)
		return
	}

	//Хэшируем пароль пользователя
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Ошибка хэширования пароля :( ")
		panic(err)
	}

	//Генерируем uuid пользователя
	userUUID := uuid.New()

	//Генерируем токены
	access_token, err := jwttokens.GenerateAccessToken(userUUID)
	if err != nil {
		fmt.Printf("Ошибка генерации токена :( ")
		panic(err)
	}
	refresh_token, err := jwttokens.GenerateRefreshToken(userUUID)
	if err != nil {
		fmt.Printf("Ошибка генерации токена :( ")
		panic(err)
	}

	_, err = db.Exec("INSERT INTO"+
		" users (userid, email, passwordhash, role_, refreshtoken, accesstoken)"+
		" VALUES ($1, $2, $3, $4, $5, $6)",
		userUUID, email, passwordHash, role, refresh_token, access_token)
	if err != nil {
		fmt.Printf("Ошибка сохранения значений в бд :( ")
		panic(err)
	}

	//Передаем в состояние из main.go айди аутентифицированного пользователя
	models.AuthMap[models.AuthKey] = userUUID.String()

	http.Redirect(w, r, "/api", http.StatusSeeOther)
}

func Handle_exit(w http.ResponseWriter, r *http.Request) {
	models.AuthMap[models.AuthKey] = ""

	http.Redirect(w, r, "/api", http.StatusSeeOther)
}
