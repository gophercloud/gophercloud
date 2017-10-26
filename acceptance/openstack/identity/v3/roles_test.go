// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
)

func TestRolesList(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	listOpts := roles.ListOpts{
		DomainID: "default",
	}

	allPages, err := roles.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to list roles: %v", err)
	}

	allRoles, err := roles.ExtractRoles(allPages)
	if err != nil {
		t.Fatalf("Unable to extract roles: %v", err)
	}

	for _, role := range allRoles {
		tools.PrintResource(t, role)
	}
}

func TestRolesGet(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	allPages, err := roles.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list roles: %v", err)
	}

	allRoles, err := roles.ExtractRoles(allPages)
	if err != nil {
		t.Fatalf("Unable to extract roles: %v", err)
	}

	role := allRoles[0]
	p, err := roles.Get(client, role.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get role: %v", err)
	}

	tools.PrintResource(t, p)
}

func TestRoleCRUD(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	createOpts := roles.CreateOpts{
		Name:     "testrole",
		DomainID: "default",
		Extra: map[string]interface{}{
			"description": "test role description",
		},
	}

	// Create Role in the default domain
	role, err := CreateRole(t, client, &createOpts)
	if err != nil {
		t.Fatalf("Unable to create role: %v", err)
	}

	tools.PrintResource(t, role)
	tools.PrintResource(t, role.Extra)
}
