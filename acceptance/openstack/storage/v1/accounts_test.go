// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/openstack/storage/v1/accounts"
)

func TestAccounts(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	err = accounts.Update(client, accounts.UpdateOpts{
		Metadata: metadata,
	})
	if err != nil {
		t.Error(err)
		return
	}
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

	gr, err := accounts.Get(client, accounts.GetOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	am := accounts.ExtractMetadata(gr)
	for k := range metadata {
		if am[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
