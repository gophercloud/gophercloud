package datastores

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List will list all available datastore types that instances can use.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, baseURL(client), func(r pagination.PageResult) pagination.Page {
		return DatastorePage{pagination.SinglePageBase(r)}
	})
}

// Get will retrieve the details of a specified datastore type.
func Get(ctx context.Context, client *gophercloud.ServiceClient, datastoreID string) (r GetResult) {
	resp, err := client.GetWithContext(ctx, resourceURL(client, datastoreID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListVersions will list all of the available versions for a specified
// datastore type.
func ListVersions(client *gophercloud.ServiceClient, datastoreID string) pagination.Pager {
	return pagination.NewPager(client, versionsURL(client, datastoreID), func(r pagination.PageResult) pagination.Page {
		return VersionPage{pagination.SinglePageBase(r)}
	})
}

// GetVersion will retrieve the details of a specified datastore version.
func GetVersion(ctx context.Context, client *gophercloud.ServiceClient, datastoreID, versionID string) (r GetVersionResult) {
	resp, err := client.GetWithContext(ctx, versionURL(client, datastoreID, versionID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
