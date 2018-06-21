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

func TestSecretsMetadataCRUD(t *testing.T) {
	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	payload := tools.RandomString("SUPERSECRET-", 8)
	secret, err := CreateSecretWithPayload(t, client, payload)
	th.AssertNoErr(t, err)
	secretID, err := ParseSecretID(secret.SecretRef)
	th.AssertNoErr(t, err)
	defer DeleteSecret(t, client, secretID)

	// Create some metadata
	createOpts := secrets.MetadataOpts{
		"foo":       "bar",
		"something": "something else",
	}

	ref, err := secrets.CreateMetadata(client, secretID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, ref["metadata_ref"], secret.SecretRef+"/metadata")

	// Get the metadata
	metadata, err := secrets.GetMetadata(client, secretID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, metadata)
	th.AssertEquals(t, metadata["foo"], "bar")
	th.AssertEquals(t, metadata["something"], "something else")

	// Add a single metadatum
	metadatumOpts := secrets.MetadatumOpts{
		Key:   "bar",
		Value: "baz",
	}

	err = secrets.CreateMetadatum(client, secretID, metadatumOpts).ExtractErr()
	th.AssertNoErr(t, err)

	metadata, err = secrets.GetMetadata(client, secretID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, metadata)
	th.AssertEquals(t, len(metadata), 3)
	th.AssertEquals(t, metadata["foo"], "bar")
	th.AssertEquals(t, metadata["something"], "something else")
	th.AssertEquals(t, metadata["bar"], "baz")

	// Update a metadatum
	metadatumOpts.Key = "foo"
	metadatumOpts.Value = "foo"

	metadatum, err := secrets.UpdateMetadatum(client, secretID, metadatumOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, metadatum)
	th.AssertDeepEquals(t, metadatum.Key, "foo")
	th.AssertDeepEquals(t, metadatum.Value, "foo")

	metadata, err = secrets.GetMetadata(client, secretID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, metadata)
	th.AssertEquals(t, len(metadata), 3)
	th.AssertEquals(t, metadata["foo"], "foo")
	th.AssertEquals(t, metadata["something"], "something else")
	th.AssertEquals(t, metadata["bar"], "baz")

	// Delete a metadatum
	err = secrets.DeleteMetadatum(client, secretID, "foo").ExtractErr()
	th.AssertNoErr(t, err)

	metadata, err = secrets.GetMetadata(client, secretID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, metadata)
	th.AssertEquals(t, len(metadata), 2)
	th.AssertEquals(t, metadata["something"], "something else")
	th.AssertEquals(t, metadata["bar"], "baz")
}
