package openstack

import (
	"testing"

	"github.com/rackspace/gophercloud/testhelper"
)

func TestAuthenticate(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
}
