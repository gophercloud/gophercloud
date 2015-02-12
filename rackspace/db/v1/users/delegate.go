package users

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
)

func Create(client *gophercloud.ServiceClient, instanceID string, opts os.CreateOptsBuilder) os.CreateResult {
	return os.Create(client, instanceID, opts)
}

func List(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	return os.List(client, instanceID)
}

func Delete(client *gophercloud.ServiceClient, instanceID, userName string) os.DeleteResult {
	return os.Delete(client, instanceID, userName)
}
