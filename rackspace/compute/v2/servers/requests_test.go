package servers

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/diskconfig"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestCreateOpts(t *testing.T) {
	opts := CreateOpts{
		Name:       "createdserver",
		ImageRef:   "image-id",
		FlavorRef:  "flavor-id",
		KeyPair:    "mykey",
		DiskConfig: diskconfig.Manual,
	}

	expected := `
	{
		"server": {
			"name": "createdserver",
			"imageRef": "image-id",
			"flavorRef": "flavor-id",
			"key_name": "mykey",
			"OS-DCF:diskConfig": "MANUAL"
		}
	}
	`
	th.CheckJSONEquals(t, expected, opts.ToServerCreateMap())
}
