package quotasets

import (
	"github.com/gophercloud/gophercloud"
)

// Get returns public data about a previously created QuotaSet.
func Get(client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, projectID), &r.Body, nil)
	return
}

// GetDefaults returns public data about the project's default block storage quotas.
func GetDefaults(client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	_, r.Err = client.Get(getDefaultsURL(client, projectID), &r.Body, nil)
	return
}

// GetUsage returns detailed public data about a previously created QuotaSet.
func GetUsage(client *gophercloud.ServiceClient, projectID string) (r GetUsageResult) {
	_, r.Err = client.Get(getUsageURL(client, projectID), &r.Body, nil)
	return
}
