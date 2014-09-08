package monitoring

import (
	"github.com/rackspace/gophercloud"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
)

type Options struct {
	Endpoint       string
	AuthOptions    gophercloud.AuthOptions
	Authentication identity.AuthResults
}
