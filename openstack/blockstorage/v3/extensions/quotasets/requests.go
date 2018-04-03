package quotasets

import (
	"github.com/gophercloud/gophercloud"
)

// Get returns public data about a previously created QuotaSet.
func Get(client *gophercloud.ServiceClient, projectID string) GetResult {
	var res GetResult
	_, res.Err = client.Get(getURL(client, projectID), &res.Body, nil)
	return res
}

// Get returns public data about the project's default block storage quotas.
func GetDefaults(client *gophercloud.ServiceClient, projectID string) GetResult {
	var res GetResult
	_, res.Err = client.Get(getDefaultsURL(client, projectID), &res.Body, nil)
	return res
}

// GetDetail returns detailed public data about a previously created QuotaSet.
func GetDetail(client *gophercloud.ServiceClient, projectID string) GetDetailResult {
	var res GetDetailResult
	_, res.Err = client.Get(getDetailURL(client, projectID), &res.Body, nil)
	return res
}
