package attachments

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// WaitForStatus will continually poll the resource, checking for a particular
// status. It will do this for the amount of seconds defined.
func WaitForStatus(ctx context.Context, c *gophercloud.ServiceClient, id, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := Get(ctx, c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
