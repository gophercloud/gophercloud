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
		ProjectName: "test-project",
		DomainID:    "domain-id-123",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"name":   &scope.ProjectName,
			"domain": map[string]any{"id": &scope.DomainID},
		},
	}
	th.AssertDeepEquals(t, expected, scopeMap)
}

func TestScopeToScopeMapProjectNameWithDomainName(t *testing.T) {
	scope := &auth.Scope{
		ProjectName: "test-project",
		DomainName:  "test-domain",
	}

	scopeMap, err := scope.ToScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"name":   &scope.ProjectName,
			"domain": map[string]any{"name": &scope.DomainName},
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

// Error cases for scope validation
func TestScopeToScopeMapProjectNameWithoutDomain(t *testing.T) {
	scope := &auth.Scope{
		ProjectName: "test-project",
	}

	_, err := scope.ToScopeMap()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeDomainIDOrDomainName)
	th.AssertEquals(t, true, ok)
}

func TestScopeToScopeMapProjectNameWithProjectID(t *testing.T) {
	scope := &auth.Scope{
		ProjectName: "test-project",
		ProjectID:   "project-id-123",
		DomainID:    "domain-id-123",
	}

	_, err := scope.ToScopeMap()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeProjectIDOrProjectName)
	th.AssertEquals(t, true, ok)
}

func TestScopeToScopeMapProjectIDWithDomainID(t *testing.T) {
	scope := &auth.Scope{
		ProjectID: "project-id-123",
		DomainID:  "domain-id-123",
	}

	_, err := scope.ToScopeMap()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeProjectIDAlone)
	th.AssertEquals(t, true, ok)
}

func TestScopeToScopeMapProjectIDWithDomainName(t *testing.T) {
	scope := &auth.Scope{
		ProjectID:  "project-id-123",
		DomainName: "test-domain",
	}

	_, err := scope.ToScopeMap()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeProjectIDAlone)
	th.AssertEquals(t, true, ok)
}

func TestScopeToScopeMapDomainIDWithDomainName(t *testing.T) {
	scope := &auth.Scope{
		DomainID:   "domain-id-123",
		DomainName: "test-domain",
	}

	_, err := scope.ToScopeMap()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeDomainIDOrDomainName)
	th.AssertEquals(t, true, ok)
}
