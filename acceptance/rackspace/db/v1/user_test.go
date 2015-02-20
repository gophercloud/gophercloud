// +build acceptance db

package v1

import (
	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/identity/v2/users"
	"github.com/rackspace/gophercloud/pagination"
	db "github.com/rackspace/gophercloud/rackspace/db/v1/databases"
	u "github.com/rackspace/gophercloud/rackspace/db/v1/users"
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
	c.Logf("Listing users on instance %s", c.instanceID)

	err := users.List(c.client, c.instanceID).EachPage(func(page pagination.Page) (bool, error) {
		uList, err := u.ExtractUsers(page)
		c.AssertNoErr(err)

		for _, u := range uList {
			c.Logf("User: %#v", u)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) deleteUsers() {
	for _, id := range c.users {
		err := u.Delete(c.client, c.instanceID, id).ExtractErr()
		c.CheckNoErr(err)
		t.Logf("Deleted user %s", id)
	}
}

func (c context) changeUserPwd() {
	opts := u.BatchCreateOpts{}

	for _, id := range c.users[:1] {
		opts = append(opts, u.CreateOpts{Password: tools.PrefixString("", 5)})
	}

	err := u.UpdatePassword(c.client, c.instanceID, opts).ExtractErr()
	c.Logf("Updated 2 users' passwords")
	c.AssertNoErr(err)
}

func (c context) getUser() {
	user, err := u.Get(c.client, c.instanceID, c.users[0]).Extract()
	c.Logf("Getting user %s", user)
	c.AssertNoErr(err)
}

func (c context) updateUser() {
	opts := u.CreateOpts{Name: tools.PrefixString("new_name_", 5)}
	user, err := u.Update(c.client, c.instanceID, c.users[0], opts).Extract()
	c.Logf("Updated user %s", user)
	c.AssertNoErr(err)
}

func (c context) listUserAccess() {
	err := u.ListAccess(c.client, c.instanceId, c.users[0]).EachPage(func(page pagination.Page) (bool, error) {
		dbList, err := db.ExtractDBs(page)
		c.AssertNoErr(err)

		for _, db := range dbList {
			c.Logf("User %s has access to DB: %#v", db)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) grantUserAccess() {
	userID, dbID := c.users[0], c.DBIDS[0]
	err := u.GrantUserAccess(c.client, c.instanceID, userID, dbID)
	c.Logf("Granted access for user %s to DB %s", userID, dbID)
}

func (c context) revokeUserAccess() {
	userID, dbID := c.users[0], c.DBIDS[0]
	err := u.RevokeUserAccess(c.client, c.instanceID, userID, dbID)
	c.Logf("Revoked access for user %s to DB %s", userID, dbID)
}
