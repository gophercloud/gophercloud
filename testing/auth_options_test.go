package testing

import (
	"reflect"
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
		opts     gophercloud.AuthOptions
		expected map[string]interface{}
	}{
		// System-scoped
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					System: true,
				},
			},
			map[string]interface{}{
				"system": map[string]interface{}{
					"all": true,
				},
			},
		},
		// Project-scoped (ID)
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectID: projectID,
				},
			},
			map[string]interface{}{
				"project": map[string]interface{}{
					"id": &projectID,
				},
			},
		},
		// Project-scoped (name)
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectName: projectName,
					DomainName:  domainName,
				},
			},
			map[string]interface{}{
				"project": map[string]interface{}{
					"name": &projectName,
					"domain": map[string]interface{}{
						"name": &domainName,
					},
				},
			},
		},
		// Domain-scoped (ID)
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					DomainID: domainID,
				},
			},
			map[string]interface{}{
				"domain": map[string]interface{}{
					"id": &domainID,
				},
			},
		},
		// Domain-scoped (name)
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					DomainName: domainName,
				},
			},
			map[string]interface{}{
				"domain": map[string]interface{}{
					"name": &domainName,
				},
			},
		},
		// Empty with project fallback (ID)
		{
			gophercloud.AuthOptions{
				TenantID: projectID,
				Scope:    nil,
			},
			map[string]interface{}{
				"project": map[string]interface{}{
					"id": &projectID,
				},
			},
		},
		// Empty with project fallback (name)
		{
			gophercloud.AuthOptions{
				TenantName: projectName,
				DomainName: domainName,
				Scope:      nil,
			},
			map[string]interface{}{
				"project": map[string]interface{}{
					"name": &projectName,
					"domain": map[string]interface{}{
						"name": &domainName,
					},
				},
			},
		},
		// Empty without fallback
		{
			gophercloud.AuthOptions{
				Scope: nil,
			},
			nil,
		},
	}
	for _, successCase := range successCases {
		actual, err := successCase.opts.ToTokenV3ScopeMap()
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, successCase.expected, actual)
	}

	var failCases = []struct {
		opts     gophercloud.AuthOptions
		expected error
	}{
		// Project-scoped with name but missing domain ID/name
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectName: "admin",
				},
			},
			gophercloud.ErrScopeDomainIDOrDomainName{},
		},
		// Project-scoped with both project name and project ID
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectName: "admin",
					ProjectID:   "685038cd-3c25-4faf-8f9b-78c18e503190",
					DomainName:  "Default",
				},
			},
			gophercloud.ErrScopeProjectIDOrProjectName{},
		},
		// Project-scoped with name and unnecessary domain ID
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectID: "685038cd-3c25-4faf-8f9b-78c18e503190",
					DomainID:  "e4b515b8-e453-49d8-9cce-4bec244fa84e",
				},
			},
			gophercloud.ErrScopeProjectIDAlone{},
		},
		// Project-scoped with name and unnecessary domain name
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					ProjectID:  "685038cd-3c25-4faf-8f9b-78c18e503190",
					DomainName: "Default",
				},
			},
			gophercloud.ErrScopeProjectIDAlone{},
		},
		// Domain-scoped with both domain name and domain ID
		{
			gophercloud.AuthOptions{
				Scope: &gophercloud.AuthScope{
					DomainID:   "e4b515b8-e453-49d8-9cce-4bec244fa84e",
					DomainName: "Default",
				},
			},
			gophercloud.ErrScopeDomainIDOrDomainName{},
		},
	}
	for _, failCase := range failCases {
		_, err := failCase.opts.ToTokenV3ScopeMap()
		th.AssertDeepEquals(t, reflect.TypeOf(failCase.expected), reflect.TypeOf(err))
	}
}
