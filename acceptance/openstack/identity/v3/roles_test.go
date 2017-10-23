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

	role, err := FindRole(t, client)
	if err != nil {
		t.Fatalf("Unable to find a role: %v", err)
	}

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

func TestAssignToUserOnProject(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an indentity client: %v", err)
	}

	project, err := FindProject(t, client)
	if err != nil {
		t.Fatalf("Unable to get a project: %v", err)
	}

	role, err := FindRole(t, client)
	if err != nil {
		t.Fatalf("Unable to get a role: %v", err)
	}

	user, err := CreateUser(t, client, nil)
	if err != nil {
		t.Fatalf("Unable to create user: %v", err)
	}
	defer DeleteUser(t, client, user.ID)

	err = AssignRoleToUserOnProject(t, client, role, user, project)
	if err != nil {
		t.Fatalf("Unable to assign a role to a user on a project: %v", err)
	}
	defer UnassignRoleFromUserOnProject(t, client, role, user, project)

	allPages, err := roles.ListAssignments(client, roles.ListAssignmentsOpts{
		RoleID:         role.ID,
		ScopeProjectID: project.ID,
		UserID:         user.ID,
	}).AllPages()
	if err != nil {
		t.Fatalf("Unable to list role assignments: %v", err)
	}

	allRoleAssignments, err := roles.ExtractRoleAssignments(allPages)
	if err != nil {
		t.Fatalf("Unable to extract role assignments: %v", err)
	}

	t.Logf("Role assignments of user %s on project %s:", user.Name, project.Name)
	for _, roleAssignment := range allRoleAssignments {
		tools.PrintResource(t, roleAssignment)
	}
}
