//go:build acceptance || compute || limits

package v2

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/limits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestLimits(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	limits, err := limits.Get(context.TODO(), client, nil).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, limits)

	th.AssertEquals(t, limits.Absolute.MaxPersonalitySize, 10240)
}

func TestLimitsForTenant(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	// I think this is the easiest way to get the tenant ID while being
	// agnostic to Identity v2 and v3.
	// Technically we're just returning the limits for ourselves, but it's
	// the fact that we're specifying a tenant ID that is important here.
	endpointParts := strings.Split(client.Endpoint, "/")
	tenantID := endpointParts[4]

	getOpts := limits.GetOpts{
		TenantID: tenantID,
	}

	limits, err := limits.Get(context.TODO(), client, getOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, limits)

	th.AssertEquals(t, limits.Absolute.MaxPersonalitySize, 10240)
}
