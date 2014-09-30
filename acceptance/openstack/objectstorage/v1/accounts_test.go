// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/accounts"
)

func TestAccounts(t *testing.T) {
	// Create a provider client for making the HTTP requests.
	// See common.go in this directory for more information.
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	// Update an account's metadata.
	err = accounts.Update(client, accounts.UpdateOpts{
		Metadata: metadata,
	})
	if err != nil {
		t.Error(err)
		return
	}
	// Defer the deletion of the metadata set above.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		err = accounts.Update(client, accounts.UpdateOpts{
			Metadata: tempMap,
		})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// Retrieve account metadata.
	gr, err := accounts.Get(client, accounts.GetOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	// Extract the custom metadata from the 'Get' response.
	am := accounts.ExtractMetadata(gr)
	for k := range metadata {
		if am[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
