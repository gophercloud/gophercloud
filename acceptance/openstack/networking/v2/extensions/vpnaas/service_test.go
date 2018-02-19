// +build acceptance networking fwaas

package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	layer3 "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/acceptance/tools"
)

func TestServiceCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router.ID)

	service, err := CreateService(t, client, router.ID)
	if err != nil {
		t.Fatalf("Unable to create service: %v", err)
	}
	defer DeleteService(t, client, service.ID)

	tools.PrintResource(t, service)
}
