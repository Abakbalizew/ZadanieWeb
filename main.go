package main

import (
	"net/http"
	"zadanieweb/handlers"
	"zadanieweb/models"

	"github.com/gorilla/mux"
)

func main() {
	//Обнуляем ошибки, чтобы они не выводились на экранах регистрации и входа в аккаунт просто так
	models.AuthMap[models.AuthKey] = ""

	//Router чтобы сделать динамические страницы
	rtr := mux.NewRouter()

	//Главная страница
	rtr.HandleFunc("/api", handlers.Handle_api).Methods("GET")
	//Страница входа в аккаунт
	rtr.HandleFunc("/api/auth/getlogin", handlers.Handle_getLogin).Methods("GET")
	//Страница регистрации
	rtr.HandleFunc("/api/auth/getregister", handlers.Handle_getRegister).Methods("GET")
	//Выход из аккаунта
	rtr.HandleFunc("/api/auth/exit", handlers.Handle_exit).Methods("GET")

	//Посты
	rtr.HandleFunc("/api/posts", handlers.Handle_posts).Methods("GET")

	rtr.HandleFunc("/api/auth/login", handlers.Handle_postLogin).Methods("POST")
	rtr.HandleFunc("/api/auth/register", handlers.Handle_postRegister).Methods("POST")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
