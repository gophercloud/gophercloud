package nodes

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// Delete deletes the specified node ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
