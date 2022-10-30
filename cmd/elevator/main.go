package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/idtoken"
	"gopkg.in/yaml.v3"
)

type UpdateIamBindingRequest struct {
	Project string `json:"project"`
	Role    string `json:"role"`
	Minutes int    `json:"minutes"`
	Reason  string `json:"reason"`
}

type EntityPolicies []EntityPolicy

type EntityPolicy struct {
	Type     string   `yaml:"type"`
	Name     string   `yaml:"name"`
	Policies []Policy `yaml:"policies"`
}

type Policy struct {
	Project     string `yaml:"project"`
	Role        string `yaml:"role"`
	MaxMinutes	int    `yaml:"max_minutes"`
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!\n"))
}

func UpdateIamBinding(w http.ResponseWriter, r *http.Request) {
	var updateIamBindingRequest UpdateIamBindingRequest
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&updateIamBindingRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := r.Header.Get("user-email")

	if !ValidateRequestAgainstPolicy(user, updateIamBindingRequest) {
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

func ValidateRequestAgainstPolicy(user string, updateIamBindingRequest UpdateIamBindingRequest) bool {
	var entityPolicies EntityPolicies
	file, err := os.ReadFile("configs/policies.yaml")

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(file, &entityPolicies)

	for _, ep := range entityPolicies {
		if ep.Name == user {
			for _, p := range ep.Policies {
				if
					p.Project == updateIamBindingRequest.Project &&
					p.Role == updateIamBindingRequest.Role &&
					(p.MaxMinutes == 0 || p.MaxMinutes >= updateIamBindingRequest.Minutes) {
					LogRequestResult(user, updateIamBindingRequest, true)
					return true
				}
			}
		}
	}

	LogRequestResult(user, updateIamBindingRequest, false)
	return false
}

func RequireToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			ReturnUnauthorized(w)
		} else {
			idToken := authHeader[1]
			parsedToken, err := idtoken.Validate(ctx, idToken, "")

			if err != nil {
				fmt.Println("Error: ", err)
				ReturnUnauthorized(w)
			} else {
				r.Header.Add("user-email", parsedToken.Claims["email"].(string))
				next(w, r)
			}
		}
	}
}

func ReturnUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Unauthorized\n"))
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", IndexPage).Methods("GET")
	r.HandleFunc("/updateiam", RequireToken(UpdateIamBinding)).Methods("POST")
	http.ListenAndServe(":8080", r)
}
