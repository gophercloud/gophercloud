// +build acceptance networking vpnaas

package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
)

func TestIPSecPolicyCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	policy, err := CreateIPSecPolicy(t, client)
	if err != nil {
		t.Fatalf("Unable to create IPSec policy: %v", err)
	}
	defer DeleteIPSecPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	newPolicy, err := ipsecpolicies.Get(client, policy.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get IPSec policy: %v", err)
	}
	tools.PrintResource(t, newPolicy)
}
