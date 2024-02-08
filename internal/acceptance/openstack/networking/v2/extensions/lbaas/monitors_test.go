//go:build acceptance || networking || lbaas || monitors
// +build acceptance networking lbaas monitors

package lbaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/lbaas/monitors"
)

func TestMonitorsList(t *testing.T) {
	t.Skip("Neutron LBaaS was replaced by Octavia and the API will be removed in a future release")
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := monitors.List(client, monitors.ListOpts{}).AllPages()
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

func TestMonitorsCRUD(t *testing.T) {
	t.Skip("Neutron LBaaS was replaced by Octavia and the API will be removed in a future release")
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	monitor, err := CreateMonitor(t, client)
	if err != nil {
		t.Fatalf("Unable to create monitor: %v", err)
	}
	defer DeleteMonitor(t, client, monitor.ID)

	tools.PrintResource(t, monitor)

	updateOpts := monitors.UpdateOpts{
		Delay: 999,
	}

	_, err = monitors.Update(client, monitor.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update monitor: %v", err)
	}

	newMonitor, err := monitors.Get(client, monitor.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get monitor: %v", err)
	}

	tools.PrintResource(t, newMonitor)
}
