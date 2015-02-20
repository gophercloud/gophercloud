// +build acceptance db

package v1

import (
	"github.com/rackspace/gophercloud/acceptance/tools"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/pagination"
)

func (c context) createDBs() {
	dbs := []string{
		tools.RandomString("db_"),
		tools.RandomString("db_"),
		tools.RandomString("db_"),
	}

	opts := db.BatchCreateOpts{
		db.CreateOpts{Name: dbs[0]},
		db.CreateOpts{Name: dbs[1]},
		db.CreateOpts{Name: dbs[2]},
	}

	res := db.Create(c.client, c.instanceID, opts)
	c.Logf("Created three databases on instance %s: %s, %s, %s", c.instanceID, dbs[0], dbs[1], dbs[2])
	c.DBIDs = dbs
}

func (c context) listDBs() {
	c.Logf("Listing databases on instance %s", c.instanceID)

	err := dbs.List(c.client, c.instanceID).EachPage(func(page pagination.Page) (bool, error) {
		dbList, err := db.ExtractDBs(page)
		c.AssertNoErr(err)

		for _, db := range dbList {
			c.Logf("DB: %#v", db)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) deleteDBs() {
	for _, id := range c.DBIDs {
		err := db.Delete(c.client, c.instanceID, id).ExtractErr()
		c.CheckNoErr(err)
		t.Logf("Deleted DB %s", id)
	}
}
