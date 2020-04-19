package pauseunpause

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions"
)

// Pause is the operation responsible for pausing a Compute server.
func Pause(client *gophercloud.ServiceClient, id string) (r PauseResult) {
	resp, err := client.Post(extensions.ActionURL(client, id), map[string]interface{}{"pause": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unpause is the operation responsible for unpausing a Compute server.
func Unpause(client *gophercloud.ServiceClient, id string) (r UnpauseResult) {
	resp, err := client.Post(extensions.ActionURL(client, id), map[string]interface{}{"unpause": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
