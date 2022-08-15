package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/allocations"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	bmvolume "github.com/gophercloud/gophercloud/openstack/baremetal/v1/volume"
)

// CreateNode creates a basic node with a randomly generated name.
func CreateNode(t *testing.T, client *gophercloud.ServiceClient) (*nodes.Node, error) {
	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create bare metal node: %s", name)

	node, err := nodes.Create(client, nodes.CreateOpts{
		Name:          name,
		Driver:        "ipmi",
		BootInterface: "ipxe",
		RAIDInterface: "agent",
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

// CreateAllocation creates an allocation
func CreateAllocation(t *testing.T, client *gophercloud.ServiceClient) (*allocations.Allocation, error) {
	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create bare metal allocation: %s", name)

	allocation, err := allocations.Create(client, allocations.CreateOpts{
		Name:          name,
		ResourceClass: "baremetal",
	}).Extract()

	return allocation, err
}

// DeleteAllocation deletes a bare metal allocation via its UUID.
func DeleteAllocation(t *testing.T, client *gophercloud.ServiceClient, allocation *allocations.Allocation) {
	err := allocations.Delete(client, allocation.UUID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete allocation %s: %s", allocation.UUID, err)
	}

	t.Logf("Deleted allocation: %s", allocation.UUID)
}

// CreateFakeNode creates a node with fake-hardware to use for port tests.
func CreateFakeNode(t *testing.T, client *gophercloud.ServiceClient) (*nodes.Node, error) {
	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create bare metal node: %s", name)

	node, err := nodes.Create(client, nodes.CreateOpts{
		Name:          name,
		Driver:        "fake-hardware",
		BootInterface: "fake",
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

// CreatePort - creates a port for a node with a fixed Address
func CreatePort(t *testing.T, client *gophercloud.ServiceClient, node *nodes.Node) (*ports.Port, error) {
	mac := "e6:72:1f:52:00:f4"
	t.Logf("Attempting to create Port for Node: %s with Address: %s", node.UUID, mac)

	iTrue := true
	port, err := ports.Create(client, ports.CreateOpts{
		NodeUUID:   node.UUID,
		Address:    mac,
		PXEEnabled: &iTrue,
	}).Extract()

	return port, err
}

// DeletePort - deletes a port via its UUID
func DeletePort(t *testing.T, client *gophercloud.ServiceClient, port *ports.Port) {
	err := ports.Delete(client, port.UUID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete port %s: %s", port.UUID, err)
	}

	t.Logf("Deleted port: %s", port.UUID)

}

func UpdateNodeStorageInterface(client *gophercloud.ServiceClient, nodeId string, storage string) (*nodes.Node, error) {
	updated, err := nodes.Update(client, nodeId, nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/storage_interface",
			Value: storage,
		},
	}).Extract()
	return updated, err
}

func SetNodePowerOff(client *gophercloud.ServiceClient, nodeId string) error {
	opts := nodes.PowerStateOpts{
		Target:  nodes.PowerOff,
		Timeout: 100,
	}
	err := nodes.ChangePowerState(client, nodeId, opts).ExtractErr()
	return err
}

func CreateVolumeConnector(t *testing.T, client *gophercloud.ServiceClient, node *nodes.Node) (*bmvolume.Connector, error) {
	connectorCreateOpts := bmvolume.CreateConnectorOpts{}
	connectorCreateOpts.NodeUUID = node.UUID
	connectorCreateOpts.ConnectorType = "iqn"
	connectorCreateOpts.ConnectorId = "iqn.2017-07.org.openstack." + node.UUID
	t.Logf("Attempting to create volume connector for Node: %s with connector: %s", node.UUID, connectorCreateOpts.ConnectorId)
	connector, err := bmvolume.CreateConnector(client, connectorCreateOpts).Extract()
	return connector, err
}

func DeleteVolumeConnector(t *testing.T, client *gophercloud.ServiceClient, connector *bmvolume.Connector) {
	err := bmvolume.DeleteConnector(client, connector.UUID).Err
	if err != nil {
		t.Fatalf("Unable to delete volume connector %s", connector.UUID)
	}
	t.Logf("Deleted volume connector: %s", connector.UUID)
}

func CreateVolumeTarget(t *testing.T, client *gophercloud.ServiceClient, node *nodes.Node, volumeId string) (*bmvolume.Target, error) {
	targetCreateOpts := bmvolume.CreateTargetOpts{}
	targetCreateOpts.NodeUUID = node.UUID
	targetCreateOpts.BootIndex = "0"
	targetCreateOpts.VolumeType = "iscsi"
	targetCreateOpts.VolumeId = volumeId
	t.Logf("Attempting to create volume target for Node: %s with volumeId: %s", node.UUID, volumeId)
	target, err := bmvolume.CreateTarget(client, targetCreateOpts).Extract()
	return target, err
}

func DeleteVolumeTarget(t *testing.T, client *gophercloud.ServiceClient, target *bmvolume.Target) {
	err := bmvolume.DeleteTarget(client, target.UUID).Err
	if err != nil {
		t.Fatalf("Unable to delete volume target %s", target.UUID)
	}
	t.Logf("Deleted volume target: %s", target.UUID)
}
