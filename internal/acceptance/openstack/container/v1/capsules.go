package v1

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/container/v1/capsules"
)

// WaitForCapsuleStatus will poll a capsule's status until it either matches
// the specified status or the status becomes Failed.
func WaitForCapsuleStatus(client *gophercloud.ServiceClient, uuid, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		v, err := capsules.Get(ctx, client, uuid).Extract()
		if err != nil {
			return false, err
		}

		var newStatus string
		if capsule, ok := v.(*capsules.Capsule); ok {
			newStatus = capsule.Status
		}

		if capsule, ok := v.(*capsules.CapsuleV132); ok {
			newStatus = capsule.Status
		}

		fmt.Println(status)
		fmt.Println(newStatus)

		if newStatus == status {
			// Success!
			return true, nil
		}

		if newStatus == "Failed" {
			return false, fmt.Errorf("Capsule in FAILED state")
		}

		if newStatus == "Error" {
			return false, fmt.Errorf("Capsule in ERROR state")
		}

		return false, nil
	})
}
