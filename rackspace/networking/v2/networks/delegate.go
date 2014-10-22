package networks

import (
  "github.com/rackspace/gophercloud"
  os "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
)

// Delete accepts a unique ID and deletes the network associated with it.
func Delete(c *gophercloud.ServiceClient, id string) os.DeleteResult {
  return os.Delete(c, id)
}
