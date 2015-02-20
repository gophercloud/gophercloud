// +build acceptance db

package v1

import "github.com/rackspace/gophercloud/pagination"

func (c context) createBackup() {
	opts := backups.CreateOpts{
		Name:       tools.PrefixString("backup_", 5),
		InstanceID: c.instanceID,
	}

	backup, err := backups.Create(c.client, opts)

	c.Logf("Created backup %#v", backup)
	c.AssertNoErr(t, err)

	c.backupID = backup.ID
}

func (c context) getBackup() {
	backup, err := backups.Get(c.client, c.backupID).Extract()
	c.AssertNoErr(err)
	c.Logf("Getting backup %s", backup.ID)
}

func (c context) listAllBackups() {
	c.Logf("Listing backups")

	err := backups.List(c.client).EachPage(func(page pagination.Page) (bool, error) {
		backupList, err := backups.ExtractBackups(page)
		c.AssertNoErr(err)

		for _, b := range backupList {
			c.Logf("Backup: %#v", b)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) listInstanceBackups() {
	c.Logf("Listing backups for instance %s", c.instanceID)

	err := instances.ListBackups(c.client).EachPage(func(page pagination.Page) (bool, error) {
		backupList, err := backups.ExtractBackups(page)
		c.AssertNoErr(err)

		for _, b := range backupList {
			c.Logf("Backup: %#v", b)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) deleteBackup() {
	err := backups.Delete(c.client, c.backupID).ExtractErr()
	c.AssertNoErr(err)
	c.Logf("Deleted backup %s", c.backupID)
}
