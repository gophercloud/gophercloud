package lockunlock

import "github.com/gophercloud/gophercloud"

func actionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

// Lock is the operation responsible for locking a Compute server.
func Lock(client *gophercloud.ServiceClient, id string) (r gophercloud.ErrResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"lock": nil}, nil, nil)
	return
}

// Unlock is the operation responsible for unlocking a Compute server.
func Unlock(client *gophercloud.ServiceClient, id string) (r gophercloud.ErrResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"unlock": nil}, nil, nil)
	return
}
