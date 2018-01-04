package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
)

// CreateSubnetPool will create a subnetpool. An error will be returned if the
// subnetpool could not be created.
func CreateSubnetPool(t *testing.T, client *gophercloud.ServiceClient) (*subnetpools.SubnetPool, error) {
	subnetPoolName := tools.RandomString("TESTACC-", 8)
	subnetPoolPrefixes := []string{
		"10.0.0.0/8",
	}
	createOpts := subnetpools.CreateOpts{
		Name:     subnetPoolName,
		Prefixes: subnetPoolPrefixes,
	}

	t.Logf("Attempting to create a subnetpool: %s", subnetPoolName)

	subnetPool, err := subnetpools.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the subnetpool.")
	return subnetPool, nil
}
