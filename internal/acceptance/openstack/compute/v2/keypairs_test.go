//go:build acceptance || compute || keypairs

package v2

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	identity "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/identity/v3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// keyPairsPaginationOpts supplies the page size missing from keypairs.ListOpts.
type keyPairsPaginationOpts struct{}

func (keyPairsPaginationOpts) ToKeyPairListQuery() (string, error) {
	return "?limit=1", nil
}

func TestKeyPairsPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.35"

	expected := make(map[string]int)
	for range 2 {
		keyPair, err := keypairs.Create(ctx, client, keypairs.CreateOpts{
			Name: tools.RandomString("keypair_pagination_", 16),
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cleanupCancel()
			if err := keypairs.Delete(cleanupCtx, client, keyPair.Name, nil).ExtractErr(); err != nil {
				t.Errorf("Unable to delete keypair %s: %v", keyPair.Name, err)
			}
		})
		expected[keyPair.Name] = 1
		t.Logf("Created keypair %s", keyPair.Name)
	}

	t.Run("EachPage", func(t *testing.T) {
		actual := make(map[string]int)
		pages := 0
		err := keypairs.List(client, keyPairsPaginationOpts{}).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			items, err := keypairs.ExtractKeyPairs(page)
			if err != nil {
				return false, err
			}
			pages++
			th.CheckEquals(t, 1, len(items))
			for _, item := range items {
				if _, ok := expected[item.Name]; ok {
					actual[item.Name]++
				}
			}
			return true, nil
		})
		th.AssertNoErr(t, err)
		t.Logf("Visited %d non-empty pages", pages)
		th.CheckTrue(t, pages >= 2)
		th.CheckDeepEquals(t, expected, actual)
	})

	t.Run("AllPages", func(t *testing.T) {
		page, err := keypairs.List(client, keyPairsPaginationOpts{}).AllPages(ctx)
		th.AssertNoErr(t, err)
		items, err := keypairs.ExtractKeyPairs(page)
		th.AssertNoErr(t, err)
		actual := make(map[string]int)
		for _, item := range items {
			if _, ok := expected[item.Name]; ok {
				actual[item.Name]++
			}
		}
		th.CheckDeepEquals(t, expected, actual)
	})
}

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

	th.AssertTrue(t, found)
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

	th.AssertEquals(t, keyPair.Name, server.KeyName)
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

	th.AssertTrue(t, found)

	deleteOpts := keypairs.DeleteOpts{
		UserID: user.ID,
	}

	err = keypairs.Delete(context.TODO(), computeClient, keyPair.Name, deleteOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
