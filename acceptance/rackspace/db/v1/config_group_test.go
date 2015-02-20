// +build acceptance db

package v1

import (
	config "github.com/rackspace/gophercloud/openstack/db/v1/configurations"
	"github.com/rackspace/gophercloud/pagination"
)

func (c context) createConfigGrp() {
	opts := config.CreateOpts{
		Name: tools.PrefixString("config_", 5),
		Values: map[string]interface{}{
			"connect_timeout":  300,
			"join_buffer_size": 900000,
		},
	}

	cg, err := config.Create(c.client, opts)

	c.AssertNoErr(err)
	c.Logf("Created config group %#v", cg)

	c.configGroupID = cg.ID
}

func (c context) getConfigGrp() {
	cg, err := config.Get(c.client, c.configGroupID)
	c.Logf("Getting config group: %#v", cg)
	c.AssertNoErr(err)
}

func (c context) updateConfigGrp() {
	opts := config.UpdateOpts{
		Name: tools.PrefixString("new_name_", 5),
		Values: map[string]interface{}{
			"connect_timeout": 250,
		},
	}
	err := config.Update(c.client, c.configGroupID, opts).ExtractErr()
	c.Logf("Updated config group %s", c.configGroupID)
	c.AssertNoErr(err)
}

func (c context) replaceConfigGrp() {
	opts := config.UpdateOpts{
		Values: map[string]interface{}{
			"expire_logs_days": 7,
		},
	}

	err := config.Replace(c.client, c.configGroupID, opts).ExtractErr()
	c.Logf("Replaced values for config group %s", c.configGroupID)
	c.AssertNoErr(err)
}

func (c context) associateInstanceWithConfigGrp() {
	err := config.AssociateWithConfigGroup(c.client, c.instanceID, c.configGroupID).ExtractErr()
	c.Logf("Associated instance %s with config group %s", c.instanceID, c.configGroupID)
	c.AssertNoErr(err)
}

func (c context) listConfigGrpInstances() {
	c.Logf("Listing all instances associated with config group %s", c.configGroupID)

	err := config.ListInstances(c.client, c.configGroupID).EachPage(func(page pagination.Page) (bool, error) {
		instanceList, err := instances.ExtractInstances(page)
		c.AssertNoErr(err)

		for _, n := range networkList {
			c.Logf("Instance: %#v", instance)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) deleteConfigGrp() {
	err := config.Delete(c.client, c.configGroupID).ExtractErr()
	c.Logf("Deleted config group %s", c.configGroupID)
	c.AssertNoErr(err)
}
