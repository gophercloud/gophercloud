//go:build acceptance || identity || applicationcredentials

package v3

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/applicationcredentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestApplicationCredentialsCRD(t *testing.T) {
	// maps are required, because Application Credential roles are returned in a random order
	rolesToMap := func(roles []applicationcredentials.Role) map[string]string {
		rolesMap := map[string]string{}
		for _, role := range roles {
			rolesMap[role.Name] = role.Name
			rolesMap[role.ID] = role.ID
		}
		return rolesMap
	}

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		UserID:     ao.UserID,
		Password:   ao.Password,
		DomainName: ao.DomainName,
		DomainID:   ao.DomainID,
		Scope: tokens.Scope{
			ProjectID:   ao.TenantID,
			ProjectName: ao.TenantName,
			DomainID:    ao.DomainID,
			DomainName:  ao.DomainName,
		},
	}

	token, err := tokens.Create(context.TODO(), client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := tokens.Get(context.TODO(), client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	roles, err := tokens.Get(context.TODO(), client, token.ID).ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, roles)

	project, err := tokens.Get(context.TODO(), client, token.ID).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)

	// prepare create parameters
	var apRoles []applicationcredentials.Role
	for i, role := range roles {
		if i%2 == 0 {
			apRoles = append(apRoles, applicationcredentials.Role{Name: role.Name})
		} else {
			apRoles = append(apRoles, applicationcredentials.Role{ID: role.ID})
		}
		if i > 4 {
			break
		}
	}
	tools.PrintResource(t, apRoles)

	// restricted, limited TTL, with limited roles, autogenerated secret
	expiresAt := time.Now().Add(time.Minute).Truncate(time.Millisecond).UTC()
	createOpts := applicationcredentials.CreateOpts{
		Name:        "test-ac",
		Description: "test application credential",
		Roles:       apRoles,
		ExpiresAt:   &expiresAt,
	}

	applicationCredential, err := applicationcredentials.Create(context.TODO(), client, user.ID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer applicationcredentials.Delete(context.TODO(), client, user.ID, applicationCredential.ID)
	tools.PrintResource(t, applicationCredential)

	if applicationCredential.Secret == "" {
		t.Fatalf("Application credential secret was not generated")
	}

	th.AssertEquals(t, applicationCredential.ExpiresAt, expiresAt)
	th.AssertEquals(t, applicationCredential.Name, createOpts.Name)
	th.AssertEquals(t, applicationCredential.Description, createOpts.Description)
	th.AssertEquals(t, applicationCredential.Unrestricted, false)
	th.AssertEquals(t, applicationCredential.ProjectID, project.ID)

	checkACroles := rolesToMap(applicationCredential.Roles)
	for i, role := range roles {
		if i%2 == 0 {
			th.AssertEquals(t, checkACroles[role.Name], role.Name)
		} else {
			th.AssertEquals(t, checkACroles[role.ID], role.ID)
		}
		if i > 4 {
			break
		}
	}

	// Get an application credential
	getApplicationCredential, err := applicationcredentials.Get(context.TODO(), client, user.ID, applicationCredential.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getApplicationCredential)

	if getApplicationCredential.Secret != "" {
		t.Fatalf("Application credential secret should not be returned by a GET request")
	}

	th.AssertEquals(t, getApplicationCredential.ExpiresAt, expiresAt)
	th.AssertEquals(t, getApplicationCredential.Name, createOpts.Name)
	th.AssertEquals(t, getApplicationCredential.Description, createOpts.Description)
	th.AssertEquals(t, getApplicationCredential.Unrestricted, false)
	th.AssertEquals(t, getApplicationCredential.ProjectID, project.ID)

	checkACroles = rolesToMap(getApplicationCredential.Roles)
	for i, role := range roles {
		if i%2 == 0 {
			th.AssertEquals(t, checkACroles[role.Name], role.Name)
		} else {
			th.AssertEquals(t, checkACroles[role.ID], role.ID)
		}
		if i > 4 {
			break
		}
	}

	// unrestricted, unlimited TTL, with all possible roles, with a custom secret
	createOpts = applicationcredentials.CreateOpts{
		Name:         "super-test-ac",
		Description:  "test unrestricted application credential",
		Unrestricted: true,
		Secret:       "myprecious",
	}

	newApplicationCredential, err := applicationcredentials.Create(context.TODO(), client, user.ID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer applicationcredentials.Delete(context.TODO(), client, user.ID, newApplicationCredential.ID)
	tools.PrintResource(t, newApplicationCredential)

	th.AssertEquals(t, newApplicationCredential.ExpiresAt, time.Time{})
	th.AssertEquals(t, newApplicationCredential.Name, createOpts.Name)
	th.AssertEquals(t, newApplicationCredential.Description, createOpts.Description)
	th.AssertEquals(t, newApplicationCredential.Secret, createOpts.Secret)
	th.AssertEquals(t, newApplicationCredential.Unrestricted, true)
	th.AssertEquals(t, newApplicationCredential.ExpiresAt, time.Time{})
	th.AssertEquals(t, newApplicationCredential.ProjectID, project.ID)

	checkACroles = rolesToMap(newApplicationCredential.Roles)
	for _, role := range roles {
		th.AssertEquals(t, checkACroles[role.Name], role.Name)
		th.AssertEquals(t, checkACroles[role.ID], role.ID)
	}
}

func TestApplicationCredentialsAccessRules(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		UserID:     ao.UserID,
		Password:   ao.Password,
		DomainName: ao.DomainName,
		DomainID:   ao.DomainID,
		Scope: tokens.Scope{
			ProjectID:   ao.TenantID,
			ProjectName: ao.TenantName,
			DomainID:    ao.DomainID,
			DomainName:  ao.DomainName,
		},
	}

	token, err := tokens.Create(context.TODO(), client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := tokens.Get(context.TODO(), client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	// prepare create parameters
	apAccessRules := []applicationcredentials.AccessRule{
		{
			Path:    "/v2.0/metrics",
			Service: "monitoring",
			Method:  "GET",
		},
		{
			Path:    "/v2.0/metrics",
			Service: "monitoring",
			Method:  "PUT",
		},
	}

	tools.PrintResource(t, apAccessRules)

	// restricted, limited TTL, with limited roles, autogenerated secret
	expiresAt := time.Now().Add(time.Minute).Truncate(time.Millisecond).UTC()
	createOpts := applicationcredentials.CreateOpts{
		Name:        "test-ac",
		Description: "test application credential",
		AccessRules: apAccessRules,
		ExpiresAt:   &expiresAt,
	}

	applicationCredential, err := applicationcredentials.Create(context.TODO(), client, user.ID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer applicationcredentials.Delete(context.TODO(), client, user.ID, applicationCredential.ID)
	tools.PrintResource(t, applicationCredential)

	if applicationCredential.Secret == "" {
		t.Fatalf("Application credential secret was not generated")
	}

	th.AssertEquals(t, applicationCredential.ExpiresAt, expiresAt)
	th.AssertEquals(t, applicationCredential.Name, createOpts.Name)
	th.AssertEquals(t, applicationCredential.Description, createOpts.Description)
	th.AssertEquals(t, applicationCredential.Unrestricted, false)

	for i, rule := range applicationCredential.AccessRules {
		th.AssertEquals(t, rule.Path, apAccessRules[i].Path)
		th.AssertEquals(t, rule.Service, apAccessRules[i].Service)
		th.AssertEquals(t, rule.Method, apAccessRules[i].Method)
	}

	// Get an application credential
	getApplicationCredential, err := applicationcredentials.Get(context.TODO(), client, user.ID, applicationCredential.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getApplicationCredential)

	if getApplicationCredential.Secret != "" {
		t.Fatalf("Application credential secret should not be returned by a GET request")
	}

	th.AssertEquals(t, getApplicationCredential.ExpiresAt, expiresAt)
	th.AssertEquals(t, getApplicationCredential.Name, createOpts.Name)
	th.AssertEquals(t, getApplicationCredential.Description, createOpts.Description)
	th.AssertEquals(t, getApplicationCredential.Unrestricted, false)

	for i, rule := range applicationCredential.AccessRules {
		th.AssertEquals(t, rule.Path, apAccessRules[i].Path)
		th.AssertEquals(t, rule.Service, apAccessRules[i].Service)
		th.AssertEquals(t, rule.Method, apAccessRules[i].Method)
	}

	// test list
	allPages, err := applicationcredentials.ListAccessRules(client, user.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := applicationcredentials.ExtractAccessRules(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, getApplicationCredential.AccessRules, actual)

	// test individual get
	for i, rule := range actual {
		getRule, err := applicationcredentials.GetAccessRule(context.TODO(), client, user.ID, rule.ID).Extract()
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, actual[i], *getRule)
	}

	res := applicationcredentials.Delete(context.TODO(), client, user.ID, applicationCredential.ID)
	th.AssertNoErr(t, res.Err)

	// test delete
	for _, rule := range actual {
		res := applicationcredentials.DeleteAccessRule(context.TODO(), client, user.ID, rule.ID)
		th.AssertNoErr(t, res.Err)
	}

	allPages, err = applicationcredentials.ListAccessRules(client, user.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err = applicationcredentials.ExtractAccessRules(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(actual), 0)
}
