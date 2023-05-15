package handlers

import (
	"net/http"
	"text/template"

	"github.com/einride/elevator/internal/policy"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user-email").(string)
	data := struct {
		Name     string
		Policies []policy.Policy
	}{
		Name:     user,
		Policies: policy.GetUserPolicies(r.Context(), user),
	}

	t, _ := template.ParseFiles("web/template/index.html")
	t.Execute(w, data)
}