package policy

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
	
	"github.com/einride/elevator/internal/types"
	"github.com/einride/elevator/internal/httputil"
)

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

func ValidateRequestAgainstPolicy(user string, updateIamBindingRequest types.UpdateIamBindingRequest) bool {
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
					httputil.LogRequestResult(user, updateIamBindingRequest, true)
					return true
				}
			}
		}
	}

	httputil.LogRequestResult(user, updateIamBindingRequest, false)
	return false
}
