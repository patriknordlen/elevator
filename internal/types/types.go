package types

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
