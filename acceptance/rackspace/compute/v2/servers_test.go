// +build acceptance

package v2

import (
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/acceptance/tools"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	th "github.com/rackspace/gophercloud/testhelper"
)

func createServer(t *testing.T, client *gophercloud.ServiceClient) *os.Server {
	options, err := optionsFromEnv()
	th.AssertNoErr(t, err)

	name := tools.RandomString("Gophercloud-", 8)
	t.Logf("Creating server [%s].", name)
	s, err := servers.Create(client, &os.CreateOpts{
		Name:      name,
		ImageRef:  options.imageID,
		FlavorRef: options.flavorID,
	}).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Server created successfully.")
	return s
}

func deleteServer(t *testing.T, client *gophercloud.ServiceClient, server *os.Server) {
	t.Logf("Deleting server [%s].", server.ID)
	err := servers.Delete(client, server.ID)
	th.AssertNoErr(t, err)
	t.Logf("Server deleted successfully.")
}

func logServer(t *testing.T, server *os.Server, index int) {
	if index == -1 {
		t.Logf("             id=[%s]", server.ID)
	} else {
		t.Logf("[%02d]             id=[%s]", index, server.ID)
	}
	t.Logf("           name=[%s]", server.Name)
	t.Logf("      tenant ID=[%s]", server.TenantID)
	t.Logf("        user ID=[%s]", server.UserID)
	t.Logf("        updated=[%s]", server.Updated)
	t.Logf("        created=[%s]", server.Created)
	t.Logf("        host ID=[%s]", server.HostID)
	t.Logf("    access IPv4=[%s]", server.AccessIPv4)
	t.Logf("    access IPv6=[%s]", server.AccessIPv6)
	t.Logf("          image=[%v]", server.Image)
	t.Logf("         flavor=[%v]", server.Flavor)
	t.Logf("      addresses=[%v]", server.Addresses)
	t.Logf("       metadata=[%v]", server.Metadata)
	t.Logf("          links=[%v]", server.Links)
	t.Logf("        keyname=[%s]", server.KeyName)
	t.Logf(" admin password=[%s]", server.AdminPass)
	t.Logf("         status=[%s]", server.Status)
	t.Logf("       progress=[%d]", server.Progress)
}

func TestCreateServer(t *testing.T) {
	t.Parallel()

	client, err := newClient()
	th.AssertNoErr(t, err)

	s := createServer(t, client)
	defer deleteServer(t, client, s)

	t.Logf("Waiting for server to become active ...")
	err = servers.WaitForStatus(client, s.ID, "ACTIVE", 300)
	th.AssertNoErr(t, err)

	t.Logf("Server launched.")
	logServer(t, s, -1)

	t.Logf("Getting additional server details.")
	details, err := servers.Get(client, s.ID).Extract()
	logServer(t, details, -1)
}

func TestListServers(t *testing.T) {
	t.Parallel()

	client, err := newClient()
	th.AssertNoErr(t, err)

	count := 0
	err = servers.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		t.Logf("-- Page %02d --", count)

		t.Logf("\n%s", page.(os.ServerPage).PrettyPrintJSON())

		s, err := servers.ExtractServers(page)
		th.AssertNoErr(t, err)
		for index, server := range s {
			logServer(t, &server, index)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}
