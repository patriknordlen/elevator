package types

type UpdateIamBindingRequest struct {
	Project string `json:"project"`
	Role    string `json:"role"`
	Minutes int    `json:"minutes"`
	Reason  string `json:"reason"`
}
