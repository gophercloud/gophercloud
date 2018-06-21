// +build acceptance clustering policies

package v1

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestSecretsCRUD(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	payload := tools.RandomString("SUPERSECRET-", 8)
	secret, err := CreateSecretWithPayload(t, client, payload)
	th.AssertNoErr(t, err)
	secretID, err := ParseSecretID(secret.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	// Test payload retrieval
	actual, err := secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, payload, string(actual))

	// Test listing secrets
	createdQuery := &secrets.DateQuery{
		Date:   time.Date(2049, 6, 7, 1, 2, 3, 0, time.UTC),
		Filter: secrets.DateFilterLT,
	}

	listOpts := secrets.ListOpts{
		CreatedQuery: createdQuery,
	}

	allPages, err := secrets.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allSecrets, err := secrets.ExtractSecrets(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allSecrets {
		if v.SecretRef == secret.SecretRef {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestSecretsDelayedPayload(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	secret, err := CreateEmptySecret(t, client)
	th.AssertNoErr(t, err)
	secretID, err := ParseSecretID(secret.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	payload := tools.RandomString("SUPERSECRET-", 8)
	updateOpts := secrets.UpdateOpts{
		ContentType: "text/plain",
		Payload:     payload,
	}

	err = secrets.Update(client, secretID, updateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	// Test payload retrieval
	actual, err := secrets.GetPayload(client, secretID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, payload, string(actual))
}
