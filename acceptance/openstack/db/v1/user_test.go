// +build acceptance db

package v1

import (
	"github.com/rackspace/gophercloud/acceptance/tools"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	u "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
)

func (c context) createUsers() {
	users := []string{
		tools.RandomString("user_"),
		tools.RandomString("user_"),
		tools.RandomString("user_"),
	}

	db1 := db.CreateOpt{Name: c.DBIDs[0]}
	db2 := db.CreateOpt{Name: c.DBIDs[1]}
	db3 := db.CreateOpt{Name: c.DBIDs[2]}

	opts := u.BatchCreateOpts{
		u.CreateOpts{
			Name:      users[0],
			Password:  tools.RandomString(),
			databases: db.BatchCreateOpts{db1, db2, db3},
		},
		u.CreateOpts{
			Name:      users[1],
			Password:  tools.RandomString(),
			databases: db.BatchCreateOpts{db1, db2},
		},
		u.CreateOpts{
			Name:      users[2],
			Password:  tools.RandomString(),
			databases: db.BatchCreateOpts{db3},
		},
	}

	err := u.Create(c.client, c.instanceID, opts).ExtractErr()
	c.Logf("Created three users on instance %s: %s, %s, %s", c.instanceID, users[0], users[1], users[2])
	c.users = users
}

func (c context) listUsers() {
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

func (c context) deleteUsers() {
	for _, id := range c.DBIDs {
		err := db.Delete(c.client, c.instanceID, id).ExtractErr()
		c.CheckNoErr(err)
		t.Logf("Deleted DB %s", id)
	}
}
