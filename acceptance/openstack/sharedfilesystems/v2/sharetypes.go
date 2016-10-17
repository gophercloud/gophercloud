package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharetypes"
)

// CreateShareType will create a share type with a random name. An
// error will be returned if the share type was unable to be created.
func CreateShareType(t *testing.T, client *gophercloud.ServiceClient) (*sharetypes.ShareType, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires share type creation in short mode.")
	}

	shareTypeName := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create share type: %s", shareTypeName)

	extraSpecsOps := sharetypes.ExtraSpecsOpts{
		DriverHandlesShareServers: true,
	}

	createOpts := sharetypes.CreateOpts{
		Name: shareTypeName,
		OSShareTypeAccessIsPublic: true,
		ExtraSpecs:                extraSpecsOps,
	}

	shareType, err := sharetypes.Create(client, createOpts).Extract()
	if err != nil {
		return shareType, err
	}

	return shareType, nil
}

// PrintShareType will print a share type and all of its attributes.
func PrintShareType(t *testing.T, shareType *sharetypes.ShareType) {
	t.Logf("Name: %s", shareType.Name)
	t.Logf("ID: %s", shareType.ID)
	t.Logf("OS share type access is public: %t", shareType.OSShareTypeAccessIsPublic)
	t.Logf("Extra specs: %#v", shareType.ExtraSpecs)
}
