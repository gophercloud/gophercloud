package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2credentials"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListEC2Credentials(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListEC2CredentialsSuccessfully(t, fakeServer)

	count := 0
	err := ec2credentials.List(client.ServiceClient(fakeServer), userID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := ec2credentials.ExtractCredentials(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedEC2CredentialsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListEC2CredentialsAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListEC2CredentialsSuccessfully(t, fakeServer)

	allPages, err := ec2credentials.List(client.ServiceClient(fakeServer), userID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := ec2credentials.ExtractCredentials(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedEC2CredentialsSlice, actual)
}

func TestGetEC2Credential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetEC2CredentialSuccessfully(t, fakeServer)

	actual, err := ec2credentials.Get(context.TODO(), client.ServiceClient(fakeServer), userID, credentialID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, EC2Credential, *actual)
}

func TestCreateEC2Credential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateEC2CredentialSuccessfully(t, fakeServer)

	createOpts := ec2credentials.CreateOpts{
		TenantID: "6238dee2fec940a6bf31e49e9faf995a",
	}

	actual, err := ec2credentials.Create(context.TODO(), client.ServiceClient(fakeServer), userID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, EC2Credential, *actual)
}

func TestDeleteEC2Credential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteEC2CredentialSuccessfully(t, fakeServer)

	res := ec2credentials.Delete(context.TODO(), client.ServiceClient(fakeServer), userID, credentialID)
	th.AssertNoErr(t, res.Err)
}
