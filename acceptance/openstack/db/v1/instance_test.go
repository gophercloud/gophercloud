// +build acceptance db

package v1

import (
	"testing"

	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/gophercloud/pagination"
	rackspaceInst "github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestRunner(t *testing.T) {
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
	c.createDBs()
	c.listDBs()

	// USER tests
	c.createUsers()
	c.listUsers()

	// TEARDOWN
	c.deleteUsers()
	c.deleteDBs()
	c.deleteInstance()
}

func (c context) createInstance() {
	opts := rackspaceInst.CreateOpts{
		FlavorRef: "1",
		Size:      1,
		Name:      tools.RandomString("gopher_db", 5),
		Datastore: &rackspaceInst.DatastoreOpts{Version: "5.6", Type: "MySQL"},
	}

	instance, err := instances.Create(c.client, opts).Extract()
	th.AssertNoErr(c.test, err)

	c.Logf("Restarting %s. Waiting...", instance.ID)
	c.WaitUntilActive(instance.ID)
	c.Logf("Created DB %#v", instance)

	c.instanceID = instance.ID
}

func (c context) listInstances() {
	c.Logf("Listing instances")

	err := instances.List(c.client).EachPage(func(page pagination.Page) (bool, error) {
		instanceList, err := instances.ExtractInstances(page)
		c.AssertNoErr(err)

		for _, i := range instanceList {
			c.Logf("Instance: %#v", i)
		}

		return true, nil
	})

	c.AssertNoErr(err)
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
	_, err := instances.EnableRootUser(c.client, c.instanceID).Extract()
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
	err := instances.RestartService(c.client, id).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Restarting %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Restarted %s", id)
}

func (c context) resizeInstance() {
	id := c.instanceID
	err := instances.ResizeInstance(c.client, id, "2").ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Resizing %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Resized %s with flavorRef %s", id, "2")
}

func (c context) resizeVol() {
	id := c.instanceID
	err := instances.ResizeVolume(c.client, id, 2).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Resizing volume of %s. Waiting...", id)
	c.WaitUntilActive(id)
	c.Logf("Resized the volume of %s to %d GB", id, 2)
}
