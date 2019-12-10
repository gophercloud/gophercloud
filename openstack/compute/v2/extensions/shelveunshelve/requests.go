package shelveunshelve

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions"
)

// Shelve is the operation responsible for shelving a Compute server.
func Shelve(client *gophercloud.ServiceClient, id string) (r ShelveResult) {
	_, r.Err = client.Post(extensions.ActionURL(client, id), map[string]interface{}{"shelve": nil}, nil, nil)
	return
}

// ShelveOffload is the operation responsible for Shelve-Offload a Compute server.
func ShelveOffload(client *gophercloud.ServiceClient, id string) (r ShelveOffloadResult) {
	_, r.Err = client.Post(extensions.ActionURL(client, id), map[string]interface{}{"shelveOffload": nil}, nil, nil)
	return
}

// Unshelve is the operation responsible for unshelve a Compute server.
func Unshelve(client *gophercloud.ServiceClient, id string) (r UnshelveResult) {
	_, r.Err = client.Post(extensions.ActionURL(client, id), map[string]interface{}{"unshelve": nil}, nil, nil)
	return
}
