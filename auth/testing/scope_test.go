package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestScopeToScopeMapSystem(t *testing.T) {
	scope := &auth.Scope{
		System: true,
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"system": map[string]any{
			"all": true,
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapTrust(t *testing.T) {
	scope := &auth.Scope{
		TrustID: "trust-id-123",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"OS-TRUST:trust": map[string]string{
			"id": "trust-id-123",
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapProjectNameWithDomainID(t *testing.T) {
	scope := &auth.Scope{
		ProjectName:     "test-project",
		ProjectDomainID: "domain-id-123",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"name":   &scope.ProjectName,
			"domain": map[string]any{"id": &scope.ProjectDomainID},
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapProjectNameWithDomainName(t *testing.T) {
	scope := &auth.Scope{
		ProjectName:       "test-project",
		ProjectDomainName: "test-domain",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"name":   &scope.ProjectName,
			"domain": map[string]any{"name": &scope.ProjectDomainName},
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapProjectID(t *testing.T) {
	scope := &auth.Scope{
		ProjectID: "project-id-123",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"id": &scope.ProjectID,
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapDomainID(t *testing.T) {
	scope := &auth.Scope{
		DomainID: "domain-id-123",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"domain": map[string]any{
			"id": &scope.DomainID,
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapDomainName(t *testing.T) {
	scope := &auth.Scope{
		DomainName: "test-domain",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"domain": map[string]any{
			"name": &scope.DomainName,
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapEmpty(t *testing.T) {
	scope := &auth.Scope{}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)
	if scopeMap != nil {
		t.Errorf("Expected nil scope map, got %v", scopeMap)
	}
}

func TestNilScopeToScopeMap(t *testing.T) {
	var scope *auth.Scope

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, scopeMap == nil)
}

func TestScopeRejectsProjectAndDomainScope(t *testing.T) {
	tests := []struct {
		name          string
		scope         auth.Scope
		domainOption  string
		projectOption string
	}{
		{
			name:          "project ID and domain ID",
			scope:         auth.Scope{ProjectID: "project-id", DomainID: "domain-id"},
			domainOption:  "DomainID",
			projectOption: "ProjectID",
		},
		{
			name:          "project ID and domain name",
			scope:         auth.Scope{ProjectID: "project-id", DomainName: "domain-name"},
			domainOption:  "DomainName",
			projectOption: "ProjectID",
		},
		{
			name:          "project name and domain ID",
			scope:         auth.Scope{ProjectName: "project-name", ProjectDomainID: "project-domain-id", DomainID: "domain-id"},
			domainOption:  "DomainID",
			projectOption: "ProjectName",
		},
		{
			name:          "project name and domain name",
			scope:         auth.Scope{ProjectName: "project-name", ProjectDomainID: "project-domain-id", DomainName: "domain-name"},
			domainOption:  "DomainName",
			projectOption: "ProjectName",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.scope.ToScopeMap()
			th.AssertErr(t, err)
			scopeErr, ok := err.(gophercloud.ErrProjectScopeAndDomainScope)
			th.AssertEquals(t, true, ok)
			th.AssertEquals(t, test.domainOption, scopeErr.DomainOption)
			th.AssertEquals(t, test.projectOption, scopeErr.ProjectOption)
		})
	}
}

func TestScopeProjectNameRequiresProjectDomain(t *testing.T) {
	_, err := (&auth.Scope{ProjectName: "project-name"}).ToScopeMap()
	th.AssertErr(t, err)
	_, ok := err.(gophercloud.ErrScopeProjectDomainIDOrProjectDomainName)
	th.AssertEquals(t, true, ok)
}

func TestScopePrioritizesHigherLevelScopes(t *testing.T) {
	tests := []struct {
		name  string
		scope auth.Scope
		want  map[string]any
	}{
		{
			name:  "system over trust and project",
			scope: auth.Scope{System: true, TrustID: "trust-id", ProjectID: "project-id"},
			want:  map[string]any{"system": map[string]any{"all": true}},
		},
		{
			name:  "trust over project",
			scope: auth.Scope{TrustID: "trust-id", ProjectID: "project-id"},
			want:  map[string]any{"OS-TRUST:trust": map[string]string{"id": "trust-id"}},
		},
		{
			name:  "project ID over project name",
			scope: auth.Scope{ProjectID: "project-id", ProjectName: "project-name"},
			want:  map[string]any{"project": map[string]any{"id": stringPointer("project-id")}},
		},
		{
			name:  "project domain ID over project domain name",
			scope: auth.Scope{ProjectName: "project-name", ProjectDomainID: "domain-id", ProjectDomainName: "domain-name"},
			want: map[string]any{"project": map[string]any{
				"name": stringPointer("project-name"), "domain": map[string]any{"id": stringPointer("domain-id")},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.scope.ToScopeMap()
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, test.want, got)
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
