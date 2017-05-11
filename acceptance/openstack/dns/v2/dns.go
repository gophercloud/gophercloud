package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
)

// CreateZone will create a Zone with a random name. An error will
// be returned if the zone was unable to be created.
func CreateZone(t *testing.T, client *gophercloud.ServiceClient) (*zones.Zone, error) {
	zoneName := tools.RandomString("ACPTTEST", 8) + ".com."

	t.Logf("Attempting to create zone: %s", zoneName)
	createOpts := zones.CreateOpts{
		Name:        zoneName,
		Email:       "root@example.com",
		Type:        "PRIMARY",
		TTL:         7200,
		Description: "Test zone",
	}

	zone, err := zones.Create(client, createOpts).Extract()
	if err != nil {
		return zone, err
	}

	if err := WaitForZoneStatus(client, zone, "ACTIVE"); err != nil {
		return zone, err
	}

	newZone, err := zones.Get(client, zone.ID).Extract()
	if err != nil {
		return zone, err
	}

	t.Logf("Created Zone: %s", zoneName)
	return newZone, nil
}

// CreateSecondaryZone will create a Zone with a random name. An error will
// be returned if the zone was unable to be created.
//
// This is only for example purposes as it will try to do a zone transfer.
func CreateSecondaryZone(t *testing.T, client *gophercloud.ServiceClient) (*zones.Zone, error) {
	zoneName := tools.RandomString("ACPTTEST", 8) + ".com."

	t.Logf("Attempting to create zone: %s", zoneName)
	createOpts := zones.CreateOpts{
		Name:    zoneName,
		Type:    "SECONDARY",
		Masters: []string{"10.0.0.1"},
	}

	zone, err := zones.Create(client, createOpts).Extract()
	if err != nil {
		return zone, err
	}

	if err := WaitForZoneStatus(client, zone, "ACTIVE"); err != nil {
		return zone, err
	}

	newZone, err := zones.Get(client, zone.ID).Extract()
	if err != nil {
		return zone, err
	}

	t.Logf("Created Zone: %s", zoneName)
	return newZone, nil
}

// DeleteZone will delete a specified zone. A fatal error will occur if
// the zone failed to be deleted. This works best when used as a deferred
// function.
func DeleteZone(t *testing.T, client *gophercloud.ServiceClient, zone *zones.Zone) {
	_, err := zones.Delete(client, zone.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to delete zone %s: %v", zone.ID, err)
	}

	t.Logf("Deleted zone: %s", zone.ID)
}

// WaitForZoneStatus will poll a zone's status until it either matches
// the specified status or the status becomes ERROR.
func WaitForZoneStatus(client *gophercloud.ServiceClient, zone *zones.Zone, status string) error {
	return gophercloud.WaitFor(60, func() (bool, error) {
		current, err := zones.Get(client, zone.ID).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
