package shareaccessrules

import (
	"github.com/bizflycloud/gophercloud"
)

// Get retrieves details about a share access rule.
func Get(client *gophercloud.ServiceClient, accessID string) (r GetResult) {
	resp, err := client.Get(getURL(client, accessID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
