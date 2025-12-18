package handlers

import (
	"html/template"
	"net/http"
)

func Handle_api(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/api.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, "api.html")
}

func Handle_login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, "login.html")
}
