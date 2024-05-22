//go:build acceptance || identity || roles

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRolesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	listOpts := roles.ListOpts{
		DomainID: "default",
	}

	allPages, err := roles.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	for _, role := range allRoles {
		tools.PrintResource(t, role)
	}
}

func TestRolesGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	role, err := FindRole(t, client)
	th.AssertNoErr(t, err)

	p, err := roles.Get(context.TODO(), client, role.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, p)
}

func TestRolesCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := roles.CreateOpts{
		Name: "testrole",
		Extra: map[string]any{
			"description": "test role description",
		},
	}

	// Create Role in the default domain
	role, err := CreateRole(t, client, &createOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	tools.PrintResource(t, role)
	tools.PrintResource(t, role.Extra)

	listOpts := roles.ListOpts{}
	allPages, err := roles.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, r := range allRoles {
		tools.PrintResource(t, r)
		tools.PrintResource(t, r.Extra)

		if r.Name == role.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	updateOpts := roles.UpdateOpts{
		Extra: map[string]any{
			"description": "updated test role description",
		},
	}

	newRole, err := roles.Update(context.TODO(), client, role.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRole)
	tools.PrintResource(t, newRole.Extra)

	th.AssertEquals(t, newRole.Extra["description"], "updated test role description")
}

func TestRolesFilterList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := roles.CreateOpts{
		Name: "testrole",
		Extra: map[string]any{
			"description": "test role description",
		},
	}

	// Create Role in the default domain
	role, err := CreateRole(t, client, &createOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	var listOpts roles.ListOpts
	listOpts.Filters = map[string]string{
		"name__contains": "test",
	}

	allPages, err := roles.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	found := false
	for _, r := range allRoles {
		tools.PrintResource(t, r)
		tools.PrintResource(t, r.Extra)

		if r.Name == role.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	listOpts.Filters = map[string]string{
		"name__contains": "reader",
	}

	allPages, err = roles.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err = roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	found = false
	for _, r := range allRoles {
		tools.PrintResource(t, r)
		tools.PrintResource(t, r.Extra)

		if r.Name == role.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, false)
}

func TestRoleListAssignmentIncludeNamesAndSubtree(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, project.ID)

	domainID := "default"
	roleCreateOpts := roles.CreateOpts{
		DomainID: domainID,
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	user, err := CreateUser(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user.ID)

	t.Logf("Attempting to assign a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	assignOpts := roles.AssignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	})

	iTrue := true
	listAssignmentsOpts := roles.ListAssignmentsOpts{
		UserID:         user.ID,
		ScopeProjectID: domainID, // set domainID in ScopeProjectID field to list assignments on all projects in domain
		IncludeSubtree: &iTrue,
		IncludeNames:   &iTrue,
	}
	allPages, err := roles.ListAssignments(client, listAssignmentsOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoleAssignments(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments(with names) of user %s on projects in domain %s:", user.Name, domainID)
	var found bool
	for _, _role := range allRoles {
		tools.PrintResource(t, _role)
		if _role.Role.ID == role.ID &&
			_role.User.Name == user.Name &&
			_role.Scope.Project.Name == project.Name &&
			_role.Scope.Project.Domain.ID == domainID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRoleListAssignmentForUserOnProject(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, project.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	user, err := CreateUser(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user.ID)

	t.Logf("Attempting to assign a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	assignOpts := roles.AssignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	})

	listAssignmentsOnResourceOpts := roles.ListAssignmentsOnResourceOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	allPages, err := roles.ListAssignmentsOnResource(client, listAssignmentsOnResourceOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of user %s on project %s:", user.Name, project.Name)
	var found bool
	for _, _role := range allRoles {
		tools.PrintResource(t, _role)

		if _role.ID == role.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRoleListAssignmentForUserOnDomain(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	domain, err := CreateDomain(t, client, &domains.CreateOpts{
		Enabled: gophercloud.Disabled,
	})
	th.AssertNoErr(t, err)
	defer DeleteDomain(t, client, domain.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	user, err := CreateUser(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user.ID)

	t.Logf("Attempting to assign a role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

	assignOpts := roles.AssignOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	}

	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	})

	listAssignmentsOnResourceOpts := roles.ListAssignmentsOnResourceOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	}
	allPages, err := roles.ListAssignmentsOnResource(client, listAssignmentsOnResourceOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of user %s on domain %s:", user.Name, domain.Name)
	var found bool
	for _, _role := range allRoles {
		tools.PrintResource(t, _role)

		if _role.ID == role.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRoleListAssignmentForGroupOnProject(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, project.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	groupCreateOpts := &groups.CreateOpts{
		DomainID: "default",
	}
	group, err := CreateGroup(t, client, groupCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, group.ID)

	t.Logf("Attempting to assign a role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

	assignOpts := roles.AssignOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	}
	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	})

	listAssignmentsOnResourceOpts := roles.ListAssignmentsOnResourceOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	}
	allPages, err := roles.ListAssignmentsOnResource(client, listAssignmentsOnResourceOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of group %s on project %s:", group.Name, project.Name)
	var found bool
	for _, _role := range allRoles {
		tools.PrintResource(t, _role)

		if _role.ID == role.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRoleListAssignmentForGroupOnDomain(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	domain, err := CreateDomain(t, client, &domains.CreateOpts{
		Enabled: gophercloud.Disabled,
	})
	th.AssertNoErr(t, err)
	defer DeleteDomain(t, client, domain.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	groupCreateOpts := &groups.CreateOpts{
		DomainID: "default",
	}
	group, err := CreateGroup(t, client, groupCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, group.ID)

	t.Logf("Attempting to assign a role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

	assignOpts := roles.AssignOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	}

	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	})

	listAssignmentsOnResourceOpts := roles.ListAssignmentsOnResourceOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	}
	allPages, err := roles.ListAssignmentsOnResource(client, listAssignmentsOnResourceOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of group %s on domain %s:", group.Name, domain.Name)
	var found bool
	for _, _role := range allRoles {
		tools.PrintResource(t, _role)

		if _role.ID == role.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRolesAssignToUserOnProject(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, project.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	user, err := CreateUser(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user.ID)

	t.Logf("Attempting to assign a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	assignOpts := roles.AssignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	})

	iTrue := true
	lao := roles.ListAssignmentsOpts{
		RoleID:         role.ID,
		ScopeProjectID: project.ID,
		UserID:         user.ID,
		IncludeNames:   &iTrue,
	}

	allPages, err := roles.ListAssignments(client, lao).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoleAssignments, err := roles.ExtractRoleAssignments(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of user %s on project %s:", user.Name, project.Name)
	var found bool
	for _, roleAssignment := range allRoleAssignments {
		tools.PrintResource(t, roleAssignment)

		if roleAssignment.Role.ID == role.ID {
			found = true
		}

		if roleAssignment.User.Domain.ID == "" || roleAssignment.Scope.Project.Domain.ID == "" {
			found = false
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRolesAssignToUserOnDomain(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	domain, err := CreateDomain(t, client, &domains.CreateOpts{
		Enabled: gophercloud.Disabled,
	})
	th.AssertNoErr(t, err)
	defer DeleteDomain(t, client, domain.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	user, err := CreateUser(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user.ID)

	t.Logf("Attempting to assign a role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

	assignOpts := roles.AssignOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	}

	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	})

	iTrue := true
	lao := roles.ListAssignmentsOpts{
		RoleID:        role.ID,
		ScopeDomainID: domain.ID,
		UserID:        user.ID,
		IncludeNames:  &iTrue,
	}

	allPages, err := roles.ListAssignments(client, lao).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoleAssignments, err := roles.ExtractRoleAssignments(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of user %s on domain %s:", user.Name, domain.Name)
	var found bool
	for _, roleAssignment := range allRoleAssignments {
		tools.PrintResource(t, roleAssignment)

		if roleAssignment.Role.ID == role.ID {
			found = true
		}

		if roleAssignment.User.Domain.ID == "" {
			found = false
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRolesAssignToGroupOnDomain(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	domain, err := CreateDomain(t, client, &domains.CreateOpts{
		Enabled: gophercloud.Disabled,
	})
	th.AssertNoErr(t, err)
	defer DeleteDomain(t, client, domain.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	groupCreateOpts := &groups.CreateOpts{
		DomainID: "default",
	}
	group, err := CreateGroup(t, client, groupCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, group.ID)

	t.Logf("Attempting to assign a role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

	assignOpts := roles.AssignOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	}

	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	})

	iTrue := true
	lao := roles.ListAssignmentsOpts{
		RoleID:        role.ID,
		ScopeDomainID: domain.ID,
		GroupID:       group.ID,
		IncludeNames:  &iTrue,
	}

	allPages, err := roles.ListAssignments(client, lao).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoleAssignments, err := roles.ExtractRoleAssignments(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of group %s on domain %s:", group.Name, domain.Name)
	var found bool
	for _, roleAssignment := range allRoleAssignments {
		tools.PrintResource(t, roleAssignment)

		if roleAssignment.Role.ID == role.ID {
			found = true
		}

		if roleAssignment.Group.Domain.ID == "" {
			found = false
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRolesAssignToGroupOnProject(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, project.ID)

	roleCreateOpts := roles.CreateOpts{
		DomainID: "default",
	}
	role, err := CreateRole(t, client, &roleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, role.ID)

	groupCreateOpts := &groups.CreateOpts{
		DomainID: "default",
	}
	group, err := CreateGroup(t, client, groupCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, group.ID)

	t.Logf("Attempting to assign a role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

	assignOpts := roles.AssignOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	}
	err = roles.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

	defer UnassignRole(t, client, role.ID, &roles.UnassignOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	})

	iTrue := true
	lao := roles.ListAssignmentsOpts{
		RoleID:         role.ID,
		ScopeProjectID: project.ID,
		GroupID:        group.ID,
		IncludeNames:   &iTrue,
	}

	allPages, err := roles.ListAssignments(client, lao).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoleAssignments, err := roles.ExtractRoleAssignments(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Role assignments of group %s on project %s:", group.Name, project.Name)
	var found bool
	for _, roleAssignment := range allRoleAssignments {
		tools.PrintResource(t, roleAssignment)

		if roleAssignment.Role.ID == role.ID {
			found = true
		}

		if roleAssignment.Scope.Project.Domain.ID == "" || roleAssignment.Group.Domain.ID == "" {
			found = false
		}
	}

	th.AssertEquals(t, found, true)
}

func TestCRUDRoleInferenceRule(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	priorRoleCreateOpts := roles.CreateOpts{
		Name: "priorRole",
		Extra: map[string]any{
			"description": "prior_role description",
		},
	}
	// Create prior_role in the default domain
	priorRole, err := CreateRole(t, client, &priorRoleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, priorRole.ID)
	tools.PrintResource(t, priorRole)
	tools.PrintResource(t, priorRole.Extra)

	impliedRoleCreateOpts := roles.CreateOpts{
		Name: "impliedRole",
		Extra: map[string]any{
			"description": "implied_role description",
		},
	}
	// Create implied_role in the default domain
	impliedRole, err := CreateRole(t, client, &impliedRoleCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteRole(t, client, impliedRole.ID)
	tools.PrintResource(t, impliedRole)
	tools.PrintResource(t, impliedRole.Extra)

	roleInferenceRule, err := roles.CreateRoleInferenceRule(context.TODO(), client, priorRole.ID, impliedRole.ID).Extract()
	defer roles.DeleteRoleInferenceRule(context.TODO(), client, priorRole.ID, impliedRole.ID)

	th.AssertNoErr(t, err)
	tools.PrintResource(t, roleInferenceRule)

	getRoleInferenceRule, err := roles.GetRoleInferenceRule(context.TODO(), client, priorRole.ID, impliedRole.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getRoleInferenceRule)

	roleInferenceRuleList, err := roles.ListRoleInferenceRules(context.TODO(), client).Extract()
	tools.PrintResource(t, roleInferenceRuleList)
	th.AssertNoErr(t, err)

}
