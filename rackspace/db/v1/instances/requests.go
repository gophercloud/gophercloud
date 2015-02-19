package instances

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/db/v1/backups"
)

// GetDefaultConfig lists the default configuration settings from the template
// that was applied to the specified instance. In a sense, this is the vanilla
// configuration setting applied to an instance. Further configuration can be
// applied by associating an instance with a configuration group.
func GetDefaultConfig(client *gophercloud.ServiceClient, id string) ConfigResult {
	var res ConfigResult

	_, res.Err = client.Request("GET", configURL(client, id), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
	})

	return res
}

// AssociateWithConfigGroup associates a specified instance to a specified
// configuration group. If any of the parameters within a configuration group
// require a restart, then the instance will transition into a restart.
func AssociateWithConfigGroup(client *gophercloud.ServiceClient, instanceID, configGroupID string) UpdateResult {
	reqBody := map[string]string{
		"configuration": configGroupID,
	}

	var res UpdateResult

	_, res.Err = client.Request("PUT", resourceURL(client, instanceID), gophercloud.RequestOpts{
		JSONBody: map[string]map[string]string{"instance": reqBody},
		OkCodes:  []int{202},
	})

	return res
}

// ListBackups will list all the backups for a specified database instance.
func ListBackups(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return backups.BackupPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, backupsURL(client, instanceID), pageFn)
}

// DetachReplica will detach a specified replica instance from its source
// instance, effectively allowing it to operate independently. Detaching a
// replica will restart the MySQL service on the instance.
func DetachReplica(client *gophercloud.ServiceClient, replicaID string) DetachResult {
	var res DetachResult

	_, res.Err = client.Request("PATCH", resourceURL(client, replicaID), gophercloud.RequestOpts{
		JSONBody: map[string]interface{}{"instance": map[string]string{"replica_of": "", "slave_of": ""}},
		OkCodes:  []int{202},
	})

	return res
}
