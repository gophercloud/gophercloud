package attributestags

import (
	"github.com/gophercloud/gophercloud"
)

// ReplaceAllOptsBuilder allows extensions to add additional parameters to
// the ReplaceAll request.
type ReplaceAllOptsBuilder interface {
	ToAttributeTagsReplaceAllMap() (map[string]interface{}, error)
}

// ReplaceAllOpts provides options used to create Tags on a Resource
type ReplaceAllOpts struct {
	Tags []string `json:"tags" required:"true"`
}

// ToAttributeTagsReplaceAllMap formats a ReplaceAllOpts into the body of the
// replace request
func (opts ReplaceAllOpts) ToAttributeTagsReplaceAllMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ReplaceAll updates all tags on a resource, replacing any existing tags
func ReplaceAll(client *gophercloud.ServiceClient, resourceType string, resourceID string, opts ReplaceAllOptsBuilder) (r ReplaceAllResult) {
	b, err := opts.ToAttributeTagsReplaceAllMap()
	url := replaceURL(client, resourceType, resourceID)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(url, &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// List all tags on a resource
func List(client *gophercloud.ServiceClient, resourceType string, resourceID string) (r ListResult) {
	url := listURL(client, resourceType, resourceID)
	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
