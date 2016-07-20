package bays

import (
	"github.com/gophercloud/gophercloud"
)

// Get retrieves a specific bay based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}
