package shareaccessrules

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retrieves details about a share access rule.
func Get(ctx context.Context, client *gophercloud.ServiceClient, accessID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, accessID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List gets all access rules of a share.
func List(ctx context.Context, client *gophercloud.ServiceClient, shareID string) (r ListResult) {
	resp, err := client.Get(ctx, listURL(client, shareID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMetadataOptsBuilder allows extensions to add additional parameters to the
// UpdateMetadata request.
type UpdateMetadataOptsBuilder interface {
	ToUpdateMetadataMap() (map[string]any, error)
}

// UpdateMetadataOpts contains options for updating share access rule metadata.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share access rule metadata, Update share access rule metadata documentation.
type UpdateMetadataOpts struct {
	// One or more access rule metadata key and value pairs as a dictionary of strings.
	Metadata map[string]string `json:"metadata"`
}

// ToUpdateMetadataMap assembles a request body based on the contents of an
// UpdateMetadataOpts.
func (opts UpdateMetadataOpts) ToUpdateMetadataMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// UpdateMetadata updates the metadata of the specified share access rule.
// Existing metadata items that are not present in the request are left
// untouched; items present in the request are created or overwritten.
// To extract the updated metadata from the response, call the Extract
// method on the MetadataResult.
// Client must have Microversion set; minimum supported microversion for
// UpdateMetadata is 2.45.
func UpdateMetadata(ctx context.Context, client *gophercloud.ServiceClient, accessID string, opts UpdateMetadataOptsBuilder) (r MetadataResult) {
	b, err := opts.ToUpdateMetadataMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, metadataURL(client, accessID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMetadatum deletes a single key-value pair from the metadata of the
// specified share access rule.
// Client must have Microversion set; minimum supported microversion for
// DeleteMetadatum is 2.45.
func DeleteMetadatum(ctx context.Context, client *gophercloud.ServiceClient, accessID, key string) (r DeleteMetadatumResult) {
	resp, err := client.Delete(ctx, metadatumURL(client, accessID, key), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
