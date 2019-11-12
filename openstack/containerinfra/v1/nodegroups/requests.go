package nodegroups

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// Get makes a request to the Magnum API to retrieve a node group
// with the given ID/name belonging to the given cluster.
// Use the Extract method of the returned GetResult to extract the
// node group from the result.
func Get(client *gophercloud.ServiceClient, clusterID, nodeGroupID string) (r GetResult) {
	var response *http.Response
	response, r.Err = client.Get(getURL(client, clusterID, nodeGroupID), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	if r.Err == nil {
		r.Header = response.Header
	}
	return
}
