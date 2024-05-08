//go:build acceptance || networking || qos || ruletypes

package ruletypes

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/ruletypes"
)

func TestRuleTypes(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
		return
	}

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "qos")

	page, err := ruletypes.ListRuleTypes(client).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Failed to list rule types pages: %v", err)
		return
	}

	ruleTypes, err := ruletypes.ExtractRuleTypes(page)
	if err != nil {
		t.Fatalf("Failed to list rule types: %v", err)
		return
	}

	tools.PrintResource(t, ruleTypes)

	if len(ruleTypes) > 0 {
		t.Logf("Trying to get rule type: %s", ruleTypes[0].Type)

		ruleType, err := ruletypes.GetRuleType(context.TODO(), client, ruleTypes[0].Type).Extract()
		if err != nil {
			t.Fatalf("Failed to get rule type %s: %s", ruleTypes[0].Type, err)
		}

		tools.PrintResource(t, ruleType)
	}
}
