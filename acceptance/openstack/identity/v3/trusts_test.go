// +build acceptance identity trusts

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestTrustCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	// Create a trustor with project.
	trustorProject, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, trustorProject.ID)

	tools.PrintResource(t, trustorProject)

	trustorUserCreateOpts := users.CreateOpts{
		DefaultProjectID: trustorProject.ID,
		Password:         "secret",
		DomainID:         "default",
	}
	trustorUser, err := CreateUser(t, client, &trustorUserCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, trustorUser.ID)

	// Create a trustee.
	trusteeUserCreateOpts := users.CreateOpts{
		Password: "secret",
		DomainID: "default",
	}
	trusteeUser, err := CreateUser(t, client, &trusteeUserCreateOpts)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, trusteeUser.ID)

	// Create a trust.
	trust, err := CreateTrust(t, client, trusts.CreateOpts{
		TrusteeUserID: trusteeUser.ID,
		TrustorUserID: trustorUser.ID,
		ProjectID:     trustorProject.ID,
	})
	th.AssertNoErr(t, err)
	defer DeleteTrust(t, client, trust.ID)
}
