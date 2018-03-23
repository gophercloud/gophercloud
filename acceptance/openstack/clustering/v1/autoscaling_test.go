package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"

	// TODO: Dependent on nodes:Get() Uncomment this if using proper validation
	//th "github.com/gophercloud/gophercloud/testhelper"

	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
)

var testName string

func TestAutoScaling(t *testing.T) {
	testName = tools.RandomString("TESTACC-", 8)
	testName = "node-Fbxesyrn"
	nodeUpdate(t)
}

func nodeUpdate(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	if err != nil {
		t.Fatalf("Unable to create clustering client: %v", err)
	}

	nodeName := testName
	newNodeName := nodeName + "-UPDATE_TEST"

	// Update to new node name
	actionID, err := nodes.Update(client, nodeName, nodes.UpdateOpts{Name: newNodeName}).ExtractAction()
	if err != nil {
		t.Fatalf("Unable to update node: %v", err)
	}

	WaitForNodeToUpdate(client, actionID, 15)

	// TODO: Dependent on nodes:Get()
	/*
		node, err := nodes.Get(client, actionID).Extract()
		if err != nil {
			t.Fatalf("Unable to get node: %v", err)
		}
		th.AssertEquals(t, newNodeName, node.Name)
	*/

	// Revert back to original node name
	actionID, err = nodes.Update(client, newNodeName, nodes.UpdateOpts{Name: nodeName}).ExtractAction()
	if err != nil {
		t.Fatalf("Unable to update node: %v", err)
	}

	WaitForNodeToUpdate(client, actionID, 15)

	// TODO: Dependent on nodes:Get()
	/*
		node, err := nodes.Get(client, actionID).Extract()
		if err != nil {
			t.Fatalf("Unable to get node: %v", err)
		}
		th.AssertEquals(t, nodeName, node.Name)
	*/
}

func WaitForNodeToUpdate(client *gophercloud.ServiceClient, actionID string, secs int) error {

	time.Sleep(time.Duration(secs) * time.Second)
	return nil

	// TODO: Dependent on actions:Get() Proper way of handling update but requires actions
	/*
		return gophercloud.WaitFor(secs, func() (bool, error) {
			action, err := actions.Get(client, actionID).Extract()
			if err != nil {
				return false, err
			}
			switch actions.Status {
			case "SUCCEEDED":
				return true, nil
			case "RUNNING":
				return false, nil
			default:
				return false, fmt.Errorf("Error WaitFor. Received status=%v", actions.Status)
			}
		})
	*/
}
