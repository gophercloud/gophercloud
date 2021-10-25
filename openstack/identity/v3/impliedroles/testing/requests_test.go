package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/impliedroles"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListImpliedRoles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListImpliedRolesSuccessfully(t)

	roles, err := impliedroles.GetImpliesRoles(client.ServiceClient(), "42c764f0c19146728dbfe73a49cc35c3").Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, GetImpliedRole, *roles)
}
