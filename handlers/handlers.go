package handlers

import (
	"html/template"
	"net/http"
	myerrors "zadanieweb/errors"
	"zadanieweb/users"

	"github.com/google/uuid"
)

// Главная страница
func Handle_api(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/api.html")
	if err != nil {
		panic(err)
	}

	userUUIDstring := users.AuthMap[users.AuthKey]
	if userUUIDstring != "" {
		cur_user := users.User{}

		userUUID, err := uuid.Parse(userUUIDstring)
		if err != nil {
			panic(err)
		}

		cur_user.FillUserWithUUID(userUUID)

		t.Execute(w, cur_user)
	} else {
		t.Execute(w, users.EmptyUser())
	}

}

// Обработчик страницы входа в аккаунт
func Handle_getLogin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getlogin.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, myerrors.Cur_error)
	myerrors.Cur_error.ErrMsg = ""
}

// Обработчик страницы входа в аккаунт
func Handle_getRegister(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getregister.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, myerrors.Cur_error)
	myerrors.Cur_error.ErrMsg = ""
}

// Обработчик страницы с постами
func Handle_posts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		panic(err)
	}

}
