package httputil

import (
	"net/http"
	"log"
	"context"
	"strings"
	"google.golang.org/api/idtoken"

	"github.com/einride/elevator/internal/types"
)

func LogRequestResult(user string, updateIamBindingRequest types.UpdateIamBindingRequest, allowed bool) {
	var action string
	if allowed {
		action = "allow"
	} else {
		action = "reject"

	}

	log.Printf(`Elevation request: user="%s" project="%s" role="%s" minutes="%d" reason="%s" action="%s"`,
		user,
		updateIamBindingRequest.Project,
		updateIamBindingRequest.Role,
		updateIamBindingRequest.Minutes,
		updateIamBindingRequest.Reason,
		action)
}

func RequireToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) == 2 && authHeader[0] == "Bearer" {
			idToken := authHeader[1]
			parsedToken, err := idtoken.Validate(ctx, idToken, "")

			if err != nil {
				log.Println("Error: ", err)
				ReturnUnauthorized(w)
			} else {
				r.Header.Add("user-email", parsedToken.Claims["email"].(string))
				next(w, r)
			}
		} else {
			ReturnUnauthorized(w)
		}
	}
}

func ReturnUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Unauthorized\n"))
}
