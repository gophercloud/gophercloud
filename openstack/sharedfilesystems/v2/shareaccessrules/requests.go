package shareaccessrules

import (
	"github.com/gophercloud/gophercloud"
)

// Get retrieves details about a share access rule.
func Get(client *gophercloud.ServiceClient, accessID string) (r GetResult) {
	resp, err := client.Get(getURL(client, accessID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List gets all access rules of a share.
func List(client *gophercloud.ServiceClient, shareID string) (r ListResult) {
	resp, err := client.Get(listURL(client, shareID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
