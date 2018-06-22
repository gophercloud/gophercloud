// +build acceptance keymanager containers

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGenericContainersCRUD(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	payload := tools.RandomString("SUPERSECRET-", 8)
	secret, err := CreateSecretWithPayload(t, client, payload)
	th.AssertNoErr(t, err)
	secretID, err := ParseID(secret.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	container, err := CreateGenericContainer(t, client, secret)
	th.AssertNoErr(t, err)
	containerID, err := ParseID(container.ContainerRef)
	th.AssertNoErr(t, err)
	defer DeleteContainer(t, client, containerID)

	allPages, err := containers.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allContainers, err := containers.ExtractContainers(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allContainers {
		if v.ContainerRef == container.ContainerRef {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestCertificateContainer(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	pass := tools.RandomString("", 16)
	priv, cert, err := CreateCertificate(t, pass)
	th.AssertNoErr(t, err)

	private, err := CreatePrivateSecret(t, client, priv)
	th.AssertNoErr(t, err)
	secretID, err := ParseID(private.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload, err := secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Private Payload: %s", string(payload))

	certificate, err := CreateCertificateSecret(t, client, cert)
	th.AssertNoErr(t, err)
	secretID, err = ParseID(certificate.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload, err = secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Certificate Payload: %s", string(payload))

	passphrase, err := CreatePassphraseSecret(t, client, pass)
	th.AssertNoErr(t, err)
	secretID, err = ParseID(passphrase.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload, err = secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Passphrase Payload: %s", string(payload))

	container, err := CreateCertificateContainer(t, client, passphrase, private, certificate)
	th.AssertNoErr(t, err)
	containerID, err := ParseID(container.ContainerRef)
	defer DeleteContainer(t, client, containerID)
}

func TestRSAContainer(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	pass := tools.RandomString("", 16)
	priv, pub, err := CreateRSAKeyPair(t, pass)
	th.AssertNoErr(t, err)

	private, err := CreatePrivateSecret(t, client, priv)
	th.AssertNoErr(t, err)
	secretID, err := ParseID(private.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload, err := secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Private Payload: %s", string(payload))

	public, err := CreatePublicSecret(t, client, pub)
	th.AssertNoErr(t, err)
	secretID, err = ParseID(public.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload, err = secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Public Payload: %s", string(payload))

	passphrase, err := CreatePassphraseSecret(t, client, pass)
	th.AssertNoErr(t, err)
	secretID, err = ParseID(passphrase.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload, err = secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Passphrase Payload: %s", string(payload))

	container, err := CreateRSAContainer(t, client, passphrase, private, public)
	th.AssertNoErr(t, err)
	containerID, err := ParseID(container.ContainerRef)
	defer DeleteContainer(t, client, containerID)
}
