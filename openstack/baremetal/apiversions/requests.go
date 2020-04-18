package apiversions

import (
	"github.com/gophercloud/gophercloud"
)

// List lists all the API versions available to end users.
func List(client *gophercloud.ServiceClient) (r ListResult) {
	resp, err := client.Get(listURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get will get a specific API version, specified by major ID.
func Get(client *gophercloud.ServiceClient, v string) (r GetResult) {
	resp, err := client.Get(getURL(client, v), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
