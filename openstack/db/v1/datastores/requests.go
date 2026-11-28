package datastores

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateVersionOptsBuilder is the top-level interface for create version
// options.
type CreateVersionOptsBuilder interface {
	ToVersionCreateMap() (map[string]any, error)
}

// CreateVersionOpts represents options for registering a datastore version.
//
// This is an admin-only API. When the datastore specified by DatastoreName does
// not exist, Trove creates it automatically.
type CreateVersionOpts struct {
	// The name of the datastore version.
	Name string `json:"name" required:"true"`
	// The name of the datastore.
	DatastoreName string `json:"datastore_name" required:"true"`
	// The datastore manager type.
	DatastoreManager string `json:"datastore_manager" required:"true"`
	// The ID of an image. Either Image or ImageTags must be provided.
	Image string `json:"image,omitempty" or:"ImageTags"`
	// Image tags used to find the latest matching image. Either Image or
	// ImageTags must be provided.
	ImageTags []string `json:"image_tags,omitempty"`
	// Whether the datastore version is enabled.
	Active gophercloud.EnabledState `json:"active" required:"true"`
	// Whether this datastore version is the default for the datastore.
	Default gophercloud.EnabledState `json:"default,omitempty"`
	// Packages associated with the datastore version.
	Packages []string `json:"packages,omitempty"`
	// The database engine version number.
	Version string `json:"version,omitempty"`
}

// ToVersionCreateMap converts a CreateVersionOpts struct into a request body.
func (opts CreateVersionOpts) ToVersionCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "version")
}

// CreateVersion registers a new datastore version. This is an admin-only API.
// Trove automatically creates the datastore when DatastoreName does not already
// exist.
func CreateVersion(ctx context.Context, client *gophercloud.ServiceClient, opts CreateVersionOptsBuilder) (r CreateVersionResult) {
	b, err := opts.ToVersionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createVersionURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAllVersions lists all datastore versions. This is an admin-only API.
func ListAllVersions(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, createVersionURL(client), func(r pagination.PageResult) pagination.Page {
		return VersionPage{pagination.SinglePageBase(r)}
	})
}

// GetVersionByID retrieves a datastore version by ID. This is an admin-only
// API.
func GetVersionByID(ctx context.Context, client *gophercloud.ServiceClient, versionID string) (r GetVersionResult) {
	resp, err := client.Get(ctx, mgmtVersionURL(client, versionID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateVersionOptsBuilder is the top-level interface for update version
// options.
type UpdateVersionOptsBuilder interface {
	ToVersionUpdateMap() (map[string]any, error)
}

// UpdateVersionOpts represents options for updating a datastore version.
type UpdateVersionOpts struct {
	// The name of the datastore version.
	Name string `json:"name,omitempty"`
	// The datastore manager type.
	DatastoreManager string `json:"datastore_manager,omitempty"`
	// The ID of an image.
	Image string `json:"image,omitempty"`
	// Image tags used to find the latest matching image.
	ImageTags []string `json:"image_tags,omitempty"`
	// Whether the datastore version is enabled.
	Active gophercloud.EnabledState `json:"active,omitempty"`
	// Whether this datastore version is the default for the datastore.
	Default gophercloud.EnabledState `json:"default,omitempty"`
}

// ToVersionUpdateMap converts an UpdateVersionOpts struct into a request body.
func (opts UpdateVersionOpts) ToVersionUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// UpdateVersion updates a datastore version. This is an admin-only API.
func UpdateVersion(ctx context.Context, client *gophercloud.ServiceClient, versionID string, opts UpdateVersionOptsBuilder) (r UpdateVersionResult) {
	b, err := opts.ToVersionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, mgmtVersionURL(client, versionID), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteVersion deletes a datastore version. This is an admin-only API.
func DeleteVersion(ctx context.Context, client *gophercloud.ServiceClient, versionID string) (r DeleteVersionResult) {
	resp, err := client.Delete(ctx, mgmtVersionURL(client, versionID), &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List will list all available datastore types that instances can use.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, baseURL(client), func(r pagination.PageResult) pagination.Page {
		return DatastorePage{pagination.SinglePageBase(r)}
	})
}

// Get will retrieve the details of a specified datastore type.
func Get(ctx context.Context, client *gophercloud.ServiceClient, datastoreID string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, datastoreID), &r.Body, nil)
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
	resp, err := client.Get(ctx, versionURL(client, datastoreID, versionID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
