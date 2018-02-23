// +build acceptance networking vpnaas

package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
)

func TestIPSecPolicyCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	policy, err := CreateIPSecPolicy(t, client)
	if err != nil {
		t.Fatalf("Unable to create policy: %v", err)
	}
	defer DeleteIPSecPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)
}
