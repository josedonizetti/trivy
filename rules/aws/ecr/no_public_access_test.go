package ecr

import (
	"testing"

	"github.com/aquasecurity/defsec/providers/aws/iam"
	"github.com/liamg/iamgo"

	"github.com/aquasecurity/defsec/parsers/types"
	"github.com/aquasecurity/defsec/providers/aws/ecr"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/state"
	"github.com/stretchr/testify/assert"
)

func TestCheckNoPublicAccess(t *testing.T) {
	tests := []struct {
		name     string
		input    ecr.ECR
		expected bool
	}{
		{
			name: "ECR repository policy with wildcard principal",
			input: ecr.ECR{
				Metadata: types.NewTestMetadata(),
				Repositories: []ecr.Repository{
					{
						Metadata: types.NewTestMetadata(),
						Policies: func() []iam.Policy {

							sb := iamgo.NewStatementBuilder()
							sb.WithSid("new policy")
							sb.WithEffect("Allow")
							sb.WithAllPrincipals(true)
							sb.WithActions([]string{
								"ecr:GetDownloadUrlForLayer",
								"ecr:BatchGetImage",
								"ecr:BatchCheckLayerAvailability",
								"ecr:PutImage",
								"ecr:InitiateLayerUpload",
								"ecr:UploadLayerPart",
								"ecr:CompleteLayerUpload",
								"ecr:DescribeRepositories",
								"ecr:GetRepositoryPolicy",
								"ecr:ListImages",
								"ecr:DeleteRepository",
								"ecr:BatchDeleteImage",
								"ecr:SetRepositoryPolicy",
								"ecr:DeleteRepositoryPolicy",
							})

							builder := iamgo.NewPolicyBuilder()
							builder.WithVersion("2021-10-07")
							builder.WithStatement(sb.Build())

							return []iam.Policy{
								{
									Document: iam.Document{
										Metadata: types.NewTestMetadata(),
										Parsed:   builder.Build(),
									},
								},
							}
						}(),
					},
				},
			},
			expected: true,
		},
		{
			name: "ECR repository policy with specific principal",
			input: ecr.ECR{
				Metadata: types.NewTestMetadata(),
				Repositories: []ecr.Repository{
					{
						Metadata: types.NewTestMetadata(),
						Policies: func() []iam.Policy {

							sb := iamgo.NewStatementBuilder()
							sb.WithSid("new policy")
							sb.WithEffect("Allow")
							sb.WithAWSPrincipals([]string{"arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"})
							sb.WithActions([]string{
								"ecr:GetDownloadUrlForLayer",
								"ecr:BatchGetImage",
								"ecr:BatchCheckLayerAvailability",
								"ecr:PutImage",
								"ecr:InitiateLayerUpload",
								"ecr:UploadLayerPart",
								"ecr:CompleteLayerUpload",
								"ecr:DescribeRepositories",
								"ecr:GetRepositoryPolicy",
								"ecr:ListImages",
								"ecr:DeleteRepository",
								"ecr:BatchDeleteImage",
								"ecr:SetRepositoryPolicy",
								"ecr:DeleteRepositoryPolicy",
							})

							builder := iamgo.NewPolicyBuilder()
							builder.WithVersion("2021-10-07")
							builder.WithStatement(sb.Build())

							return []iam.Policy{
								{
									Document: iam.Document{
										Metadata: types.NewTestMetadata(),
										Parsed:   builder.Build(),
									},
								},
							}
						}(),
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.ECR = test.input
			results := CheckNoPublicAccess.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == rules.StatusFailed && result.Rule().LongID() == CheckNoPublicAccess.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}