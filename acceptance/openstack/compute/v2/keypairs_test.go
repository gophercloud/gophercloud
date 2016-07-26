// +build acceptance compute keypairs

package v2

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"golang.org/x/crypto/ssh"
)

const keyName = "gophercloud_test_key_pair"

func TestKeypairsList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := keypairs.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve keypairs: %s", err)
	}

	allKeys, err := keypairs.ExtractKeyPairs(allPages)
	if err != nil {
		t.Fatalf("Unable to extract keypairs results: %s", err)
	}

	for _, keypair := range allKeys {
		printKeyPair(t, &keypair)
	}
}

func TestKeypairsCreate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	createOpts := keypairs.CreateOpts{
		Name: keyName,
	}
	keyPair, err := keypairs.Create(client, createOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to create keypair: %s", err)
	}
	defer deleteKeyPair(t, client, keyPair)

	printKeyPair(t, keyPair)
}

func TestKeypairsImportPublicKey(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	publicKey, err := createKey()
	if err != nil {
		t.Fatalf("Unable to create public key: %s", err)
	}

	createOpts := keypairs.CreateOpts{
		Name:      keyName,
		PublicKey: publicKey,
	}
	keyPair, err := keypairs.Create(client, createOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to create keypair: %s", err)
	}
	defer deleteKeyPair(t, client, keyPair)

	printKeyPair(t, keyPair)
}

func TestKeypairsServerCreateWithKey(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	publicKey, err := createKey()
	if err != nil {
		t.Fatalf("Unable to create public key: %s", err)
	}

	createOpts := keypairs.CreateOpts{
		Name:      keyName,
		PublicKey: publicKey,
	}

	keyPair, err := keypairs.Create(client, createOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to create keypair: %s", err)
	}
	defer deleteKeyPair(t, client, keyPair)

	networkID, err := getNetworkIDFromTenantNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatalf("Failed to obtain network ID: %v", err)
	}

	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create server: %s", name)

	serverCreateOpts := servers.CreateOpts{
		Name:      name,
		FlavorRef: choices.FlavorID,
		ImageRef:  choices.ImageID,
		Networks: []servers.Network{
			servers.Network{UUID: networkID},
		},
	}

	server, err := servers.Create(client, keypairs.CreateOptsExt{
		serverCreateOpts,
		keyName,
	}).Extract()
	if err != nil {
		t.Fatalf("Unable to create server: %s", err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	server, err = servers.Get(client, server.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve server: %s", err)
	}

	if server.KeyName != keyName {
		t.Fatalf("key name of server %s is %s, not %s", server.ID, server.KeyName, keyName)
	}
}

func createKey() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.PublicKey
	pub, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return "", err
	}

	pubBytes := ssh.MarshalAuthorizedKey(pub)
	pk := string(pubBytes)
	return pk, nil
}

func deleteKeyPair(t *testing.T, client *gophercloud.ServiceClient, keyPair *keypairs.KeyPair) {
	err := keypairs.Delete(client, keyPair.Name).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete keypair %s: %v", keyPair.Name, err)
	}

	t.Logf("Deleted keypair: %s", keyPair.Name)
}

func printKeyPair(t *testing.T, keypair *keypairs.KeyPair) {
	t.Logf("Name: %s", keypair.Name)
	t.Logf("Fingerprint: %s", keypair.Fingerprint)
	t.Logf("Public Key: %s", keypair.PublicKey)
	t.Logf("Private Key: %s", keypair.PrivateKey)
	t.Logf("UserID: %s", keypair.UserID)
}
