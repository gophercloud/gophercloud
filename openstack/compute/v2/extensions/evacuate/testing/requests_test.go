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

	_, err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{
		Host:            "derp",
		AdminPass:       "MySecretPass",
		OnSharedStorage: false,
	}).ExtractAdminPass()
	if err != nil {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}

func TestEvacuateWithHost(t *testing.T) {
	const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockEvacuateResponseWithHost(t, serverID)

	_, err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{
		Host: "derp",
	}).ExtractAdminPass()
	if err != nil {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}

func TestEvacuateWithNoOpts(t *testing.T) {
	const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockEvacuateResponseWithNoOpts(t, serverID)

	_, err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{}).ExtractAdminPass()
	if err != nil {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}

func TestEvacuateAdminpassResponse(t *testing.T) {
	const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockEvacuateAdminpassResponse(t, serverID)

	actual, err := evacuate.Evacuate(client.ServiceClient(), serverID, evacuate.EvacuateOpts{}).ExtractAdminPass()
	th.CheckEquals(t, "MySecretPass", actual)
	if err != nil {
		t.Fatalf("Unable to evacuate to server: %s", err)
	}
}
