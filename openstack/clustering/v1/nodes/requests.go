package nodes

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// Get makes a request against senlin to get a details of a node type
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
