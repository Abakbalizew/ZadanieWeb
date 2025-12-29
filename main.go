package main

import (
	"context"
	"net/http"
	"zadanieweb/handlers"

	"github.com/gorilla/mux"
)

func main() {
	//Контекст, где мы будем хранить id юзера, если он аутентифицирован
	ctx := context.Background()
	//Router чтобы сделать динамические страницы
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api", handlers.Handle_api).Methods("GET")
	rtr.HandleFunc("/api/auth/getlogin", handlers.Handle_getLogin).Methods("GET")
	rtr.HandleFunc("/api/auth/getregister", handlers.Handle_getRegister).Methods("GET")

	rtr.HandleFunc("/api/auth/login", handlers.Handle_postLogin(ctx)).Methods("POST")
	rtr.HandleFunc("/api/auth/register", handlers.Handle_postRegister(ctx)).Methods("POST")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
