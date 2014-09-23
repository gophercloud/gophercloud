// +build acceptance

package v2

import (
	"fmt"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

func TestListServers(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	t.Logf("ID\tRegion\tName\tStatus\tIPv4\tIPv6")

	pager := servers.List(client)
	count, pages := 0, 0
	pager.EachPage(func(page pagination.Page) (bool, error) {
		pages++
		t.Logf("---")

		servers, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		for _, s := range servers {
			t.Logf("%s\t%s\t%s\t%s\t%s\t\n", s.ID, s.Name, s.Status, s.AccessIPv4, s.AccessIPv6)
			count++
		}

		return true, nil
	})

	fmt.Printf("--------\n%d servers listed on %d pages.\n", count, pages)
}

func createServer(t *testing.T, client *gophercloud.ServiceClient, choices *ComputeChoices) (*servers.Server, error) {
	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create server: %s\n", name)

	server, err := servers.Create(client, map[string]interface{}{
		"flavorRef": choices.FlavorID,
		"imageRef":  choices.ImageID,
		"name":      name,
	})
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	return servers.ExtractServer(server)
}

func TestCreateDestroyServer(t *testing.T) {
	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create server: %s\n", name)

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	defer func() {
		servers.Delete(client, server.ID)
		t.Logf("Server deleted.")
	}()

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
}

func TestUpdateServer(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}
	defer servers.Delete(client, server.ID)

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	alternateName := tools.RandomString("ACPTTEST", 16)
	for alternateName == server.Name {
		alternateName = tools.RandomString("ACPTTEST", 16)
	}

	t.Logf("Attempting to rename the server to %s.", alternateName)

	result, err := servers.Update(client, server.ID, map[string]interface{}{
		"name": alternateName,
	})
	if err != nil {
		t.Fatalf("Unable to rename server: %v", err)
	}
	updated, err := servers.ExtractServer(result)
	if err != nil {
		t.Fatalf("Unable to extract server: %v", err)
	}

	if updated.ID != server.ID {
		t.Errorf("Updated server ID [%s] didn't match original server ID [%s]!", updated.ID, server.ID)
	}

	err = tools.WaitFor(func() (bool, error) {
		result, err := servers.Get(client, updated.ID)
		if err != nil {
			return false, err
		}
		latest, err := servers.ExtractServer(result)
		if err != nil {
			return false, err
		}

		return latest.Name == alternateName, nil
	})
}

func TestActionChangeAdminPassword(t *testing.T) {
	t.Parallel()

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}
	defer servers.Delete(client, server.ID)

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	randomPassword := tools.MakeNewPassword(server.AdminPass)
	err = servers.ChangeAdminPassword(client, server.ID, randomPassword)
	if err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "PASSWORD"); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestActionReboot(t *testing.T) {
	t.Parallel()

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}
	defer servers.Delete(client, server.ID)

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	err = servers.Reboot(client, server.ID, "aldhjflaskhjf")
	if err == nil {
		t.Fatal("Expected the SDK to provide an ArgumentError here")
	}

	t.Logf("Attempting reboot of server %s", server.ID)
	err = servers.Reboot(client, server.ID, servers.OSReboot)
	if err != nil {
		t.Fatalf("Unable to reboot server: %v", err)
	}

	if err = waitForStatus(client, server, "REBOOT"); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestActionRebuild(t *testing.T) {
	t.Parallel()

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}
	defer servers.Delete(client, server.ID)

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	t.Logf("Attempting to rebuild server %s", server.ID)

	newPassword := tools.MakeNewPassword(server.AdminPass)
	newName := tools.RandomString("ACPTTEST", 16)
	result, err := servers.Rebuild(client, server.ID, newName, newPassword, choices.ImageID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rebuilt, err := servers.ExtractServer(result)
	if err != nil {
		t.Fatal(err)
	}
	if rebuilt.ID != server.ID {
		t.Errorf("Expected rebuilt server ID of [%s]; got [%s]", server.ID, rebuilt.ID)
	}

	if err = waitForStatus(client, rebuilt, "REBUILD"); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, rebuilt, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func resizeServer(t *testing.T, client *gophercloud.ServiceClient, server *servers.Server, choices *ComputeChoices) {
	if err := waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	t.Logf("Attempting to resize server [%s]", server.ID)

	if err := servers.Resize(client, server.ID, choices.FlavorIDResize); err != nil {
		t.Fatal(err)
	}

	if err := waitForStatus(client, server, "VERIFY_RESIZE"); err != nil {
		t.Fatal(err)
	}
}

func TestActionResizeConfirm(t *testing.T) {
	t.Parallel()

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}
	defer servers.Delete(client, server.ID)
	resizeServer(t, client, server, choices)

	t.Logf("Attempting to confirm resize for server %s", server.ID)

	if err = servers.ConfirmResize(client, server.ID); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestActionResizeRevert(t *testing.T) {
	t.Parallel()

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}
	defer servers.Delete(client, server.ID)
	resizeServer(t, client, server, choices)

	t.Logf("Attempting to revert resize for server %s", server.ID)

	if err := servers.RevertResize(client, server.ID); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}
