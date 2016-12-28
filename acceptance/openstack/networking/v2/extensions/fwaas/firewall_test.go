// +build acceptance networking fwaas

package fwaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	layer3 "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/firewalls"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/routerinsertion"
)

func TestFirewallList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := firewalls.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list firewalls: %v", err)
	}

	allFirewalls, err := firewalls.ExtractFirewalls(allPages)
	if err != nil {
		t.Fatalf("Unable to extract firewalls: %v", err)
	}

	for _, firewall := range allFirewalls {
		PrintFirewall(t, &firewall)
	}
}

func TestFirewallCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router.ID)

	rule, err := CreateRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer DeleteRule(t, client, rule.ID)

	PrintRule(t, rule)

	policy, err := CreatePolicy(t, client, rule.ID)
	if err != nil {
		t.Fatalf("Unable to create policy: %v", err)
	}
	defer DeletePolicy(t, client, policy.ID)

	PrintPolicy(t, policy)

	firewall, err := CreateFirewallOnRouter(t, client, policy.ID, router.ID)
	if err != nil {
		t.Fatalf("Unable to create firewall: %v", err)
	}
	defer DeleteFirewall(t, client, firewall.ID)

	PrintFirewall(t, firewall)

	router2, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router2.ID)

	firewallUpdateOpts := firewalls.UpdateOpts{
		PolicyID:    policy.ID,
		Description: "Some firewall description",
	}

	updateOpts := routerinsertion.UpdateOptsExt{
		firewallUpdateOpts,
		[]string{router2.ID},
	}

	_, err = firewalls.Update(client, firewall.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall: %v", err)
	}

	newFirewall, err := firewalls.Get(client, firewall.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get firewall: %v", err)
	}

	PrintFirewall(t, newFirewall)
}
