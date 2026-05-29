//go:build acceptance || keymanager || acls

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestACLCRUD(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	payload := tools.RandomString("SUPERSECRET-", 8)
	secret, err := CreateSecretWithPayload(t, client, payload)
	th.AssertNoErr(t, err)
	secretID, err := ParseID(secret.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	user := tools.RandomString("", 32)
	users := []string{user}
	iFalse := false
	setOpts := acls.SetOpts{
		acls.SetOpt{
			Type:          "read",
			Users:         &users,
			ProjectAccess: &iFalse,
		},
	}

	aclRef, err := acls.SetSecretACL(context.TODO(), client, secretID, setOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, aclRef)
	defer func() {
		err := acls.DeleteSecretACL(context.TODO(), client, secretID).ExtractErr()
		th.AssertNoErr(t, err)
		acl, err := acls.GetSecretACL(context.TODO(), client, secretID).Extract()
		th.AssertNoErr(t, err)
		tools.PrintResource(t, acl)
	}()

	acl, err := acls.GetSecretACL(context.TODO(), client, secretID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, acl)
	tools.PrintResource(t, (*acl)["read"].Created)
	th.AssertEquals(t, len((*acl)["read"].Users), 1)
	th.AssertEquals(t, (*acl)["read"].ProjectAccess, false)

	newUsers := []string{}
	updateOpts := acls.SetOpts{
		acls.SetOpt{
			Type:  "read",
			Users: &newUsers,
		},
	}

	aclRef, err = acls.UpdateSecretACL(context.TODO(), client, secretID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, aclRef)

	acl, err = acls.GetSecretACL(context.TODO(), client, secretID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, acl)
	tools.PrintResource(t, (*acl)["read"].Created)
	th.AssertEquals(t, len((*acl)["read"].Users), 0)
	th.AssertEquals(t, (*acl)["read"].ProjectAccess, false)

	container, err := CreateGenericContainer(t, client, secret)
	th.AssertNoErr(t, err)
	containerID, err := ParseID(container.ContainerRef)
	th.AssertNoErr(t, err)
	defer DeleteContainer(t, client, containerID)

	aclRef, err = acls.SetContainerACL(context.TODO(), client, containerID, setOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, aclRef)
	defer func() {
		err := acls.DeleteContainerACL(context.TODO(), client, containerID).ExtractErr()
		th.AssertNoErr(t, err)
		acl, err := acls.GetContainerACL(context.TODO(), client, containerID).Extract()
		th.AssertNoErr(t, err)
		tools.PrintResource(t, acl)
	}()

	acl, err = acls.GetContainerACL(context.TODO(), client, containerID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, acl)
	tools.PrintResource(t, (*acl)["read"].Created)
	th.AssertEquals(t, len((*acl)["read"].Users), 1)
	th.AssertEquals(t, (*acl)["read"].ProjectAccess, false)

	aclRef, err = acls.UpdateContainerACL(context.TODO(), client, containerID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, aclRef)

	acl, err = acls.GetContainerACL(context.TODO(), client, containerID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, acl)
	tools.PrintResource(t, (*acl)["read"].Created)
	th.AssertEquals(t, len((*acl)["read"].Users), 0)
	th.AssertEquals(t, (*acl)["read"].ProjectAccess, false)
}
