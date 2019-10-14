package quotas

import "github.com/gophercloud/gophercloud"

// Get returns Networking Quotas for a project.
func Get(client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, projectID), &r.Body, nil)
	return
}
