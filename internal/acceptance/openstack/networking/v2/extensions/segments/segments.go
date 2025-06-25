package segments

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func CreateSegment(t *testing.T, client *gophercloud.ServiceClient, networkID string) (*segments.Segment, error) {
	name := tools.RandomString("TESTACC-SEGMENT-", 8)
	desc := "test segment description"

	opts := segments.CreateOpts{
		NetworkID:   networkID,
		NetworkType: "geneve",
		Name:        name,
		Description: desc,
	}

	segment, err := segments.Create(context.TODO(), client, opts).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, segment)

	th.AssertEquals(t, segment.Name, name)
	th.AssertEquals(t, segment.Description, desc)
	th.AssertEquals(t, segment.NetworkType, "geneve")
	th.AssertEquals(t, segment.NetworkID, networkID)

	return segment, nil
}

func DeleteSegment(t *testing.T, client *gophercloud.ServiceClient, segmentID string) {
	t.Logf("Attempting to delete segment %s", segmentID)

	err := segments.Delete(context.TODO(), client, segmentID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete segment %s: %v", segmentID, err)
	}

	t.Logf("Deleted segment %s", segmentID)
}
