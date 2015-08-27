package roles

import "github.com/rackspace/gophercloud"

func roleAssignmentsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("role_assignments")
}
