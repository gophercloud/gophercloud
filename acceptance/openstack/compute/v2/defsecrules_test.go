// +build acceptance compute defsecrules

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	dsr "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/defsecrules"
)

func TestDefSecRulesList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := dsr.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list default rules: %v", err)
	}

	allDefaultRules, err := dsr.ExtractDefaultRules(allPages)
	if err != nil {
		t.Fatalf("Unable to extract default rules: %v", err)
	}

	for _, defaultRule := range allDefaultRules {
		printDefaultRule(t, &defaultRule)
	}
}

func TestDefSecRulesCreate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	defaultRule, err := createDefaultRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create default rule: %v", err)
	}
	defer deleteDefaultRule(t, client, defaultRule)

	printDefaultRule(t, &defaultRule)
}

func TestDefSecRulesGet(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	defaultRule, err := createDefaultRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create default rule: %v", err)
	}
	defer deleteDefaultRule(t, client, defaultRule)

	newDefaultRule, err := dsr.Get(client, defaultRule.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get default rule %s: %v", defaultRule.ID, err)
	}

	printDefaultRule(t, newDefaultRule)
}

func createDefaultRule(t *testing.T, client *gophercloud.ServiceClient) (dsr.DefaultRule, error) {
	createOpts := dsr.CreateOpts{
		FromPort:   tools.RandomInt(80, 89),
		ToPort:     tools.RandomInt(90, 99),
		IPProtocol: "TCP",
		CIDR:       "0.0.0.0/0",
	}

	defaultRule, err := dsr.Create(client, createOpts).Extract()
	if err != nil {
		return *defaultRule, err
	}

	t.Logf("Created default rule: %s", defaultRule.ID)

	return *defaultRule, nil
}

func deleteDefaultRule(t *testing.T, client *gophercloud.ServiceClient, defaultRule dsr.DefaultRule) {
	err := dsr.Delete(client, defaultRule.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete default rule %s: %v", defaultRule.ID, err)
	}

	t.Logf("Deleted default rule: %s", defaultRule.ID)
}

func printDefaultRule(t *testing.T, defaultRule *dsr.DefaultRule) {
	t.Logf("\tID: %s", defaultRule.ID)
	t.Logf("\tFrom Port: %d", defaultRule.FromPort)
	t.Logf("\tTo Port: %d", defaultRule.ToPort)
	t.Logf("\tIP Protocol: %s", defaultRule.IPProtocol)
	t.Logf("\tIP Range: %s", defaultRule.IPRange.CIDR)
	t.Logf("\tParent Group ID: %s", defaultRule.ParentGroupID)
	t.Logf("\tGroup Tenant ID: %s", defaultRule.Group.TenantID)
	t.Logf("\tGroup Name: %s", defaultRule.Group.Name)
}
