// +build acceptance rackspace networking v2

package v2

import (
  "testing"

  osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
  "github.com/rackspace/gophercloud/pagination"
  "github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
  "github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
  "github.com/rackspace/gophercloud/rackspace/networking/v2/virtualinterfaces"
  th "github.com/rackspace/gophercloud/testhelper"
)

func TestVirtualInterfaces(t *testing.T) {
  Setup(t)
  defer Teardown()

  // Get a Server
  var serverID string
  pager := servers.List(Client, osServers.ListOpts{Limit:1})
  err := pager.EachPage(func(page pagination.Page) (bool, error) {
    servers, err := servers.ExtractServers(page)
    if err != nil {
      return false, err
    }
    serverID = servers[0].ID
    return true, nil
  })
  th.AssertNoErr(t, err)

  t.Logf("ServerID: %s", serverID)

  // Create a network
  n, err := networks.Create(Client, networks.CreateOpts{Label: "sample_network", CIDR: "172.20.0.0/24"}).Extract()
  th.AssertNoErr(t, err)
  defer networks.Delete(Client, n.ID)
  networkID := n.ID

  t.Logf("NetworkID: %s", networkID)

  // Create a virtual interface
  vi, err := virtualinterfaces.Create(Client, serverID, networkID).Extract()
  th.AssertNoErr(t, err)
  t.Logf("Created virtual interface: %+v\n", vi)
  defer virtualinterfaces.Delete(Client, serverID, vi.ID)

  // List virtual interfaces
  pager = virtualinterfaces.List(Client, serverID)
  err = pager.EachPage(func(page pagination.Page) (bool, error) {
    t.Logf("--- Page ---")

    virtualinterfacesList, err := virtualinterfaces.ExtractVirtualInterfaces(page)
    th.AssertNoErr(t, err)

    for _, vi := range virtualinterfacesList {
      t.Logf("Virtual Interface: ID [%s] MAC Address [%s] IP Addresses [%v]",
        vi.ID, vi.MACAddress, vi.IPAddresses)
    }

    return true, nil
  })
  th.CheckNoErr(t, err)
}
