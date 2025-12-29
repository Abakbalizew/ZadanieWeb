package handlers

import (
	"html/template"
	"net/http"
	myerrors "zadanieweb/errors"
	// myerrors "zadanieweb/errors"
)

// Главная страница
func Handle_api(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/api.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, "api.html")
}

// Обработчик страницы входа в аккаунт
func Handle_getLogin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getlogin.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, myerrors.Cur_error)
}

// Обработчик страницы входа в аккаунт
func Handle_getRegister(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getregister.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, myerrors.Cur_error)
}