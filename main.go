package main

import (
	"net/http"
	"zadanieweb/handlers"

	"github.com/gorilla/mux"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api", handlers.Handle_api).Methods("GET")
	rtr.HandleFunc("/api/auth/login", handlers.Handle_login).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)

}
