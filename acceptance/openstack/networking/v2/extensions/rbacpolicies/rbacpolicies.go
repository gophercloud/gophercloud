package rbacpolicies

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbacpolicies"
)

// CreateRBAC will create a rbac-policy. An error will be returned if the
// rbac-policy could not be created.
func CreateRBAC(t *testing.T, client *gophercloud.ServiceClient, tenantID, networkID string) (*rbacpolicies.RBAC, error) {
	createOpts := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.ActionAccessShared,
		ObjectType:   "network",
		TargetTenant: tenantID,
		ObjectID:     networkID,
	}

	t.Logf("Trying to create rbac_policy")

	rbacPolicy, err := rbacpolicies.Create(client, createOpts).Extract()
	if err != nil {
		return rbacPolicy, err
	}

	t.Logf("Successfully created rbac_policy")
	return rbacPolicy, nil
}
