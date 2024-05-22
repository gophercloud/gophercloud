//go:build acceptance

package openstack

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/credentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthenticatedClient(t *testing.T) {
	// Obtain credentials from the environment.
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		t.Fatalf("Unable to acquire credentials: %v", err)
	}

	client, err := openstack.AuthenticatedClient(context.TODO(), ao)
	if err != nil {
		t.Fatalf("Unable to authenticate: %v", err)
	}

	if client.TokenID == "" {
		t.Errorf("No token ID assigned to the client")
	}

	t.Logf("Client successfully acquired a token: %v", client.TokenID)

	// Find the storage service in the service catalog.
	storage, err := openstack.NewObjectStorageV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		t.Errorf("Unable to locate a storage service: %v", err)
	} else {
		t.Logf("Located a storage service at endpoint: [%s]", storage.Endpoint)
	}
}

func TestEC2AuthMethod(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		Password:   ao.Password,
		DomainName: ao.DomainName,
		DomainID:   ao.DomainID,
		// We need a scope to get the token roles list
		Scope: tokens.Scope{
			ProjectID:   ao.TenantID,
			ProjectName: ao.TenantName,
			DomainID:    ao.DomainID,
			DomainName:  ao.DomainName,
		},
	}
	token, err := tokens.Create(context.TODO(), client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := tokens.Get(context.TODO(), client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	project, err := tokens.Get(context.TODO(), client, token.ID).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)

	createOpts := credentials.CreateOpts{
		ProjectID: project.ID,
		Type:      "ec2",
		UserID:    user.ID,
		Blob:      "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
	}

	// Create a credential
	credential, err := credentials.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)

	// Delete a credential
	defer credentials.Delete(context.TODO(), client, credential.ID)
	tools.PrintResource(t, credential)

	newClient, err := clients.NewIdentityV3UnauthenticatedClient()
	th.AssertNoErr(t, err)

	ec2AuthOptions := &ec2tokens.AuthOptions{
		Access: "181920",
		Secret: "secretKey",
	}

	err = openstack.AuthenticateV3(context.TODO(), newClient.ProviderClient, ec2AuthOptions, gophercloud.EndpointOpts{})
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newClient.TokenID)
}

func TestReauth(t *testing.T) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		t.Fatalf("Unable to obtain environment auth options: %v", err)
	}

	// Allow reauth
	ao.AllowReauth = true

	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		t.Fatalf("Unable to create provider: %v", err)
	}

	err = openstack.Authenticate(context.TODO(), provider, ao)
	if err != nil {
		t.Fatalf("Unable to authenticate: %v", err)
	}

	t.Logf("Creating a compute client")
	_, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		t.Fatalf("Unable to create compute client: %v", err)
	}

	t.Logf("Sleeping for 1 second")
	time.Sleep(1 * time.Second)
	t.Logf("Attempting to reauthenticate")

	err = provider.ReauthFunc(context.TODO())
	if err != nil {
		t.Fatalf("Unable to reauthenticate: %v", err)
	}

	t.Logf("Creating a compute client")
	_, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		t.Fatalf("Unable to create compute client: %v", err)
	}
}
