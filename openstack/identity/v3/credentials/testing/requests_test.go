package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/credentials"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListCredentials(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListCredentialsSuccessfully(t, fakeServer)

	count := 0
	err := credentials.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := credentials.ExtractCredentials(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedCredentialsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListCredentialsAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListCredentialsSuccessfully(t, fakeServer)

	allPages, err := credentials.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := credentials.ExtractCredentials(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedCredentialsSlice, actual)
	th.AssertDeepEquals(t, ExpectedCredentialsSlice[0].Blob, "{\"access\":\"181920\",\"secret\":\"secretKey\"}")
	th.AssertDeepEquals(t, ExpectedCredentialsSlice[1].Blob, "{\"access\":\"7da79ff0aa364e1396f067e352b9b79a\",\"secret\":\"7a18d68ba8834b799d396f3ff6f1e98c\"}")
}

func TestGetCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetCredentialSuccessfully(t, fakeServer)

	actual, err := credentials.Get(context.TODO(), client.ServiceClient(fakeServer), credentialID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Credential, *actual)
}

func TestCreateCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateCredentialSuccessfully(t, fakeServer)

	createOpts := credentials.CreateOpts{
		ProjectID: projectID,
		Type:      "ec2",
		UserID:    userID,
		Blob:      "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
	}

	CredentialResponse := Credential

	actual, err := credentials.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, CredentialResponse, *actual)
}

func TestDeleteCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteCredentialSuccessfully(t, fakeServer)

	res := credentials.Delete(context.TODO(), client.ServiceClient(fakeServer), credentialID)
	th.AssertNoErr(t, res.Err)
}

func TestUpdateCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateCredentialSuccessfully(t, fakeServer)

	updateOpts := credentials.UpdateOpts{
		ProjectID: "731fc6f265cd486d900f16e84c5cb594",
		Type:      "ec2",
		UserID:    "bb5476fd12884539b41d5a88f838d773",
		Blob:      "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
	}

	actual, err := credentials.Update(context.TODO(), client.ServiceClient(fakeServer), "2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondCredentialUpdated, *actual)
}
