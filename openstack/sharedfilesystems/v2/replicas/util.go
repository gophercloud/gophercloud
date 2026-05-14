package replicas

import (
	"github.com/gophercloud/gophercloud"
)

// WaitForStatus will continually poll the resource, checking for a particular
// status. It will do this for the amount of seconds defined.
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

// WaitForState will continually poll the resource, checking for a particular
// state. It will do this for the amount of seconds defined.
func WaitForState(c *gophercloud.ServiceClient, id, state string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := Get(c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.State == state {
			return true, nil
		}

		return false, nil
	})
}
