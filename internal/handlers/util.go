package handlers

import (
	"log"
	"net/http"
)

func LogRequestResult(user string, updateIamBindingRequest UpdateIamBindingRequest, allowed bool) {
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

func ReturnUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	if _, err := w.Write([]byte("Unauthorized\n")); err != nil {
		log.Println(err)
	}
}
