package receivers

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// Get retrieves details of a single receiver. Use Extract to convert its result into a Receiver.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	if r.Err == nil {
		r.Header = result.Header
	}
	return

}
