package main

import (
	"log"
	"net/http"
	"time"

	"github.com/einride/elevator/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RequireToken(handlers.IndexPage)).Methods("GET")
	r.HandleFunc("/updateiam", handlers.RequireToken(handlers.HandleUpdateIamBindingRequest)).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
