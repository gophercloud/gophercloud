package instances

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/db/v1/backups"
)

func GetDefaultConfig(client *gophercloud.ServiceClient, id string) ConfigResult {
	var res ConfigResult

	_, res.Err = perigee.Request("GET", configURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

func AssociateWithConfigGroup(client *gophercloud.ServiceClient, instanceID, configGroupID string) UpdateResult {
	reqBody := map[string]string{
		"configuration": configGroupID,
	}

	var res UpdateResult

	_, res.Err = perigee.Request("PUT", resourceURL(client, instanceID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     map[string]map[string]string{"instance": reqBody},
		OkCodes:     []int{202},
	})

	return res
}

func ListBackups(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return backups.BackupPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, backupsURL(client, instanceID), pageFn)
}

func DetachReplica(client *gophercloud.ServiceClient, replicaID string) DetachResult {
	var res DetachResult

	_, res.Err = perigee.Request("PATCH", resourceURL(client, replicaID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     map[string]interface{}{"instance": map[string]string{"replica_of": "", "slave_of": ""}},
		OkCodes:     []int{202},
	})

	return res
}
