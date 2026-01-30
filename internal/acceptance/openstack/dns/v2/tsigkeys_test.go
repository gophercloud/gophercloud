//go:build acceptance || dns || tsigkeys

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/tsigkeys"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTSIGKeysCRUD(t *testing.T) {
	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	tsigkey, err := CreateTSIGKey(t, client)
	th.AssertNoErr(t, err)
	defer DeleteTSIGKey(t, client, tsigkey)

	tools.PrintResource(t, &tsigkey)

	allPages, err := tsigkeys.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allTSIGKeys, err := tsigkeys.ExtractTSIGKeys(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, k := range allTSIGKeys {
		tools.PrintResource(t, &k)

		if tsigkey.Name == k.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	updateOpts := tsigkeys.UpdateOpts{
		Name:   tsigkey.Name + "-updated",
		Secret: "updated-test-secret-key==",
	}

	newTSIGKey, err := tsigkeys.Update(context.TODO(), client, tsigkey.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, &newTSIGKey)

	th.AssertEquals(t, newTSIGKey.Name, tsigkey.Name+"-updated")
	th.AssertEquals(t, newTSIGKey.Secret, "updated-test-secret-key==")
}
