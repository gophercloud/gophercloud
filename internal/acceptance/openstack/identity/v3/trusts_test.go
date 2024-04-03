//go:build acceptance || identity || trusts

package v3

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/trusts"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTrustCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	// Generate a token and obtain the Admin user's ID from it.
	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		Password:   ao.Password,
		DomainName: ao.DomainName,
		DomainID:   ao.DomainID,
	}

	token, err := tokens.Create(context.TODO(), client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	adminUser, err := tokens.Get(context.TODO(), client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)

	// Get the admin and member role IDs.
	adminRoleID := ""
	memberRoleID := ""
	allPages, err := roles.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allRoles {
		if v.Name == "admin" {
			adminRoleID = v.ID
		}

		if v.Name == "member" {
			memberRoleID = v.ID
		}
	}

	// Create a project to apply the trust.
	trusteeProject, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, trusteeProject.ID)

	tools.PrintResource(t, trusteeProject)

	// Add the admin user to the trustee project.
	assignOpts := roles.AssignOpts{
		UserID:    adminUser.ID,
		ProjectID: trusteeProject.ID,
	}

	err = roles.Assign(context.TODO(), client, adminRoleID, assignOpts).ExtractErr()
	th.AssertNoErr(t, err)

	// Create a user as the trustee.
	trusteeUserCreateOpts := users.CreateOpts{
		Password: "secret",
		DomainID: "default",
	}
	trusteeUser, err := CreateUser(t, client, &trusteeUserCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, trusteeUser.ID)

	expiresAt := time.Now().Add(time.Minute).Truncate(time.Second).UTC()
	// Create a trust.
	trust, err := CreateTrust(t, client, trusts.CreateOpts{
		TrusteeUserID: trusteeUser.ID,
		TrustorUserID: adminUser.ID,
		ProjectID:     trusteeProject.ID,
		ExpiresAt:     &expiresAt,
		Roles: []trusts.Role{
			{
				ID: memberRoleID,
			},
		},
	})
	th.AssertNoErr(t, err)
	defer DeleteTrust(t, client, trust.ID)

	trust, err = FindTrust(t, client)
	th.AssertNoErr(t, err)

	// Get trust
	p, err := trusts.Get(context.TODO(), client, trust.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, p.ExpiresAt, expiresAt)
	th.AssertEquals(t, p.DeletedAt.IsZero(), true)

	tools.PrintResource(t, p)

	// List trust roles
	rolesPages, err := trusts.ListRoles(client, p.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allTrustRoles, err := trusts.ExtractRoles(rolesPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(allTrustRoles), 1)
	th.AssertEquals(t, allTrustRoles[0].ID, memberRoleID)

	// Get trust role
	role, err := trusts.GetRole(context.TODO(), client, p.ID, memberRoleID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, role.ID, memberRoleID)

	// Check trust role
	err = trusts.CheckRole(context.TODO(), client, p.ID, memberRoleID).ExtractErr()
	th.AssertNoErr(t, err)
}
