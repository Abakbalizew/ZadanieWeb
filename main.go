package main

import (
	"net/http"
	"zadanieweb/handlers"

	"github.com/gorilla/mux"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api", handlers.Handle_api).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)

}
