package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestQoS(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	qs, err := CreateQoS(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQoS(t, client, qs)
}
