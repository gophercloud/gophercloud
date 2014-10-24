// +build acceptance

package v2

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/smashwilson/gophercloud/acceptance/tools"
)

func TestBootFromVolume(t *testing.T) {
	client, err := newClient()
	th.AssertNoErr(t, err)

	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	options, err := optionsFromEnv()
	th.AssertNoErr(t, err)

	name := tools.RandomString("Gophercloud-", 8)
	t.Logf("Creating server [%s].", name)

	bd := bootfromvolume.BlockDevice{
		UUID:       options.imageID,
		SourceType: "image",
	}

	server, err := bootfromvolume.Create(client, bootfromvolume.CreateOptsExt{
		Name:        name,
		BlockDevice: bd,
	}).Extract()
	th.AssertNoErr(t, err)
	//defer deleteServer(t, client, server)
}
