// +build acceptance

package v1

import (
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
	"github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestStackResources(t *testing.T) {
	// Create a provider client for making the HTTP requests.
	// See common.go in this directory for more information.
	client := newClient(t)

	stackName := "postman_stack_2"

	createOpts := stacks.CreateOpts{
		Name:     stackName,
		Template: template,
		Timeout:  5,
	}
	stack, err := stacks.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Created stack: %+v\n", stack)
	defer func() {
		err := stacks.Delete(client, stackName, stack.ID).ExtractErr()
		th.AssertNoErr(t, err)
		t.Logf("Deleted stack (%s)", stackName)
	}()
	err = gophercloud.WaitFor(60, func() (bool, error) {
		getStack, err := stacks.Get(client, stackName, stack.ID).Extract()
		if err != nil {
			return false, err
		}
		if getStack.Status == "CREATE_COMPLETE" {
			return true, nil
		}
		return false, nil
	})

	var resourceName string
	pager := stackresources.List(client, stackName, stack.ID, stackresources.ListOpts{})
	pager.EachPage(func(page pagination.Page) (bool, error) {
		resources, err := stackresources.ExtractResources(page)
		if err != nil {
			return false, err
		}

		t.Logf("%+v\n", resources)

		resourceName = resources[0].Name

		return true, nil
	})

	/*
		resource, err := stackresources.Find(client, stackName).Extract()
		th.AssertNoErr(t, err)
		t.Logf("Got stack resource: %+v\n", resource)
	*/

	resource, err := stackresources.Get(client, stackName, stack.ID, resourceName).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Got stack resource: %+v\n", resource)

	metadata, err := stackresources.Metadata(client, stackName, stack.ID, resourceName).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Got stack resource metadata: %+v\n", metadata)

}
