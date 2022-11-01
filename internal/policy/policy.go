package policy

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
	
	"github.com/einride/elevator/internal/types"
	"github.com/einride/elevator/internal/httputil"
)

func ValidateRequestAgainstPolicy(user string, updateIamBindingRequest types.UpdateIamBindingRequest) bool {
	for _, p := range GetPoliciesForUser(user) {
		if
			p.Project == updateIamBindingRequest.Project &&
			p.Role == updateIamBindingRequest.Role &&
			(p.MaxMinutes == 0 || p.MaxMinutes >= updateIamBindingRequest.Minutes) {
			httputil.LogRequestResult(user, updateIamBindingRequest, true)
			return true
		}
	}

	httputil.LogRequestResult(user, updateIamBindingRequest, false)
	return false
}

func GetPoliciesForUser(user string) []types.Policy {
	var entityPolicies types.EntityPolicies
	var userPolicies []types.Policy

	file, err := os.ReadFile("configs/policies.yaml")

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(file, &entityPolicies)

	for _, ep := range entityPolicies {
		if ep.Name == user {
			userPolicies = append(userPolicies, ep.Policies...)
		}
	}

	return userPolicies
}