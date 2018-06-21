package v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateEmptySecret will create a random secret with no payload. An error will
// be returned if the secret could not be created.
func CreateEmptySecret(t *testing.T, client *gophercloud.ServiceClient) (*secrets.Secret, error) {
	secretName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create secret %s", secretName)

	createOpts := secrets.CreateOpts{
		Algorithm:  "aes",
		BitLength:  256,
		Mode:       "cbc",
		Name:       secretName,
		SecretType: secrets.OpaqueSecret,
	}

	secret, err := secrets.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created secret: %s", secret.SecretRef)

	secretID, err := ParseSecretID(secret.SecretRef)
	if err != nil {
		return nil, err
	}

	secret, err = secrets.Get(client, secretID).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, secret)

	th.AssertEquals(t, secret.Name, secretName)
	th.AssertEquals(t, secret.Algorithm, "aes")

	return secret, nil
}

// CreateSecretWithPayload will create a random secret with a given payload.
// An error will be returned if the secret could not be created.
func CreateSecretWithPayload(t *testing.T, client *gophercloud.ServiceClient, payload string) (*secrets.Secret, error) {
	secretName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create secret %s", secretName)

	createOpts := secrets.CreateOpts{
		Algorithm:          "aes",
		BitLength:          256,
		Mode:               "cbc",
		Name:               secretName,
		Payload:            payload,
		PayloadContentType: "text/plain",
		SecretType:         secrets.OpaqueSecret,
	}

	secret, err := secrets.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created secret: %s", secret.SecretRef)

	secretID, err := ParseSecretID(secret.SecretRef)
	if err != nil {
		return nil, err
	}

	secret, err = secrets.Get(client, secretID).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, secret)

	th.AssertEquals(t, secret.Name, secretName)
	th.AssertEquals(t, secret.Algorithm, "aes")

	return secret, nil
}

// DeleteSecret will delete a secret. A fatal error will occur if the secret
// could not be deleted. This works best when used as a deferred function.
func DeleteSecret(t *testing.T, client *gophercloud.ServiceClient, id string) {
	t.Logf("Attempting to delete secret %s", id)

	err := secrets.Delete(client, id).ExtractErr()
	if err != nil {
		t.Fatalf("Could not delete secret: %s", err)
	}

	t.Logf("Successfully deleted secret %s", id)
}

func ParseSecretID(secretRef string) (string, error) {
	parts := strings.Split(secretRef, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("Could not parse %s", secretRef)
	}

	return parts[len(parts)-1], nil
}
