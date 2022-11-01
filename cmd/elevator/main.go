package main

import (
	"net/http"
	"text/template"
	"github.com/einride/elevator/internal/httputil"
	"github.com/einride/elevator/internal/iam"
	"github.com/gorilla/mux"
)


func IndexPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	} {
		Name: r.Header.Get("user-email"),
	}

	t, _ := template.ParseFiles("web/template/index.html")
	t.Execute(w, data)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", httputil.RequireToken(IndexPage)).Methods("GET")
	r.HandleFunc("/updateiam", httputil.RequireToken(iam.HandleUpdateIamBindingRequest)).Methods("POST")
	http.ListenAndServe(":8080", r)
}
