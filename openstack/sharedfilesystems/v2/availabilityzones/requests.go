package availabilityzones

import "github.com/gophercloud/gophercloud"

// List will return the existing availability zones.
func List(client *gophercloud.ServiceClient) (r ListResult) {
	_, r.Err = client.Get(listURL(client), &r.Body, nil)
	return
}
