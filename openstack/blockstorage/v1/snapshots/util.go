package snapshots

import (
	"github.com/rackspace/gophercloud"
)

func WaitForStatus(c *gophercloud.ServiceClient, id, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := Get(c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
