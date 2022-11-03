package iam

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/cloudresourcemanager/v3"

)

func UpdateIamBinding(ctx context.Context, user string, project string, role string, minutes int, reason string) error {
	crmService, err := cloudresourcemanager.NewService(ctx)

	policy, err := crmService.Projects.GetIamPolicy(
		fmt.Sprintf("projects/%s", project),
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: 3,
			},
		},
	).Do()

	if err != nil {
		return err
	}

	newBinding := &cloudresourcemanager.Binding{
		Role:    fmt.Sprintf("roles/%s", role),
		Members: []string{fmt.Sprintf("user:%s", ctx.Value("user-email"))},
		Condition: &cloudresourcemanager.Expr{
			Title:       fmt.Sprintf("Added by elevator %s", time.Now().Format(time.RFC3339)),
			Description: fmt.Sprintf("Reason supplied by user:\n%s", reason),
			Expression:  fmt.Sprintf(`request.time < timestamp("%s")`, time.Now().Add(time.Duration(minutes)*time.Minute).Format(time.RFC3339Nano)),
		},
	}

	policy.Bindings = append(policy.Bindings, newBinding)
	policy.Version = 3
	setIamPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
	_, err = crmService.Projects.SetIamPolicy(fmt.Sprintf("projects/%s", project), setIamPolicyRequest).Do()

	return err
}

