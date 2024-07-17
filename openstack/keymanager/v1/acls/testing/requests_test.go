package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetSecretACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSecretACLSuccessfully(t, fakeServer)

	actual, err := acls.GetSecretACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedACL, *actual)
}

func TestGetContainerACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetContainerACLSuccessfully(t, fakeServer)

	actual, err := acls.GetContainerACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedACL, *actual)
}

func TestSetSecretACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetSecretACLSuccessfully(t, fakeServer)

	users := []string{"GG27dVwR9gBMnsOaRoJ1DFJmZfdVjIdW"}
	iFalse := false
	setOpts := acls.SetOpts{
		acls.SetOpt{
			Type:          "read",
			Users:         &users,
			ProjectAccess: &iFalse,
		},
	}

	actual, err := acls.SetSecretACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", setOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedSecretACLRef, *actual)
}

func TestSetContainerACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetContainerACLSuccessfully(t, fakeServer)

	users := []string{"GG27dVwR9gBMnsOaRoJ1DFJmZfdVjIdW"}
	iFalse := false
	setOpts := acls.SetOpts{
		acls.SetOpt{
			Type:          "read",
			Users:         &users,
			ProjectAccess: &iFalse,
		},
	}

	actual, err := acls.SetContainerACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", setOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedContainerACLRef, *actual)
}

func TestDeleteSecretACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSecretACLSuccessfully(t, fakeServer)

	res := acls.DeleteSecretACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
}

func TestDeleteContainerACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteContainerACLSuccessfully(t, fakeServer)

	res := acls.DeleteContainerACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateSecretACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSecretACLSuccessfully(t, fakeServer)

	newUsers := []string{}
	updateOpts := acls.SetOpts{
		acls.SetOpt{
			Type:  "read",
			Users: &newUsers,
		},
	}

	actual, err := acls.UpdateSecretACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedSecretACLRef, *actual)
}

func TestUpdateContainerACL(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateContainerACLSuccessfully(t, fakeServer)

	newUsers := []string{}
	updateOpts := acls.SetOpts{
		acls.SetOpt{
			Type:  "read",
			Users: &newUsers,
		},
	}

	actual, err := acls.UpdateContainerACL(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedContainerACLRef, *actual)
}
