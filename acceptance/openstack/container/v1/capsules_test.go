package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/container/v1/capsules"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCapsuleGet(t *testing.T) {
	client, err := clients.NewContainerV1Client()
	if err != nil {
		t.Fatalf("Unable to create an container v1 client: %v", err)
	}
	th.AssertNoErr(t, err)
	capsuleUUID := "e6c913bb-b4e4-409d-8b71-3e029f196458"
	if capsuleUUID == "" {
		t.Fatalf("In order to retrieve a capsule, the CapsuleUUID must be set")
	}
	capsule, err := capsules.Get(client, capsuleUUID).Extract()
	// Get a capsule

	th.AssertNoErr(t, err)
	th.AssertEquals(t, capsule.Status, "Running")
	th.AssertEquals(t, capsule.MetaName, "template")
	th.AssertEquals(t, capsule.CPU, float64(2.0))
}
