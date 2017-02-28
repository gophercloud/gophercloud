// +build acceptance dns zones
package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestZonesList(t *testing.T) {
	client, err := clients.NewDNSV2Client()
	if err != nil {
		t.Fatalf("Unable to create a DNS client: %v", err)
	}

	// This works
	var allZones []zones.Zone
	pager := zones.List(client)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		zones, err := zones.ExtractZones(page)
		if err != nil {
			return false, err
		}

		for _, i := range zones {
			allZones = append(allZones, i)
		}

		return true, nil
	})

	// This doesn't
	/*
		allPages, err := zones.List(client).AllPages()
		if err != nil {
			t.Fatalf("Unable to retrieve zones: %v", err)
		}

		allZones, err = zones.ExtractZones(allPages)
		if err != nil {
			t.Fatalf("Unable to extract zones: %v", err)
		}
	*/

	for _, zone := range allZones {
		tools.PrintResource(t, &zone)
	}
}
