// +build acceptance compute rescueunrescue

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServerRescue(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	err = RescueServer(t, client, server)
	th.AssertNoErr(t, err)
}
