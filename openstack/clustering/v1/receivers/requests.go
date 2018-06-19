package receivers

import (
	"github.com/gophercloud/gophercloud"
)

// Delete deletes the specified receiver ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
