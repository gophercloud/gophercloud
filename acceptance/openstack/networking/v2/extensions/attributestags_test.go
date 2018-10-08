// +build acceptance networking tags

package extensions

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestTags(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	tagReplaceAllOpts := attributestags.ReplaceAllOpts{
		Tags: []string{"abc", "123"},
	}
	tags, err := attributestags.ReplaceAll(client, "networks", network.ID, tagReplaceAllOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, tags, []string{"abc", "123"})

	// FIXME(shardy) - when the networks Get schema supports tags we should
	// verify the tags are set in the Get response
}
