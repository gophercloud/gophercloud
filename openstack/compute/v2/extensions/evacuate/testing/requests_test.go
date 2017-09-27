package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/evacuate"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"

func TestEvacuate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockMigrateResponse(t, serverID)

	err := evacuate.Evacuate(client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}
