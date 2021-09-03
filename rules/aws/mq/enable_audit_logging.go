package mq

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckEnableAuditLogging = rules.Register(
	rules.Rule{
		Provider:    provider.AWSProvider,
		Service:     "mq",
		ShortCode:   "enable-audit-logging",
		Summary:     "MQ Broker should have audit logging enabled",
		Impact:      "Without audit logging it is difficult to trace activity in the MQ broker",
		Resolution:  "Enable audit logging",
		Explanation: `Logging should be enabled to allow tracing of issues and activity to be investigated more fully. Logs provide additional information and context which is often invalauble during investigation`,
		Links: []string{ 
		},
		Severity: severity.Medium,
	},
	func(s *state.State) (results rules.Results) {
		for _, x := range s.AWS.S3.Buckets {
			if x.Encryption.Enabled.IsFalse() {
				results.Add(
					"",
					x.Encryption.Enabled.Metadata(),
					x.Encryption.Enabled.Value(),
				)
			}
		}
		return
	},
)