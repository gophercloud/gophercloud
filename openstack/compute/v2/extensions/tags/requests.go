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

// Check if a tag exists on a server.
func Check(client *gophercloud.ServiceClient, serverID, tag string) (r CheckResult) {
	url := checkURL(client, serverID, tag)
	_, r.Err = client.Get(url, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	return
}
