package testing

import (
	"testing"

	bmvolume "github.com/gophercloud/gophercloud/openstack/baremetal/v1/volume"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListVolumes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleVolumeListSuccessfully(t)
	actual, err := bmvolume.List(client.ServiceClient(), bmvolume.ListOpts{}).Extract()
	th.AssertNoErr(t, err)
	if len(actual.Connectors) <= 0 || len(actual.Targets) <= 0 || len(actual.Links) <= 0 {
		t.Fatalf("Expected Connectors Targets Links, but got not all!")
	}
	th.AssertEquals(t, "http://127.0.0.1:6385/v1/volume/connectors", actual.Connectors[0].(map[string]interface{})["href"].(string))
	th.AssertEquals(t, "http://127.0.0.1:6385/v1/volume/targets", actual.Targets[0].(map[string]interface{})["href"].(string))
	th.AssertEquals(t, "http://127.0.0.1:6385/v1/volume/", actual.Links[0].(map[string]interface{})["href"].(string))
}

func TestListConnectors(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleConnectorListSuccessfully(t)
	pages := 0
	connectorListOpts := bmvolume.ListConnectorsOpts{}
	connectorListOpts.Node = "6d85703a-565d-469a-96ce-30b6de53079d"
	err := bmvolume.ListConnectors(client.ServiceClient(), connectorListOpts).EachPage(func(page pagination.Page) (bool, error) {
		pages++
		actual, err := bmvolume.ExtractConnectors(page)
		if err != nil {
			return false, err
		}
		if len(actual) != 1 {
			t.Fatalf("Expected 1 connectors, but got not!")
		}
		th.AssertEquals(t, "iqn.2017-07.org.openstack:01:d9a51732c3f", actual[0].ConnectorId)
		th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual[0].NodeUUID)
		return true, nil
	})
	th.AssertNoErr(t, err)
	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestCreateConnector(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleConnectorCreationSuccessfully(t, SingleConnectorBody)
	connectorCreateOpts := bmvolume.CreateConnectorOpts{}
	connectorCreateOpts.NodeUUID = "6d85703a-565d-469a-96ce-30b6de53079d"
	connectorCreateOpts.ConnectorType = "iqn"
	connectorCreateOpts.ConnectorId = "iqn.2017-07.org.openstack:01:d9a51732c3f"
	actual, err := bmvolume.CreateConnector(client.ServiceClient(), connectorCreateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual.NodeUUID)
	th.AssertEquals(t, "iqn", actual.ConnectorType)
	th.AssertEquals(t, "iqn.2017-07.org.openstack:01:d9a51732c3f", actual.ConnectorId)
	th.AssertEquals(t, "9bf93e01-d728-47a3-ad4b-5e66a835037c", actual.UUID)
}

func TestDeleteConnector(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleConnectorDeletionSuccessfully(t)
	res := bmvolume.DeleteConnector(client.ServiceClient(), "9bf93e01-d728-47a3-ad4b-5e66a835037c")
	th.AssertNoErr(t, res.Err)
}
func TestGetConnector(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleConnectorGetSuccessfully(t)
	actual, err := bmvolume.GetConnector(client.ServiceClient(), "9bf93e01-d728-47a3-ad4b-5e66a835037c").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get Connector error: %v", err)
	}
	th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual.NodeUUID)
	th.AssertEquals(t, "iqn", actual.ConnectorType)
	th.AssertEquals(t, "iqn.2017-07.org.openstack:01:d9a51732c3f", actual.ConnectorId)
	th.AssertEquals(t, "9bf93e01-d728-47a3-ad4b-5e66a835037c", actual.UUID)
}

func TestUpdateConnector(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleConnectorUpdateSuccessfully(t, SingleConnectorBody)
	actual, err := bmvolume.UpdateConnector(client.ServiceClient(), "9bf93e01-d728-47a3-ad4b-5e66a835037c", bmvolume.UpdateOpts{
		bmvolume.UpdateOperation{
			Op:    bmvolume.ReplaceOp,
			Path:  "/connector_id",
			Value: "iqn.2017-07.org.openstack:01:66666666666",
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update Connector error: %v", err)
	}
	th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual.NodeUUID)
	th.AssertEquals(t, "iqn", actual.ConnectorType)
	th.AssertEquals(t, "iqn.2017-07.org.openstack:01:d9a51732c3f", actual.ConnectorId)
	th.AssertEquals(t, "9bf93e01-d728-47a3-ad4b-5e66a835037c", actual.UUID)
}

func TestListTargets(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTargetListSuccessfully(t)
	pages := 0
	targetListOpts := bmvolume.ListTargetsOpts{}
	targetListOpts.Node = "6d85703a-565d-469a-96ce-30b6de53079d"
	err := bmvolume.ListTargets(client.ServiceClient(), targetListOpts).EachPage(func(page pagination.Page) (bool, error) {
		pages++
		actual, err := bmvolume.ExtractTargets(page)
		if err != nil {
			return false, err
		}
		if len(actual) != 1 {
			t.Fatalf("Expected 1 targets, but got not!")
		}
		th.AssertEquals(t, "bd4d008c-7d31-463d-abf9-6c23d9d55f7f", actual[0].UUID)
		th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual[0].NodeUUID)
		return true, nil
	})
	th.AssertNoErr(t, err)
	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestCreateTarget(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTargetCreationSuccessfully(t, SingleTargetBody)
	targetCreateOpts := bmvolume.CreateTargetOpts{}
	targetCreateOpts.BootIndex = 0
	targetCreateOpts.NodeUUID = "6d85703a-565d-469a-96ce-30b6de53079d"
	targetCreateOpts.VolumeType = "iscsi"
	targetCreateOpts.VolumeId = "04452bed-5367-4202-8bf5-de4335ac56d2"
	actual, err := bmvolume.CreateTarget(client.ServiceClient(), targetCreateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual.NodeUUID)
	th.AssertEquals(t, "iscsi", actual.VolumeType)
	th.AssertEquals(t, "04452bed-5367-4202-8bf5-de4335ac56d2", actual.VolumeId)
	th.AssertEquals(t, "bd4d008c-7d31-463d-abf9-6c23d9d55f7f", actual.UUID)
}

func TestDeleteTarget(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTargetDeletionSuccessfully(t)
	res := bmvolume.DeleteTarget(client.ServiceClient(), "bd4d008c-7d31-463d-abf9-6c23d9d55f7f")
	th.AssertNoErr(t, res.Err)
}
func TestGetTarget(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTargetGetSuccessfully(t)
	actual, err := bmvolume.GetTarget(client.ServiceClient(), "bd4d008c-7d31-463d-abf9-6c23d9d55f7f").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get Target error: %v", err)
	}
	th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual.NodeUUID)
	th.AssertEquals(t, "iscsi", actual.VolumeType)
	th.AssertEquals(t, "04452bed-5367-4202-8bf5-de4335ac56d2", actual.VolumeId)
	th.AssertEquals(t, "bd4d008c-7d31-463d-abf9-6c23d9d55f7f", actual.UUID)
}

func TestUpdateTarget(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTargetUpdateSuccessfully(t, SingleTargetBody)
	actual, err := bmvolume.UpdateTarget(client.ServiceClient(), "bd4d008c-7d31-463d-abf9-6c23d9d55f7f", bmvolume.UpdateOpts{
		bmvolume.UpdateOperation{
			Op:    bmvolume.ReplaceOp,
			Path:  "/volume_id",
			Value: "06666bed-5367-4202-8bf5-de4335ac56d2",
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update Target error: %v", err)
	}
	th.AssertEquals(t, "6d85703a-565d-469a-96ce-30b6de53079d", actual.NodeUUID)
	th.AssertEquals(t, "iscsi", actual.VolumeType)
	th.AssertEquals(t, "04452bed-5367-4202-8bf5-de4335ac56d2", actual.VolumeId)
	th.AssertEquals(t, "bd4d008c-7d31-463d-abf9-6c23d9d55f7f", actual.UUID)
}
