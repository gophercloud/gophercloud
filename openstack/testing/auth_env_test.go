package testing

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// clearAuthEnv unsets all OS_* env vars that AuthOptionsFromEnv reads.
func clearAuthEnv(t *testing.T) {
	t.Helper()
	vars := []string{
		"OS_AUTH_URL", "OS_USERNAME", "OS_USERID", "OS_PASSWORD", "OS_PASSCODE",
		"OS_TENANT_ID", "OS_TENANT_NAME", "OS_PROJECT_ID", "OS_PROJECT_NAME",
		"OS_DOMAIN_ID", "OS_DOMAIN_NAME",
		"OS_USER_DOMAIN_ID", "OS_USER_DOMAIN_NAME",
		"OS_PROJECT_DOMAIN_ID", "OS_PROJECT_DOMAIN_NAME",
		"OS_APPLICATION_CREDENTIAL_ID", "OS_APPLICATION_CREDENTIAL_NAME",
		"OS_APPLICATION_CREDENTIAL_SECRET", "OS_SYSTEM_SCOPE",
	}
	for _, v := range vars {
		t.Setenv(v, "")
		os.Unsetenv(v)
	}
}

// setAdminOpenRC sets env vars matching a default OpenStack admin-openrc.sh.
func setAdminOpenRC(t *testing.T) {
	t.Helper()
	t.Setenv("OS_AUTH_URL", "http://127.0.0.0/identity")
	t.Setenv("OS_PROJECT_ID", "b4b0fe18900348e488fbb0a78d3cd5f4")
	t.Setenv("OS_PROJECT_NAME", "admin")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")
	t.Setenv("OS_PROJECT_DOMAIN_ID", "default")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
}

// TestAuthOptionsFromEnv_AdminOpenRC mirrors a real admin-openrc.sh
// downloaded from Horizon. It has both OS_PROJECT_ID and OS_PROJECT_NAME,
// uses OS_USER_DOMAIN_NAME (not OS_DOMAIN_NAME), and OS_PROJECT_DOMAIN_ID.
func TestAuthOptionsFromEnv_AdminOpenRC(t *testing.T) {
	clearAuthEnv(t)
	setAdminOpenRC(t)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "http://127.0.0.0/identity", ao.IdentityEndpoint)
	th.AssertEquals(t, "admin", ao.Username)
	th.AssertEquals(t, "secret", ao.Password)
	th.AssertEquals(t, "b4b0fe18900348e488fbb0a78d3cd5f4", ao.TenantID)
	th.AssertEquals(t, "admin", ao.TenantName)
	th.AssertEquals(t, "Default", ao.DomainName)
	th.AssertEquals(t, "", ao.DomainID)
}

// TestAuthOptionsFromEnv_DemoOpenRC mirrors a demo-openrc.sh with a
// different project but same structure as admin.
func TestAuthOptionsFromEnv_DemoOpenRC(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "http://127.0.0.0/identity")
	t.Setenv("OS_PROJECT_ID", "885bacd9c88d41eda6626e6a462cba49")
	t.Setenv("OS_PROJECT_NAME", "demo")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")
	t.Setenv("OS_PROJECT_DOMAIN_ID", "default")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "885bacd9c88d41eda6626e6a462cba49", ao.TenantID)
	th.AssertEquals(t, "demo", ao.TenantName)
	th.AssertEquals(t, "Default", ao.DomainName)
}

// TestAuthOptionsFromEnv_UserDomainNamePrecedence verifies that
// OS_USER_DOMAIN_NAME takes precedence over OS_DOMAIN_NAME.
func TestAuthOptionsFromEnv_UserDomainNamePrecedence(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_ID", "project-id-123")
	t.Setenv("OS_DOMAIN_NAME", "legacy-domain")
	t.Setenv("OS_USER_DOMAIN_NAME", "user-domain")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "user-domain", ao.DomainName)
}

// TestAuthOptionsFromEnv_UserDomainIDPrecedence verifies that
// OS_USER_DOMAIN_ID takes precedence over OS_DOMAIN_ID.
func TestAuthOptionsFromEnv_UserDomainIDPrecedence(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_ID", "project-id-123")
	t.Setenv("OS_DOMAIN_ID", "legacy-domain-id")
	t.Setenv("OS_USER_DOMAIN_ID", "user-domain-id")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "user-domain-id", ao.DomainID)
	th.AssertEquals(t, "", ao.DomainName)
}

// TestAuthOptionsFromEnv_LegacyDomainNameFallback verifies that
// OS_DOMAIN_NAME is still used when OS_USER_DOMAIN_NAME is not set.
func TestAuthOptionsFromEnv_LegacyDomainNameFallback(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")
	t.Setenv("OS_DOMAIN_NAME", "Default")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "Default", ao.DomainName)
	th.AssertEquals(t, "myproject", ao.TenantName)
}

// TestAuthOptionsFromEnv_LegacyDomainIDFallback verifies that OS_DOMAIN_ID
// is still used when OS_USER_DOMAIN_ID is not set.
func TestAuthOptionsFromEnv_LegacyDomainIDFallback(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_ID", "project-123")
	t.Setenv("OS_DOMAIN_ID", "domain-id-123")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "domain-id-123", ao.DomainID)
}

// TestAuthOptionsFromEnv_ProjectNameWithProjectDomainID verifies that
// OS_PROJECT_DOMAIN_ID satisfies the domain requirement for project-name
// scoping, even without OS_DOMAIN_*.
func TestAuthOptionsFromEnv_ProjectNameWithProjectDomainID(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")
	t.Setenv("OS_PROJECT_DOMAIN_ID", "project-domain-id")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "Default", ao.DomainName)
	th.AssertEquals(t, "myproject", ao.TenantName)
}

// TestAuthOptionsFromEnv_ProjectNameWithProjectDomainName verifies that
// OS_PROJECT_DOMAIN_NAME satisfies the domain requirement.
func TestAuthOptionsFromEnv_ProjectNameWithProjectDomainName(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")
	t.Setenv("OS_PROJECT_DOMAIN_NAME", "ProjectDomain")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "Default", ao.DomainName)
	th.AssertEquals(t, "myproject", ao.TenantName)
}

// TestAuthOptionsFromEnv_ProjectNameWithUserDomainOnly verifies that when
// only OS_USER_DOMAIN_NAME is set (no OS_PROJECT_DOMAIN_*), the resolved
// user domain is used as scope domain fallback via ToTokenV3ScopeMap.
func TestAuthOptionsFromEnv_ProjectNameWithUserDomainOnly(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "Default", ao.DomainName)
	th.AssertEquals(t, "myproject", ao.TenantName)

	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	projectName := "myproject"
	domainName := "Default"
	th.AssertDeepEquals(t, map[string]any{
		"project": map[string]any{
			"name":   &projectName,
			"domain": map[string]any{"name": &domainName},
		},
	}, scopeMap)
}

// TestAuthOptionsFromEnv_MissingDomainWithProjectName verifies that when
// OS_PROJECT_NAME is set but no domain vars or project ID are provided,
// an error is returned.
func TestAuthOptionsFromEnv_MissingDomainWithProjectName(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")

	_, err := openstack.AuthOptionsFromEnv()
	if err == nil {
		t.Fatal("expected error when project name is set without any domain vars or project ID")
	}
}

// TestAuthOptionsFromEnv_ProjectIDBypassesDomainCheck verifies that
// OS_PROJECT_ID satisfies the validation even without any domain vars.
func TestAuthOptionsFromEnv_ProjectIDBypassesDomainCheck(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERID", "user-id-123")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")
	t.Setenv("OS_PROJECT_ID", "project-id-123")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "project-id-123", ao.TenantID)
}

// TestAuthOptionsFromEnv_SystemScope verifies OS_SYSTEM_SCOPE=all creates
// a system-scoped AuthScope.
func TestAuthOptionsFromEnv_SystemScope(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_SYSTEM_SCOPE", "all")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	if ao.Scope == nil {
		t.Fatal("expected Scope to be set for system scope")
	}
	th.AssertEquals(t, true, ao.Scope.System)
}

// TestAuthOptionsFromEnv_NoScope verifies that when no project or system
// scope vars are set, Scope remains nil.
func TestAuthOptionsFromEnv_NoScope(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "Default", ao.DomainName)
	if ao.Scope != nil {
		t.Fatalf("expected Scope to be nil, got %+v", ao.Scope)
	}
}

// TestAuthOptionsFromEnv_AppCredentialWithUserDomain verifies that
// application credential auth with username accepts OS_USER_DOMAIN_*.
func TestAuthOptionsFromEnv_AppCredentialWithUserDomain(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_APPLICATION_CREDENTIAL_NAME", "mycred")
	t.Setenv("OS_APPLICATION_CREDENTIAL_SECRET", "supersecret")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "Default", ao.DomainName)
	th.AssertEquals(t, "mycred", ao.ApplicationCredentialName)
	th.AssertEquals(t, "supersecret", ao.ApplicationCredentialSecret)
}

// TestAuthOptionsFromEnv_AppCredentialMissingDomain verifies that
// app credential auth with username but no domain vars errors with
// a message listing both OS_USER_DOMAIN_* and OS_DOMAIN_*.
func TestAuthOptionsFromEnv_AppCredentialMissingDomain(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_APPLICATION_CREDENTIAL_NAME", "mycred")
	t.Setenv("OS_APPLICATION_CREDENTIAL_SECRET", "supersecret")

	_, err := openstack.AuthOptionsFromEnv()
	if err == nil {
		t.Fatal("expected error for app credential with username but no domain")
	}
	envErr, ok := err.(gophercloud.ErrMissingAnyoneOfEnvironmentVariables)
	if !ok {
		t.Fatalf("expected ErrMissingAnyoneOfEnvironmentVariables, got %T: %v", err, err)
	}
	th.AssertEquals(t, 4, len(envErr.EnvironmentVariables))
}

// TestAuthOptionsFromEnv_BothUserDomainIDAndName verifies that when both
// OS_USER_DOMAIN_ID and OS_USER_DOMAIN_NAME are set, both are reflected
// in AuthOptions (downstream auth_options.go handles any conflict).
func TestAuthOptionsFromEnv_BothUserDomainIDAndName(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_ID", "project-123")
	t.Setenv("OS_USER_DOMAIN_ID", "domain-id")
	t.Setenv("OS_USER_DOMAIN_NAME", "DomainName")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "domain-id", ao.DomainID)
	th.AssertEquals(t, "DomainName", ao.DomainName)
}

// TestAuthOptionsFromEnv_ProjectOverridesTenant verifies that OS_PROJECT_ID
// and OS_PROJECT_NAME override OS_TENANT_ID and OS_TENANT_NAME.
func TestAuthOptionsFromEnv_ProjectOverridesTenant(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_TENANT_ID", "old-tenant-id")
	t.Setenv("OS_TENANT_NAME", "old-tenant")
	t.Setenv("OS_PROJECT_ID", "new-project-id")
	t.Setenv("OS_PROJECT_NAME", "new-project")
	t.Setenv("OS_USER_DOMAIN_NAME", "Default")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "new-project-id", ao.TenantID)
	th.AssertEquals(t, "new-project", ao.TenantName)
}

// TestAuthOptionsFromEnv_ScopeDomainFallbackToUserDomain verifies that
// when OS_PROJECT_DOMAIN_* is not set, the scope domain validation falls
// back to the resolved user domain (domainID/domainName).
func TestAuthOptionsFromEnv_ScopeDomainFallbackToUserDomain(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	t.Setenv("OS_USERNAME", "admin")
	t.Setenv("OS_PASSWORD", "secret")
	t.Setenv("OS_PROJECT_NAME", "myproject")
	t.Setenv("OS_USER_DOMAIN_ID", "user-domain-id")

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "user-domain-id", ao.DomainID)
	th.AssertEquals(t, "myproject", ao.TenantName)
}

// --- Tests for relaxed ToTokenV3ScopeMap ---

// TestToTokenV3ScopeMap_ProjectNameWithoutDomain verifies that after
// relaxing ToTokenV3ScopeMap, a ProjectName without domain info no
// longer errors.
func TestToTokenV3ScopeMap_ProjectNameWithoutDomain(t *testing.T) {
	ao := gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{
			ProjectName: "admin",
		},
	}
	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	if scopeMap != nil {
		t.Fatalf("expected nil scope map for ProjectName without domain, got %v", scopeMap)
	}
}

// TestToTokenV3ScopeMap_ProjectIDWithDomainID verifies that ProjectID
// with DomainID no longer errors.
func TestToTokenV3ScopeMap_ProjectIDWithDomainID(t *testing.T) {
	projectID := "685038cd-3c25-4faf-8f9b-78c18e503190"
	ao := gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{
			ProjectID: projectID,
			DomainID:  "e4b515b8-e453-49d8-9cce-4bec244fa84e",
		},
	}
	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"project": map[string]any{
			"id": &projectID,
		},
	}, scopeMap)
}

// TestToTokenV3ScopeMap_ProjectIDWithDomainName verifies that ProjectID
// with DomainName no longer errors.
func TestToTokenV3ScopeMap_ProjectIDWithDomainName(t *testing.T) {
	projectID := "685038cd-3c25-4faf-8f9b-78c18e503190"
	ao := gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{
			ProjectID:  projectID,
			DomainName: "Default",
		},
	}
	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"project": map[string]any{
			"id": &projectID,
		},
	}, scopeMap)
}

// TestToTokenV3ScopeMap_DomainIDAndDomainName verifies that having both
// DomainID and DomainName no longer errors — DomainID takes precedence.
func TestToTokenV3ScopeMap_DomainIDAndDomainName(t *testing.T) {
	domainID := "e4b515b8-e453-49d8-9cce-4bec244fa84e"
	ao := gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{
			DomainID:   domainID,
			DomainName: "Default",
		},
	}
	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"domain": map[string]any{
			"id": &domainID,
		},
	}, scopeMap)
}

// TestToTokenV3ScopeMap_ProjectNameAndProjectID verifies that having both
// ProjectName and ProjectID no longer errors — ProjectName takes
// precedence (checked first in the if-else chain).
func TestToTokenV3ScopeMap_ProjectNameAndProjectID(t *testing.T) {
	projectName := "admin"
	domainID := "e4b515b8-e453-49d8-9cce-4bec244fa84e"
	ao := gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{
			ProjectName: projectName,
			ProjectID:   "685038cd-3c25-4faf-8f9b-78c18e503190",
			DomainID:    domainID,
		},
	}
	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"project": map[string]any{
			"name":   &projectName,
			"domain": map[string]any{"id": &domainID},
		},
	}, scopeMap)
}

// --- End-to-end tests ---

// TestEndToEnd_AdminOpenRC verifies the full auth flow from admin-openrc
// env vars through to the V3 token create map.
func TestEndToEnd_AdminOpenRC(t *testing.T) {
	clearAuthEnv(t)
	setAdminOpenRC(t)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	scopeMap, err := ao.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	authMap, err := ao.ToTokenV3CreateMap(scopeMap)
	th.AssertNoErr(t, err)

	auth := authMap["auth"].(map[string]any)

	// Scope uses project ID (TenantID takes precedence over TenantName)
	scope := auth["scope"].(map[string]any)
	project := scope["project"].(map[string]any)
	// Auto-populated scope stores ID behind a *string pointer
	switch v := project["id"].(type) {
	case *string:
		th.AssertEquals(t, ao.TenantID, *v)
	case string:
		th.AssertEquals(t, ao.TenantID, v)
	default:
		t.Fatalf("unexpected type for project id: %T", v)
	}

	// Identity uses user domain name from OS_USER_DOMAIN_NAME
	identity := auth["identity"].(map[string]any)
	password := identity["password"].(map[string]any)
	user := password["user"].(map[string]any)
	domain := user["domain"].(map[string]any)
	switch v := domain["name"].(type) {
	case *string:
		th.AssertEquals(t, "Default", *v)
	case string:
		th.AssertEquals(t, "Default", v)
	default:
		t.Fatalf("unexpected type for domain name: %T", v)
	}
}
