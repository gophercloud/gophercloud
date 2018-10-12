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
		// Note Neutron returns tags sorted, and although the API
		// docs say list of tags, it's a set e.g no duplicates
		Tags: []string{"a", "b", "c"},
	}
	tags, err := attributestags.ReplaceAll(client, "networks", network.ID, tagReplaceAllOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, tags, []string{"a", "b", "c"})

	// Add a tag
	err = attributestags.Add(client, "networks", network.ID, "d").ExtractErr()
	th.AssertNoErr(t, err)

	// Verify the tags are set in the List response
	tags, err = attributestags.List(client, "networks", network.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, tags, []string{"a", "b", "c", "d"})

	// Delete all tags
	err = attributestags.DeleteAll(client, "networks", network.ID).ExtractErr()
	th.AssertNoErr(t, err)
	tags, err = attributestags.List(client, "networks", network.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, tags, []string{})
}
