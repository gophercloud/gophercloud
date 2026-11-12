//go:build acceptance || identity || limits

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/limits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetEnforcementModel(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	model, err := limits.GetEnforcementModel(context.TODO(), client).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, model)
}

func TestLimitsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	listOpts := limits.ListOpts{}

	allPages, err := limits.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	_, err = limits.ExtractLimits(allPages)
	th.AssertNoErr(t, err)
}
