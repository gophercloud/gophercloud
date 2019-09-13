package tags

import "github.com/gophercloud/gophercloud"

// List all tags on a server.
func List(client *gophercloud.ServiceClient, serverID string) (r ListResult) {
	url := listURL(client, serverID)
	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Check if a tag exists on a server.
func Check(client *gophercloud.ServiceClient, serverID, tag string) (r CheckResult) {
	url := checkURL(client, serverID, tag)
	_, r.Err = client.Get(url, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// ReplaceOptsBuilder allows to add additional parameters to the Replace request.
type ReplaceOptsBuilder interface {
	ToTagsReplaceMap() (map[string]interface{}, error)
}

// ReplaceOpts provides options used to replace Tags on a server.
type ReplaceOpts struct {
	Tags []string `json:"tags" required:"true"`
}

// ToTagsReplaceMap formats a ReplaceOpts into the body of the replace request.
func (opts ReplaceOpts) ToTagsReplaceMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Replace replaces all tags on a server.
func Replace(client *gophercloud.ServiceClient, serverID string, opts ReplaceOptsBuilder) (r ReplaceResult) {
	b, err := opts.ToTagsReplaceMap()
	url := replaceURL(client, serverID)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(url, &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
