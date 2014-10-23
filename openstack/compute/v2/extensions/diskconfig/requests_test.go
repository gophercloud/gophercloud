package diskconfig

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestCreateOpts(t *testing.T) {
	base := servers.CreateOpts{
		Name:      "createdserver",
		ImageRef:  "asdfasdfasdf",
		FlavorRef: "performance1-1",
	}

	ext := CreateOptsExt{
		CreateOptsBuilder: base,
		DiskConfig:        Manual,
	}

	expected := `
		{
			"server": {
				"name": "createdserver",
				"imageRef": "asdfasdfasdf",
				"flavorRef": "performance1-1",
				"OS-DCF:diskConfig": "MANUAL"
			}
		}
	`
	th.CheckJSONEquals(t, expected, ext.ToServerCreateMap())
}

func TestRebuildOpts(t *testing.T) {
	base := servers.RebuildOpts{
		Name:      "createdserver",
		AdminPass: "swordfish",
		ImageID:   "asdfasdfasdf",
	}

	ext := RebuildOptsExt{
		RebuildOptsBuilder: base,
		DiskConfig:         Auto,
	}

	actual, err := ext.ToServerRebuildMap()
	th.AssertNoErr(t, err)

	expected := `
		{
			"rebuild": {
				"name": "createdserver",
				"imageRef": "asdfasdfasdf",
				"adminPass": "swordfish",
				"OS-DCF:diskConfig": "AUTO"
			}
		}
	`
	th.CheckJSONEquals(t, expected, actual)
}
