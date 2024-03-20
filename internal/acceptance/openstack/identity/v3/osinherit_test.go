//go:build acceptance || identity || osinherit

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/osinherit"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestInheritRolesAssignToUserOnProject(t *testing.T) {
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

	t.Logf("Attempting to assign an inherited role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	assignOpts := osinherit.AssignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	err = osinherit.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	validateOpts := osinherit.ValidateOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	err = osinherit.Validate(context.TODO(), client, role.ID, validateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully validated inherited role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

	unassignOpts := osinherit.UnassignOpts{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	err = osinherit.Unassign(context.TODO(), client, role.ID, unassignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully unassigned inherited role %s to a user %s on a project %s",
		role.Name, user.Name, project.Name)

}

func TestInheritRolesAssignToUserOnDomain(t *testing.T) {
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

	assignOpts := osinherit.AssignOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	}

	err = osinherit.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

	validateOpts := osinherit.ValidateOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	}

	err = osinherit.Validate(context.TODO(), client, role.ID, validateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully validated inherited role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

	unassignOpts := osinherit.UnassignOpts{
		UserID:   user.ID,
		DomainID: domain.ID,
	}

	err = osinherit.Unassign(context.TODO(), client, role.ID, unassignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully unassigned inherited role %s to a user %s on a domain %s",
		role.Name, user.Name, domain.Name)

}

func TestInheritRolesAssignToGroupOnDomain(t *testing.T) {
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

	assignOpts := osinherit.AssignOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	}

	err = osinherit.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

	validateOpts := osinherit.ValidateOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	}

	err = osinherit.Validate(context.TODO(), client, role.ID, validateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully validated inherited role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

	unassignOpts := osinherit.UnassignOpts{
		GroupID:  group.ID,
		DomainID: domain.ID,
	}

	err = osinherit.Unassign(context.TODO(), client, role.ID, unassignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully unassigned inherited role %s to a group %s on a domain %s",
		role.Name, group.Name, domain.Name)

}

func TestInheritRolesAssignToGroupOnProject(t *testing.T) {
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

	assignOpts := osinherit.AssignOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	}
	err = osinherit.Assign(context.TODO(), client, role.ID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully assigned a role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

	validateOpts := osinherit.ValidateOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	}
	err = osinherit.Validate(context.TODO(), client, role.ID, validateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully validated inherited role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

	unassignOpts := osinherit.UnassignOpts{
		GroupID:   group.ID,
		ProjectID: project.ID,
	}
	err = osinherit.Unassign(context.TODO(), client, role.ID, unassignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Successfully unassigned inherited role %s to a group %s on a project %s",
		role.Name, group.Name, project.Name)

}
