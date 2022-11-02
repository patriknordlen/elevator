package policy

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/einride/elevator/internal/types"
	"google.golang.org/api/cloudidentity/v1"
)

func ValidateRequestAgainstPolicy(user string, updateIamBindingRequest types.UpdateIamBindingRequest) bool {
	for _, p := range GetPoliciesForUser(user) {
		if
			p.Project == updateIamBindingRequest.Project &&
			p.Role == updateIamBindingRequest.Role &&
			(p.MaxMinutes == 0 || p.MaxMinutes >= updateIamBindingRequest.Minutes) {
			return true
		}
	}

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
		if (ep.Type == "user" && ep.Name == user) || (ep.Type == "group" && UserIsMemberOfGroup(user, ep.Name)) {
			userPolicies = append(userPolicies, ep.Policies...)
		}
	}

	return userPolicies
}


func UserIsMemberOfGroup(user string, group string) bool {
	ctx := context.Background()

	cisvc, _ := cloudidentity.NewService(ctx)
	group_id, err := cisvc.Groups.Lookup().GroupKeyId(group).Do()
	res, err := cisvc.Groups.Memberships.CheckTransitiveMembership(group_id.Name).Query(fmt.Sprintf("member_key_id=='%s'", user)).Do()

	if err != nil {
		log.Println(err)
	}
	
	return res.HasMembership
}

// This could be used if permissions allow. See https://pkg.go.dev/google.golang.org/api@v0.101.0/cloudidentity/v1#GroupsMembershipsService.SearchTransitiveGroups
func GetUserGroups(user string) string {
	ctx := context.Background()
	cisvc, _ := cloudidentity.NewService(ctx)

	var groups []string
	fmt.Println(user)
	ret, err := cisvc.Groups.Memberships.SearchTransitiveGroups("groups/-").Query(fmt.Sprintf("member_key_id=='%s' && 'cloudidentity.googleapis.com/groups.discussion_forum' in labels", user)).Do()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(ret)
	}

	return strings.Join(groups, "\n")
}
