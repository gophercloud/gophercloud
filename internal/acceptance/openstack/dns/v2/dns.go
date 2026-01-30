package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/recordsets"
	transferAccepts "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/accept"
	transferRequests "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/request"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/tsigkeys"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateRecordSet will create a RecordSet with a random name. An error will
// be returned if the zone was unable to be created.
func CreateRecordSet(t *testing.T, client *gophercloud.ServiceClient, zone *zones.Zone) (*recordsets.RecordSet, error) {
	t.Logf("Attempting to create recordset: %s", zone.Name)

	createOpts := recordsets.CreateOpts{
		Name:        zone.Name,
		Type:        "A",
		TTL:         3600,
		Description: "Test recordset",
		Records:     []string{"10.1.0.2"},
	}

	rs, err := recordsets.Create(context.TODO(), client, zone.ID, createOpts).Extract()
	if err != nil {
		return rs, err
	}

	if err := WaitForRecordSetStatus(client, rs, "ACTIVE"); err != nil {
		return rs, err
	}

	newRS, err := recordsets.Get(context.TODO(), client, rs.ZoneID, rs.ID).Extract()
	if err != nil {
		return newRS, err
	}

	t.Logf("Created record set: %s", newRS.Name)

	th.AssertEquals(t, newRS.Name, zone.Name)

	return rs, nil
}

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

	zone, err := zones.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return zone, err
	}

	if err := WaitForZoneStatus(client, zone, "ACTIVE"); err != nil {
		return zone, err
	}

	newZone, err := zones.Get(context.TODO(), client, zone.ID).Extract()
	if err != nil {
		return zone, err
	}

	t.Logf("Created Zone: %s", zoneName)

	th.AssertEquals(t, newZone.Name, zoneName)
	th.AssertEquals(t, newZone.TTL, 7200)

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

	zone, err := zones.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return zone, err
	}

	if err := WaitForZoneStatus(client, zone, "ACTIVE"); err != nil {
		return zone, err
	}

	newZone, err := zones.Get(context.TODO(), client, zone.ID).Extract()
	if err != nil {
		return zone, err
	}

	t.Logf("Created Zone: %s", zoneName)

	th.AssertEquals(t, newZone.Name, zoneName)
	th.AssertEquals(t, newZone.Masters[0], "10.0.0.1")

	return newZone, nil
}

// CreateTransferRequest will create a Transfer Request to a spectified Zone. An error will
// be returned if the zone transfer request was unable to be created.
func CreateTransferRequest(t *testing.T, client *gophercloud.ServiceClient, zone *zones.Zone, targetProjectID string) (*transferRequests.TransferRequest, error) {
	t.Logf("Attempting to create Transfer Request to Zone: %s", zone.Name)

	createOpts := transferRequests.CreateOpts{
		TargetProjectID: targetProjectID,
		Description:     "Test transfer request",
	}

	transferRequest, err := transferRequests.Create(context.TODO(), client, zone.ID, createOpts).Extract()
	if err != nil {
		return transferRequest, err
	}

	if err := WaitForTransferRequestStatus(client, transferRequest, "ACTIVE"); err != nil {
		return transferRequest, err
	}

	newTransferRequest, err := transferRequests.Get(context.TODO(), client, transferRequest.ID).Extract()
	if err != nil {
		return transferRequest, err
	}

	t.Logf("Created Transfer Request for Zone: %s", zone.Name)

	th.AssertEquals(t, newTransferRequest.ZoneID, zone.ID)
	th.AssertEquals(t, newTransferRequest.ZoneName, zone.Name)

	return newTransferRequest, nil
}

// CreateTransferAccept will accept a spectified Transfer Request. An error will
// be returned if the zone transfer accept was unable to be created.
func CreateTransferAccept(t *testing.T, client *gophercloud.ServiceClient, zoneTransferRequestID string, key string) (*transferAccepts.TransferAccept, error) {
	t.Logf("Attempting to accept specified transfer reqeust: %s", zoneTransferRequestID)
	createOpts := transferAccepts.CreateOpts{
		ZoneTransferRequestID: zoneTransferRequestID,
		Key:                   key,
	}
	transferAccept, err := transferAccepts.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return transferAccept, err
	}
	if err := WaitForTransferAcceptStatus(client, transferAccept, "COMPLETE"); err != nil {
		return transferAccept, err
	}
	newTransferAccept, err := transferAccepts.Get(context.TODO(), client, transferAccept.ID).Extract()
	if err != nil {
		return transferAccept, err
	}
	t.Logf("Accepted Transfer Request: %s", zoneTransferRequestID)
	th.AssertEquals(t, newTransferAccept.ZoneTransferRequestID, zoneTransferRequestID)
	return newTransferAccept, nil
}

// DeleteTransferRequest will delete a specified zone transfer request. A fatal error will occur if
// the transfer request failed to be deleted. This works best when used as a deferred
// function.
func DeleteTransferRequest(t *testing.T, client *gophercloud.ServiceClient, tr *transferRequests.TransferRequest) {
	err := transferRequests.Delete(context.TODO(), client, tr.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete zone transfer request %s: %v", tr.ID, err)
	}
	t.Logf("Deleted zone transfer request: %s", tr.ID)
}

// CreateShare will create a zone share. An error will be returned if the
// zone share was unable to be created.
func CreateShare(t *testing.T, client *gophercloud.ServiceClient, zone *zones.Zone, targetProjectID string) (*zones.ZoneShare, error) {
	t.Logf("Attempting to share zone %s with project %s", zone.ID, targetProjectID)

	createOpts := zones.ShareZoneOpts{
		TargetProjectID: targetProjectID,
	}

	share, err := zones.Share(context.TODO(), client, zone.ID, createOpts).Extract()
	if err != nil {
		return share, err
	}

	t.Logf("Created share for zone: %s", zone.ID)

	th.AssertEquals(t, share.ZoneID, zone.ID)
	th.AssertEquals(t, share.TargetProjectID, targetProjectID)

	return share, nil
}

// UnshareZone will unshare a zone. An error will be returned if the
// zone unshare was unable to be created.
func UnshareZone(t *testing.T, client *gophercloud.ServiceClient, share *zones.ZoneShare) {
	t.Logf("Attempting to unshare zone %s with project %s", share.ZoneID, share.TargetProjectID)

	err := zones.Unshare(context.TODO(), client, share.ZoneID, share.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to unshare zone %s: %v", share.ZoneID, err)
	}

	t.Logf("Unshared zone: %s", share.ZoneID)
}

// DeleteRecordSet will delete a specified record set. A fatal error will occur if
// the record set failed to be deleted. This works best when used as a deferred
// function.
func DeleteRecordSet(t *testing.T, client *gophercloud.ServiceClient, rs *recordsets.RecordSet) {
	err := recordsets.Delete(context.TODO(), client, rs.ZoneID, rs.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete record set %s: %v", rs.ID, err)
	}

	t.Logf("Deleted record set: %s", rs.ID)
}

// DeleteZone will delete a specified zone. A fatal error will occur if
// the zone failed to be deleted. This works best when used as a deferred
// function.
func DeleteZone(t *testing.T, client *gophercloud.ServiceClient, zone *zones.Zone) {
	_, err := zones.Delete(context.TODO(), client, zone.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to delete zone %s: %v", zone.ID, err)
	}

	t.Logf("Deleted zone: %s", zone.ID)
}

// WaitForRecordSetStatus will poll a record set's status until it either matches
// the specified status or the status becomes ERROR.
func WaitForRecordSetStatus(client *gophercloud.ServiceClient, rs *recordsets.RecordSet, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		current, err := recordsets.Get(ctx, client, rs.ZoneID, rs.ID).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}

// WaitForTransferRequestStatus will poll a transfer reqeust's status until it either matches
// the specified status or the status becomes ERROR.
func WaitForTransferRequestStatus(client *gophercloud.ServiceClient, tr *transferRequests.TransferRequest, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		current, err := transferRequests.Get(ctx, client, tr.ID).Extract()
		if err != nil {
			return false, err
		}
		if current.Status == status {
			return true, nil
		}
		return false, nil
	})
}

// WaitForTransferAcceptStatus will poll a transfer accept's status until it either matches
// the specified status or the status becomes ERROR.
func WaitForTransferAcceptStatus(client *gophercloud.ServiceClient, ta *transferAccepts.TransferAccept, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		current, err := transferAccepts.Get(ctx, client, ta.ID).Extract()
		if err != nil {
			return false, err
		}
		if current.Status == status {
			return true, nil
		}
		return false, nil
	})
}

// WaitForZoneStatus will poll a zone's status until it either matches
// the specified status or the status becomes ERROR.
func WaitForZoneStatus(client *gophercloud.ServiceClient, zone *zones.Zone, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		current, err := zones.Get(ctx, client, zone.ID).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}

// CreateTSIGKey will create a TSIG key with a random name. An error will
// be returned if the TSIG key was unable to be created.
func CreateTSIGKey(t *testing.T, client *gophercloud.ServiceClient) (*tsigkeys.TSIGKey, error) {
	keyName := tools.RandomString("ACPTTEST", 8)

	t.Logf("Attempting to create TSIG key: %s", keyName)
	createOpts := tsigkeys.CreateOpts{
		Name:      keyName,
		Algorithm: "hmac-sha256",
		Secret:    "example-test-secret-key==",
		Scope:     "POOL",
	}

	tsigkey, err := tsigkeys.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return tsigkey, err
	}

	newTSIGKey, err := tsigkeys.Get(context.TODO(), client, tsigkey.ID).Extract()
	if err != nil {
		return tsigkey, err
	}

	t.Logf("Created TSIG key: %s", keyName)

	th.AssertEquals(t, newTSIGKey.Name, keyName)
	th.AssertEquals(t, newTSIGKey.Algorithm, "hmac-sha256")

	return newTSIGKey, nil
}

// DeleteTSIGKey will delete a specified TSIG key. A fatal error will occur if
// the TSIG key failed to be deleted. This works best when used as a deferred
// function.
func DeleteTSIGKey(t *testing.T, client *gophercloud.ServiceClient, tsigkey *tsigkeys.TSIGKey) {
	err := tsigkeys.Delete(context.TODO(), client, tsigkey.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete TSIG key %s: %v", tsigkey.ID, err)
	}

	t.Logf("Deleted TSIG key: %s", tsigkey.ID)
}
