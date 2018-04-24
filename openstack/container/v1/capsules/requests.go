package capsules

import (
	"github.com/gophercloud/gophercloud"
)

// Get requests details on a single capsule, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return
}
