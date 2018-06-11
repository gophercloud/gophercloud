package clusterpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// Get retrieves details of a single cluster-policy. Use Extract to convert its
// result into a Node.
func Get(client *gophercloud.ServiceClient, clusterID string, policyID string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, clusterID, policyID), &r.Body, nil)
	return
}
