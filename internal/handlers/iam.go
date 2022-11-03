package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/einride/elevator/internal/iam"
	"github.com/einride/elevator/internal/policy"
)

type UpdateIamBindingRequest struct {
	Project string `json:"project"`
	Role    string `json:"role"`
	Minutes int    `json:"minutes"`
	Reason  string `json:"reason"`
}

func HandleUpdateIamBindingRequest(w http.ResponseWriter, r *http.Request) {
	var iamReq UpdateIamBindingRequest
	user := r.Context().Value("user-email").(string)
	err := json.NewDecoder(r.Body).Decode(&iamReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !policy.ValidateRequestAgainstPolicy(r.Context(), user, iamReq.Project, iamReq.Role, iamReq.Minutes) {
		LogRequestResult(user, iamReq, false)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Forbidden by policy\n"))
		return
	}
	LogRequestResult(user, iamReq, true)

	err = iam.UpdateIamBinding(r.Context(), user, iamReq.Project, iamReq.Role, iamReq.Minutes, iamReq.Reason)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err)))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success!\n"))
	}

	return
}
