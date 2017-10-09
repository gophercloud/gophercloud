package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/evacuate"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestEvacuate(t *testing.T) {
	const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockEvacuateResponse(t, serverID)

	err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{
		Host:            "derp",
		AdminPass:       "true",
		OnSharedStorage: false,
	}).ExtractErr()
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}

func TestEvacuateWithHost(t *testing.T) {
	const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockEvacuateResponseWithHost(t, serverID)

	err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{
		Host: "derp",
	}).ExtractErr()
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}

func TestEvacuateWithNoOpts(t *testing.T) {
	const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockEvacuateResponseWithNoOpts(t, serverID)

	err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{}).ExtractErr()
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}
