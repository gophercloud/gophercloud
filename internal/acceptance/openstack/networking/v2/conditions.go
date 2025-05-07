package v2

import (
	"context"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
)

// RequireNeutronExtension will restrict a test to be only run in environments
// with the requested Neutron extension present.
func RequireNeutronExtension(t *testing.T, client *gophercloud.ServiceClient, extension string) {
	_, err := extensions.Get(context.TODO(), client, extension).Extract()
	if err != nil {
		t.Skipf("this test requires %s Neutron extension", extension)
	}
}

// RequirePortForwarding will restrict a test to only be run in environments
// that support port forwarding
func RequirePortForwarding(t *testing.T) {
	if os.Getenv("OS_PORTFORWARDING_ENVIRONMENT") == "" {
		t.Skip("this test requires support for port forwarding")
	}
}
