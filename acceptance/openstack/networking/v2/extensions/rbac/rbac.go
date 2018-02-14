package rbac

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbac"
)

// CreateRbac will create a rbac-policy. An error will be returned if the
// rbac-policy could not be created.
func CreateRbac(t *testing.T, client *gophercloud.ServiceClient, tenantID, networkID string) (*rbac.Rbac, error) {
	createOpts := rbac.CreateOpts{
		Action:       "access_as_shared",
		ObjectType:   "network",
		TargetTenant: tenantID,
		ObjectID:     networkID,
	}

	t.Logf("Trying to create rbac_policy")

	rbacPolicy, err := rbac.Create(client, createOpts).Extract()
	if err != nil {
		return rbacPolicy, err
	}

	t.Logf("Successfully created rbac_policy")
	return rbacPolicy, nil
}
