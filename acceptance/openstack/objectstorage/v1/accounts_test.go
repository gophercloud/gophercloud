// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/accounts"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestAccounts(t *testing.T) {
	// Create a provider client for making the HTTP requests.
	// See common.go in this directory for more information.
	client := newClient(t)

	// Update an account's metadata.
	res = accounts.Update(client, accounts.UpdateOpts{Metadata: metadata})
	th.AssertNoErr(t, res.Err)

	// Defer the deletion of the metadata set above.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		res = accounts.Update(client, accounts.UpdateOpts{Metadata: tempMap})
		th.AssertNoErr(t, res.Err)
	}()

	// Retrieve account metadata.
	res := accounts.Get(client, accounts.GetOpts{})
	th.AssertNoErr(res.Err)
	// Extract the custom metadata from the 'Get' response.
	am := accounts.ExtractMetadata(gr)
	for k := range metadata {
		if am[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
