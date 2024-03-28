//go:build acceptance || compute || keypairs

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	identity "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/identity/v3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"golang.org/x/crypto/ssh"
)

func TestKeyPairsParse(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	keyPair, err := CreateKeyPair(t, client)
	th.AssertNoErr(t, err)
	defer DeleteKeyPair(t, client, keyPair)

	// There was a series of OpenStack releases, between Liberty and Ocata,
	// where the returned SSH key was not parsable by Go.
	// This checks if the issue is happening again.
	_, err = ssh.ParsePrivateKey([]byte(keyPair.PrivateKey))
	th.AssertNoErr(t, err)

	tools.PrintResource(t, keyPair)
}

func TestKeyPairsCreateDelete(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	keyPair, err := CreateKeyPair(t, client)
	th.AssertNoErr(t, err)
	defer DeleteKeyPair(t, client, keyPair)

	tools.PrintResource(t, keyPair)

	allPages, err := keypairs.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allKeys, err := keypairs.ExtractKeyPairs(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, kp := range allKeys {
		tools.PrintResource(t, kp)

		if kp.Name == keyPair.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestKeyPairsImportPublicKey(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	publicKey, err := createKey()
	th.AssertNoErr(t, err)

	keyPair, err := ImportPublicKey(t, client, publicKey)
	th.AssertNoErr(t, err)
	defer DeleteKeyPair(t, client, keyPair)

	tools.PrintResource(t, keyPair)
}

func TestKeyPairsServerCreateWithKey(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	publicKey, err := createKey()
	th.AssertNoErr(t, err)

	keyPair, err := ImportPublicKey(t, client, publicKey)
	th.AssertNoErr(t, err)
	defer DeleteKeyPair(t, client, keyPair)

	server, err := CreateServerWithPublicKey(t, client, keyPair.Name)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	server, err = servers.Get(context.TODO(), client, server.ID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, server.KeyName, keyPair.Name)
}

func TestKeyPairsCreateDeleteByID(t *testing.T) {
	clients.RequireAdmin(t)

	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	computeClient.Microversion = "2.10"

	user, err := identity.CreateUser(t, identityClient, nil)
	th.AssertNoErr(t, err)
	defer identity.DeleteUser(t, identityClient, user.ID)

	keyPairName := tools.RandomString("keypair_", 5)
	createOpts := keypairs.CreateOpts{
		Name:   keyPairName,
		UserID: user.ID,
	}

	keyPair, err := keypairs.Create(context.TODO(), computeClient, createOpts).Extract()
	th.AssertNoErr(t, err)

	getOpts := keypairs.GetOpts{
		UserID: user.ID,
	}

	newKeyPair, err := keypairs.Get(context.TODO(), computeClient, keyPair.Name, getOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, keyPair.Name, newKeyPair.Name)

	listOpts := keypairs.ListOpts{
		UserID: user.ID,
	}

	allPages, err := keypairs.List(computeClient, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allKeys, err := keypairs.ExtractKeyPairs(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, kp := range allKeys {
		if kp.Name == keyPair.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	deleteOpts := keypairs.DeleteOpts{
		UserID: user.ID,
	}

	err = keypairs.Delete(context.TODO(), computeClient, keyPair.Name, deleteOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
