//go:build acceptance || compute || defsecrules

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	dsr "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/defsecrules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestDefSecRulesList(t *testing.T) {
	clients.RequireAdmin(t)
	clients.RequireNovaNetwork(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := dsr.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allDefaultRules, err := dsr.ExtractDefaultRules(allPages)
	th.AssertNoErr(t, err)

	for _, defaultRule := range allDefaultRules {
		tools.PrintResource(t, defaultRule)
	}
}

func TestDefSecRulesCreate(t *testing.T) {
	clients.RequireAdmin(t)
	clients.RequireNovaNetwork(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	defaultRule, err := CreateDefaultRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteDefaultRule(t, client, defaultRule)

	tools.PrintResource(t, defaultRule)
}

func TestDefSecRulesGet(t *testing.T) {
	clients.RequireAdmin(t)
	clients.RequireNovaNetwork(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	defaultRule, err := CreateDefaultRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteDefaultRule(t, client, defaultRule)

	newDefaultRule, err := dsr.Get(context.TODO(), client, defaultRule.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newDefaultRule)
}
