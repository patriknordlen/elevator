package main

import (
	"net/http"

	"github.com/einride/elevator/internal/handlers"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RequireToken(handlers.IndexPage)).Methods("GET")
	r.HandleFunc("/updateiam", handlers.RequireToken(handlers.HandleUpdateIamBindingRequest)).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(":8080", r)
}
