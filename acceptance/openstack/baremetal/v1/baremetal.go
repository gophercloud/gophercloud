package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

// CreateNode creates a basic node with a randomly generated name.
func CreateNode(t *testing.T, client *gophercloud.ServiceClient) (*nodes.Node, error) {
	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create bare metal node: %s", name)

	node, err := nodes.Create(client, nodes.CreateOpts{
		Name:          name,
		Driver:        "ipmi",
		BootInterface: "pxe",
		DriverInfo: map[string]interface{}{
			"ipmi_port":      "6230",
			"ipmi_username":  "admin",
			"deploy_kernel":  "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
			"ipmi_address":   "192.168.122.1",
			"deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
			"ipmi_password":  "admin",
		},
	}).Extract()

	return node, err
}

// DeleteNode deletes a bare metal node via its UUID.
func DeleteNode(t *testing.T, client *gophercloud.ServiceClient, node *nodes.Node) {
	err := nodes.Delete(client, node.UUID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete node %s: %s", node.UUID, err)
	}

	t.Logf("Deleted server: %s", node.UUID)
}

// CreatePort - creates a port with a fixed NodeUUID and Address
func CreatePort(t *testing.T, client *gophercloud.ServiceClient) (*ports.Port, error){
	nodeUUID := "1be26c0b-03f2-4d2e-ae87-c02d7f33c123"
	address := "fe:54:00:77:07:d9"
	t.Logf("Attempting to create Port for Node: %s with Address: %s", nodeUUID, address)

	port, err := ports.Create(client, ports.CreateOpts{
		UUID: "27e3153e-d5bf-4b7e-b517-fb518e17f34c",
    NodeUUID: nodeUUID,
    Address: address,
    PXEEnabled: true,
	}).Extract()

	return port, err
}

// DeletePort - deletes a port via its UUID
func DeletePort(t *testing.T, client *gophercloud.ServiceClient, port *ports.Port){
	err := ports.Delete(client, port.UUID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete port %s: %s", port.UUID, err)
	}

	t.Logf("Deleted port: %s", port.UUID)

}
