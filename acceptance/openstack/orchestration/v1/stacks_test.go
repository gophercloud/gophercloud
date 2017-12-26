// +build acceptance

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestStacks(t *testing.T) {
	// Create a provider client for making the HTTP requests.
	// See common.go in this directory for more information.
	client := newClient(t)

	stackName1 := "gophercloud-test-stack-2"
	var templateOpts = new(stacks.Template)
	templateOpts.Bin = []byte(template)

	createOpts := stacks.CreateOpts{
		Name:         stackName1,
		TemplateOpts: templateOpts,
		Timeout:      5,
	}
	stack, err := stacks.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Created stack: %+v\n", stack)
	defer func() {
		err := stacks.Delete(client, stackName1, stack.ID).ExtractErr()
		th.AssertNoErr(t, err)
		t.Logf("Deleted stack (%s)", stackName1)
	}()
	err = gophercloud.WaitFor(60, func() (bool, error) {
		getStack, err := stacks.Get(client, stackName1, stack.ID).Extract()
		if err != nil {
			return false, err
		}
		if getStack.Status == "CREATE_COMPLETE" {
			return true, nil
		}
		return false, nil
	})

	updateOpts := stacks.UpdateOpts{
		TemplateOpts: templateOpts,
		Timeout:      20,
	}
	err = stacks.Update(client, stackName1, stack.ID, updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
	err = gophercloud.WaitFor(60, func() (bool, error) {
		getStack, err := stacks.Get(client, stackName1, stack.ID).Extract()
		if err != nil {
			return false, err
		}
		if getStack.Status == "UPDATE_COMPLETE" {
			return true, nil
		}
		return false, nil
	})

	t.Logf("Updated stack")

	err = stacks.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		stackList, err := stacks.ExtractStacks(page)
		th.AssertNoErr(t, err)

		t.Logf("Got stack list: %+v\n", stackList)

		return true, nil
	})
	th.AssertNoErr(t, err)

	getStack, err := stacks.Get(client, stackName1, stack.ID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Got stack: %+v\n", getStack)

	abandonedStack, err := stacks.Abandon(client, stackName1, stack.ID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Abandonded stack %+v\n", abandonedStack)
	th.AssertNoErr(t, err)
}

// Test using the updated interface
func TestStacksNewTemplateFormat(t *testing.T) {
	// Create a provider client for making the HTTP requests.
	// See common.go in this directory for more information.
	client := newClient(t)

	stackName1 := "gophercloud-test-stack-2"
	templateOpts := new(stacks.Template)
	templateOpts.Bin = []byte(template)
	createOpts := stacks.CreateOpts{
		Name:         stackName1,
		TemplateOpts: templateOpts,
		Timeout:      5,
	}
	stack, err := stacks.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Created stack: %+v\n", stack)
	defer func() {
		err := stacks.Delete(client, stackName1, stack.ID).ExtractErr()
		th.AssertNoErr(t, err)
		t.Logf("Deleted stack (%s)", stackName1)
	}()
	err = gophercloud.WaitFor(60, func() (bool, error) {
		getStack, err := stacks.Get(client, stackName1, stack.ID).Extract()
		if err != nil {
			return false, err
		}
		if getStack.Status == "CREATE_COMPLETE" {
			return true, nil
		}
		return false, nil
	})

	updateOpts := stacks.UpdateOpts{
		TemplateOpts: templateOpts,
		Timeout:      20,
	}
	err = stacks.Update(client, stackName1, stack.ID, updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
	err = gophercloud.WaitFor(60, func() (bool, error) {
		getStack, err := stacks.Get(client, stackName1, stack.ID).Extract()
		if err != nil {
			return false, err
		}
		if getStack.Status == "UPDATE_COMPLETE" {
			return true, nil
		}
		return false, nil
	})

	t.Logf("Updated stack")

	err = stacks.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		stackList, err := stacks.ExtractStacks(page)
		th.AssertNoErr(t, err)

		t.Logf("Got stack list: %+v\n", stackList)

		return true, nil
	})
	th.AssertNoErr(t, err)

	getStack, err := stacks.Get(client, stackName1, stack.ID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Got stack: %+v\n", getStack)

	abandonedStack, err := stacks.Abandon(client, stackName1, stack.ID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Abandonded stack %+v\n", abandonedStack)
	th.AssertNoErr(t, err)
}
