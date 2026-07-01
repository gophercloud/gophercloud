package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestToTokenV3ScopeMap(t *testing.T) {
	projectID := "685038cd-3c25-4faf-8f9b-78c18e503190"
	projectName := "admin"
	domainID := "e4b515b8-e453-49d8-9cce-4bec244fa84e"
	domainName := "Default"

	var successCases = []struct {
		name     string
		opts     gophercloud.AuthOptions
		expected map[string]any
	}{
		{
			"System-scoped",
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					System: true,
				},
			},
			map[string]any{
				"system": map[string]any{
					"all": true,
				},
			},
		},
		{
			"Trust-scoped",
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					TrustID: "05144328-1f7d-46a9-a978-17eaad187077",
				},
			},
			map[string]any{
				"OS-TRUST:trust": map[string]string{
					"id": "05144328-1f7d-46a9-a978-17eaad187077",
				},
			},
		},
		{
			"Project-scoped (ID)",
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectID: projectID,
				},
			},
			map[string]any{
				"project": map[string]any{
					"id": &projectID,
				},
			},
		},
		{
			"Project-scoped (name)",
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectName: projectName,
					DomainName:  domainName,
				},
			},
			map[string]any{
				"project": map[string]any{
					"name": &projectName,
					"domain": map[string]any{
						"name": &domainName,
					},
				},
			},
		},
		{
			"Domain-scoped (ID)",
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					DomainID: domainID,
				},
			},
			map[string]any{
				"domain": map[string]any{
					"id": &domainID,
				},
			},
		},
		{
			"Domain-scoped (name)",
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					DomainName: domainName,
				},
			},
			map[string]any{
				"domain": map[string]any{
					"name": &domainName,
				},
			},
		},
		{
			"Empty with project fallback (ID)",
			gophercloud.AuthOptions{
				TenantID: projectID,
				Scope:    nil,
			},
			map[string]any{
				"project": map[string]any{
					"id": &projectID,
				},
			},
		},
		{
			"Empty with project fallback (name)",
			gophercloud.AuthOptions{
				TenantName: projectName,
				DomainName: domainName,
				Scope:      nil,
			},
			map[string]any{
				"project": map[string]any{
					"name": &projectName,
					"domain": map[string]any{
						"name": &domainName,
					},
				},
			},
		},
		{
			"Empty without fallback",
			gophercloud.AuthOptions{
				Scope: nil,
			},
			nil,
		},
	}
	for _, successCase := range successCases {
		t.Run(successCase.name, func(t *testing.T) {
			actual, err := successCase.opts.ToTokenV3ScopeMap()
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, successCase.expected, actual)
		})
	}
}
