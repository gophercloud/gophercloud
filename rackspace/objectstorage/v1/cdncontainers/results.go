package cdncontainers

import objectstorage "github.com/rackspace/gophercloud/openstack/objectstorage/v1"

// EnableResult represents the result of a get operation.
type EnableResult struct {
	objectstorage.CommonResult
}
