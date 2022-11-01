package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/einride/elevator/internal/types"
	"github.com/einride/elevator/internal/policy"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func HandleUpdateIamBindingRequest(w http.ResponseWriter, r *http.Request) {
	var updateIamBindingRequest types.UpdateIamBindingRequest
	ctx := context.Background()
	user := r.Header.Get("user-email")
	err := json.NewDecoder(r.Body).Decode(&updateIamBindingRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !policy.ValidateRequestAgainstPolicy(user, updateIamBindingRequest) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Forbidden by policy\n"))
		return
	}

	crmService, err := cloudresourcemanager.NewService(ctx)

	policy, err := crmService.Projects.GetIamPolicy(
		fmt.Sprintf("projects/%s", updateIamBindingRequest.Project),
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: 3,
			},
		},
	).Do()

	if err != nil {
		log.Println(err)
	}

	newBinding := &cloudresourcemanager.Binding{
		Role:    fmt.Sprintf("roles/%s", updateIamBindingRequest.Role),
		Members: []string{fmt.Sprintf("user:%s", user)},
		Condition: &cloudresourcemanager.Expr{
			Title:       fmt.Sprintf("Added by elevate %s", time.Now().Format(time.RFC3339)),
			Description: fmt.Sprintf("Reason supplied by user:\n%s", updateIamBindingRequest.Reason),
			Expression:  fmt.Sprintf(`request.time < timestamp("%s")`, time.Now().Add(time.Duration(updateIamBindingRequest.Minutes)*time.Minute).Format(time.RFC3339Nano)),
		},
	}

	policy.Bindings = append(policy.Bindings, newBinding)
	policy.Version = 3
	setIamPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
	_, err = crmService.Projects.SetIamPolicy(fmt.Sprintf("projects/%s", updateIamBindingRequest.Project), setIamPolicyRequest).Do()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err)))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success!\n"))
	}
}
