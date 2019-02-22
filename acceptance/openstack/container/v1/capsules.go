package v1

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/container/v1/capsules"
)

// WaitForCapsuleStatus will poll a capsule's status until it either matches
// the specified status or the status becomes Failed.
func WaitForCapsuleStatus(client *gophercloud.ServiceClient, capsule interface{}, status string) error {
	var uuid string
	if v, ok := capsule.(capsules.Capsule); ok {
		uuid = v.UUID
	}

	if v, ok := capsule.(capsules.CapsuleV132); ok {
		uuid = v.UUID
	}

	return tools.WaitFor(func() (bool, error) {
		v, err := capsules.Get(client, uuid).Extract()
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

		return false, nil
	})
}
