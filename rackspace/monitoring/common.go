package monitoring

import (
	"github.com/rackspace/gophercloud/openstack/identity"
)

type Options struct {
	Endpoint       string
	AuthOptions    identity.AuthOptions
	Authentication identity.AuthResults
}
