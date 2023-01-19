package limits

import (
	"github.com/bizflycloud/gophercloud"
)

// Get returns the limits about the currently scoped tenant.
func Get(client *gophercloud.ServiceClient) (r GetResult) {
	url := getURL(client)
	resp, err := client.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
