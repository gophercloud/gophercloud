// +build acceptance networking vpnaas

package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
)

func TestIKEPolicyCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	policy, err := CreateIKEPolicy(t, client)
	if err != nil {
		t.Fatalf("Unable to create IKE policy: %v", err)
	}

	tools.PrintResource(t, policy)
}
