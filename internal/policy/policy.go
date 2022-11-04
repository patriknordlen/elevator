package policy

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"google.golang.org/api/cloudidentity/v1"
)

type EntityPolicies []EntityPolicy

type EntityPolicy struct {
	Name     string   `yaml:"name"`
	Policies []Policy `yaml:"policies"`
}

type Policy struct {
	Type		string `yaml:"type"`
	Name        string `yaml:"name"`
	Role        string `yaml:"role"`
	MaxMinutes	int    `yaml:"max_minutes"`
}

func ValidateRequestAgainstPolicy(ctx context.Context, user string, project string, role string, minutes int) bool {
	for _, p := range GetUserPolicies(ctx, user) {
		if p.Type == "project" &&
			p.Name == project &&
			p.Role == role &&
			(p.MaxMinutes == 0 || p.MaxMinutes >= minutes) {
			return true
		}
	}

	return false
}

func GetUserPolicies(ctx context.Context, user string) []Policy {
	var entityPolicies EntityPolicies
	var userPolicies []Policy

	file, err := os.ReadFile("configs/policies.yaml")

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(file, &entityPolicies)

	for _, ep := range entityPolicies {
		// Each entity is in the format "type:name"
		entity := strings.Split(ep.Name, ":")
		if (entity[0] == "user" && entity[1] == user) || (entity[0] == "group" && UserIsMemberOfGroup(ctx, user, entity[1])) {
			userPolicies = append(userPolicies, ep.Policies...)
		}
	}

	return userPolicies
}

func UserIsMemberOfGroup(ctx context.Context, user string, group string) bool {
	cisvc, _ := cloudidentity.NewService(ctx)
	group_id, err := cisvc.Groups.Lookup().GroupKeyId(group).Do()
	res, err := cisvc.Groups.Memberships.CheckTransitiveMembership(group_id.Name).Query(fmt.Sprintf("member_key_id=='%s'", user)).Do()

	if err != nil {
		log.Println(err)
	}

	return res.HasMembership
}

// This could be used if permissions allow. See https://pkg.go.dev/google.golang.org/api@v0.101.0/cloudidentity/v1#GroupsMembershipsService.SearchTransitiveGroups
func GetUserGroups(ctx context.Context, user string) string {
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
