// +build acceptance db

package v1

import (
	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestRunner(t *testingT) {
	c := newContext(t)

	// FLAVOR tests
	c.listFlavors()
	c.getFlavor()

	// INSTANCE tests
	c.createInstance()
	c.listInstances()
	c.getInstance()
	c.isRootEnabled()
	c.enableRootUser()
	c.isRootEnabled()
	c.restartInstance()
	c.resizeInstance()
	c.resizeVol()

	// DATABASE tests
	c.createDB()
	c.listDBs()

	// USER tests
	c.createUsers()
	c.listUsers()

	// TEARDOWN
	c.deleteUsers()
	c.deleteDBs()
	c.deleteInstance(id)
}

func (c context) createInstance() {
	opts := instances.CreateOpts{
		FlavorRef: "1",
		Size:      1,
		Name:      tools.RandomString("gopher_db", 5),
	}

	instance, err := instances.Create(c.client, opts).Extract()
	th.AssertNoErr(c.test, err)

	c.Logf("Restarting %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Created DB %#v", instance)

	c.instanceID = instance.ID
}

func (c context) listInstances() {
	c.Logf("Listing instances")

	err := instances.List(c.client).EachPage(func(page pagination.Page) (bool, error) {
		instanceList, err := instances.ExtractInstances(page)
		c.AssertNoErr(err)

		for _, n := range networkList {
			c.Logf("Instance: %#v", instance)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) getInstance() {
	instance, err := instances.Get(c.client, c.instanceID).Extract()
	c.AssertNoErr(err)
	c.Logf("Getting instance: %#v", instance)
}

func (c context) deleteInstance() {
	err := instances.Delete(c.client, c.instanceID).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Deleted instance %s", c.instanceID)
}

func (c context) enableRootUser() {
	err := instances.EnableRootUser(c.client, c.instanceID).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Enabled root user on %s", c.instanceID)
}

func (c context) isRootEnabled() {
	enabled, err := instances.IsRootEnabled(c.client, c.instanceID)
	c.AssertNoErr(err)
	c.Logf("Is root enabled? %s", enabled)
}

func (c context) restartInstance() {
	id := c.instanceID
	err := instances.Restart(c.client, id).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Restarting %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Restarted %s", id)
}

func (c context) resizeInstance() {
	id := c.instanceID
	err := instances.Resize(c.client, id, "2").ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Resizing %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Resized %s with flavorRef %s", id, "2")
}

func (c context) resizeVol() {
	id := c.instanceID
	err := instances.ResizeVol(c.client, id, 2).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Resizing volume of %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Resized the volume of %s to %d GB", id, 2)
}
