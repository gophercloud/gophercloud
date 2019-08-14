package tags

import "github.com/gophercloud/gophercloud"

// List all tags on a server.
func List(client *gophercloud.ServiceClient, serverID string) (r ListResult) {
	url := listURL(client, serverID)
	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
