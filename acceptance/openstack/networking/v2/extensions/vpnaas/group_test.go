// +build acceptance networking vpnaas

package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
)

func TestGroupCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	group, err := CreateEndpointGroup(t, client)
	if err != nil {
		t.Fatalf("Unable to create Endpoint group: %v", err)
	}
	tools.PrintResource(t, group)

	newGroup, err := endpointgroups.Get(client, group.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve Endpoint group: %v", err)
	}
	tools.PrintResource(t, newGroup)

}
