package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetSecretACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSecretACLSuccessfully(t)

	actual, err := acls.GetSecretACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedACL, *actual)
}

func TestGetContainerACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetContainerACLSuccessfully(t)

	actual, err := acls.GetContainerACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedACL, *actual)
}

func TestSetSecretACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSetSecretACLSuccessfully(t)

	users := []string{"GG27dVwR9gBMnsOaRoJ1DFJmZfdVjIdW"}
	iFalse := false
	setOpts := acls.SetOpts{
		acls.SetOpt{
			Type:          "read",
			Users:         &users,
			ProjectAccess: &iFalse,
		},
	}

	actual, err := acls.SetSecretACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", setOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedSecretACLRef, *actual)
}

func TestSetContainerACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSetContainerACLSuccessfully(t)

	users := []string{"GG27dVwR9gBMnsOaRoJ1DFJmZfdVjIdW"}
	iFalse := false
	setOpts := acls.SetOpts{
		acls.SetOpt{
			Type:          "read",
			Users:         &users,
			ProjectAccess: &iFalse,
		},
	}

	actual, err := acls.SetContainerACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", setOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedContainerACLRef, *actual)
}

func TestDeleteSecretACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSecretACLSuccessfully(t)

	res := acls.DeleteSecretACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
}

func TestDeleteContainerACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteContainerACLSuccessfully(t)

	res := acls.DeleteContainerACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateSecretACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateSecretACLSuccessfully(t)

	newUsers := []string{}
	updateOpts := acls.SetOpts{
		acls.SetOpt{
			Type:  "read",
			Users: &newUsers,
		},
	}

	actual, err := acls.UpdateSecretACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedSecretACLRef, *actual)
}

func TestUpdateContainerACL(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateContainerACLSuccessfully(t)

	newUsers := []string{}
	updateOpts := acls.SetOpts{
		acls.SetOpt{
			Type:  "read",
			Users: &newUsers,
		},
	}

	actual, err := acls.UpdateContainerACL(context.TODO(), client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedContainerACLRef, *actual)
}
