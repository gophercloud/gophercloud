//go:build acceptance || networking || loadbalancer || monitors
// +build acceptance networking loadbalancer monitors

package lbaas_v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
)

func TestMonitorsList(t *testing.T) {
	t.Skip("Neutron LBaaS v2 was replaced by Octavia and the API will be removed in a future release")
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := monitors.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list monitors: %v", err)
	}

	allMonitors, err := monitors.ExtractMonitors(allPages)
	if err != nil {
		t.Fatalf("Unable to extract monitors: %v", err)
	}

	for _, monitor := range allMonitors {
		tools.PrintResource(t, monitor)
	}
}
