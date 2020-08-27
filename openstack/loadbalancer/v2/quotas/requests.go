package quotas

import "github.com/gophercloud/gophercloud"

// Get returns load balancer Quotas for a project.
func Get(client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	resp, err := client.Get(getURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
