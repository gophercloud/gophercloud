// Package v1 contains common functions for creating reservation-based
// resources for use in acceptance tests. See the `*_test.go` files for example
// usages.
package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/hypervisors"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// HypervisorHostname returns the name of a hypervisor that can be enrolled
// into the freepool. Blazar knows a host by the name Nova gives it, so there
// is no way to enrol one without asking Nova first.
func HypervisorHostname(t *testing.T, client *gophercloud.ServiceClient) string {
	allPages, err := hypervisors.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)

	if len(allHypervisors) == 0 {
		t.Skip("no hypervisor available to enrol into the freepool")
	}

	return allHypervisors[0].HypervisorHostname
}

// CreateHost will enrol a compute host into the Blazar freepool. An error will
// be returned if the host was unable to be created.
func CreateHost(t *testing.T, client *gophercloud.ServiceClient, name string) (*hosts.Host, error) {
	t.Logf("Attempting to create host: %s", name)

	createOpts := hosts.CreateOpts{
		Name: name,
		ExtraCapabilities: map[string]any{
			"gpu": "a100",
		},
	}

	host, err := hosts.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return host, err
	}

	t.Logf("Created host: %s", host.ID)

	th.AssertEquals(t, name, host.HypervisorHostname)
	th.AssertEquals(t, "a100", host.ExtraCapabilities["gpu"])

	return host, nil
}

// DeleteHost will remove a host from the Blazar freepool. A fatal error will
// occur if the host could not be deleted.
func DeleteHost(t *testing.T, client *gophercloud.ServiceClient, host *hosts.Host) {
	err := hosts.Delete(context.TODO(), client, host.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete host %s: %v", host.ID, err)
	}

	t.Logf("Deleted host: %s", host.ID)
}
