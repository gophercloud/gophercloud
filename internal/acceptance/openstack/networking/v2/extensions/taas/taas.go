package taas

import (
	"context"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/taas/tapmirrors"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateTapMirror will create a Tap Mirror with the specified portID and remoteIP. An error
// will be returned if the Tap Mirror could not be created.
func CreateTapMirror(t *testing.T, client *gophercloud.ServiceClient, portID string, remoteIP string) (*tapmirrors.TapMirror, error) {
	mirrorName := tools.RandomString("TESTACC-", 8)
	mirrorDescription := tools.RandomString("TESTACC-DESC-", 8)
	mirrorDirectionIN := tools.RandomInt(1, 1000000)
	t.Logf("Attempting to create tap mirror: %s", mirrorName)

	createopts := tapmirrors.CreateOpts{
		Name:        mirrorName,
		Description: mirrorDescription,
		PortID:      portID,
		MirrorType:  tapmirrors.MirrorTypeErspanv1,
		RemoteIP:    remoteIP,
		Directions: tapmirrors.Directions{
			In:  strconv.Itoa(mirrorDirectionIN),
			Out: strconv.Itoa(mirrorDirectionIN + 1),
		},
	}

	mirror, err := tapmirrors.Create(context.TODO(), client, createopts).Extract()
	if err != nil {
		return nil, err
	}

	th.AssertEquals(t, mirrorName, mirror.Name)
	th.AssertEquals(t, mirrorDescription, mirror.Description)

	t.Logf("Created Tap Mirror: %s", mirror.ID)
	return mirror, nil
}

// DeleteTapMirror will delete a Tap Mirror with a specified ID. A fatal error will
// occur if the delete was not successful. This works best when used as a
// deferred function.
func DeleteTapMirror(t *testing.T, client *gophercloud.ServiceClient, mirrorID string) {
	t.Logf("Attempting to delete Tap Mirror: %s", mirrorID)

	err := tapmirrors.Delete(context.TODO(), client, mirrorID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete Tap Mirror %s: %v", mirrorID, err)
	}

	t.Logf("Deleted Tap Mirror: %s", mirrorID)
}
