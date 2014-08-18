package monitoring

import (
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
)

type Options struct {
	Endpoint       string
	AuthOptions    identity.AuthOptions
	Authentication identity.AuthResults
}
