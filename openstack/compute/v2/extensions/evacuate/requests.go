package evacuate

import (
	"github.com/gophercloud/gophercloud"
)

// Evacuate will Evacuate a failed instance to another host.
func Evacuate(client *gophercloud.ServiceClient, id string) (r EvacuateResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"evacuate": nil}, nil, nil)
	return
}
