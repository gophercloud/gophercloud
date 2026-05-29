//go:build acceptance || networking || layer3 || addressscopes

package layer3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAddressScopesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create an address-scope
	addressScope, err := CreateAddressScope(t, client)
	th.AssertNoErr(t, err)
	defer DeleteAddressScope(t, client, addressScope.ID)

	tools.PrintResource(t, addressScope)

	newName := tools.RandomString("TESTACC-", 8)
	updateOpts := &addressscopes.UpdateOpts{
		Name: &newName,
	}

	_, err = addressscopes.Update(context.TODO(), client, addressScope.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newAddressScope, err := addressscopes.Get(context.TODO(), client, addressScope.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newAddressScope)
	th.AssertEquals(t, newAddressScope.Name, newName)

	allPages, err := addressscopes.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAddressScopes, err := addressscopes.ExtractAddressScopes(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, addressScope := range allAddressScopes {
		if addressScope.ID == newAddressScope.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
