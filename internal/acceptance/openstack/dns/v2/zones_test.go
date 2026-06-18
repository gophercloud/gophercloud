//go:build acceptance || dns || zones

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestZonesCRUD(t *testing.T) {
	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	tools.PrintResource(t, &zone)

	allPages, err := zones.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allZones, err := zones.ExtractZones(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, z := range allZones {
		tools.PrintResource(t, &z)

		if zone.Name == z.Name {
			found = true
		}
	}

	th.AssertTrue(t, found)

	description := ""
	updateOpts := zones.UpdateOpts{
		Description: &description,
		TTL:         3600,
	}

	newZone, err := zones.Update(context.TODO(), client, zone.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, &newZone)

	th.AssertEquals(t, newZone.Description, description)
	th.AssertEquals(t, 3600, newZone.TTL)
}

func TestZonesListWithAllProjects(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	listOpts := zones.ListOpts{AllProjects: true}
	allPages, err := zones.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allZones, err := zones.ExtractZones(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, z := range allZones {
		if zone.Name == z.Name {
			found = true
		}
	}

	th.AssertTrue(t, found)
}
