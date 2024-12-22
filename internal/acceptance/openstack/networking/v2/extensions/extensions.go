package extensions

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateExternalNetwork will create an external network. An error will be
// returned if the creation failed.
func CreateExternalNetwork(t *testing.T, client *gophercloud.ServiceClient) (*networks.Network, error) {
	networkName := tools.RandomString("TESTACC-", 8)
	networkDescription := tools.RandomString("TESTACC-DESC-", 8)

	t.Logf("Attempting to create external network: %s", networkName)

	adminStateUp := true
	isExternal := true

	networkCreateOpts := networks.CreateOpts{
		Name:         networkName,
		Description:  networkDescription,
		AdminStateUp: &adminStateUp,
	}

	createOpts := external.CreateOptsExt{
		CreateOptsBuilder: networkCreateOpts,
		External:          &isExternal,
	}

	network, err := networks.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return network, err
	}

	t.Logf("Created external network: %s", networkName)

	th.AssertEquals(t, networkName, network.Name)
	th.AssertEquals(t, networkDescription, network.Description)

	return network, nil
}

// CreatePortWithSecurityGroup will create a port with a security group
// attached. An error will be returned if the port could not be created.
func CreatePortWithSecurityGroup(t *testing.T, client *gophercloud.ServiceClient, networkID, subnetID, secGroupID string) (*ports.Port, error) {
	portName := tools.RandomString("TESTACC-", 8)
	portDescription := tools.RandomString("TESTACC-DESC-", 8)
	iFalse := false

	t.Logf("Attempting to create port: %s", portName)

	createOpts := ports.CreateOpts{
		NetworkID:      networkID,
		Name:           portName,
		Description:    portDescription,
		AdminStateUp:   &iFalse,
		FixedIPs:       []ports.IP{{SubnetID: subnetID}},
		SecurityGroups: &[]string{secGroupID},
	}

	port, err := ports.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return port, err
	}

	t.Logf("Successfully created port: %s", portName)

	th.AssertEquals(t, portName, port.Name)
	th.AssertEquals(t, portDescription, port.Description)
	th.AssertEquals(t, networkID, port.NetworkID)

	return port, nil
}

// CreateSecurityGroup will create a security group with a random name.
// An error will be returned if one was failed to be created.
func CreateSecurityGroup(t *testing.T, client *gophercloud.ServiceClient) (*groups.SecGroup, error) {
	secGroupName := tools.RandomString("TESTACC-", 8)
	secGroupDescription := tools.RandomString("TESTACC-DESC-", 8)

	t.Logf("Attempting to create security group: %s", secGroupName)

	createOpts := groups.CreateOpts{
		Name:        secGroupName,
		Description: secGroupDescription,
	}

	secGroup, err := groups.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return secGroup, err
	}

	t.Logf("Created security group: %s", secGroup.ID)

	th.AssertEquals(t, secGroupName, secGroup.Name)
	th.AssertEquals(t, secGroupDescription, secGroup.Description)

	return secGroup, nil
}

// CreateSecurityGroupRule will create a security group rule with a random name
// and random port between 80 and 99.
// An error will be returned if one was failed to be created.
func CreateSecurityGroupRule(t *testing.T, client *gophercloud.ServiceClient, secGroupID string) (*rules.SecGroupRule, error) {
	t.Logf("Attempting to create security group rule in group: %s", secGroupID)

	description := "Rule description"
	fromPort := tools.RandomInt(80, 89)
	toPort := tools.RandomInt(90, 99)

	createOpts := rules.CreateOpts{
		Description:  description,
		Direction:    "ingress",
		EtherType:    "IPv4",
		SecGroupID:   secGroupID,
		PortRangeMin: fromPort,
		PortRangeMax: toPort,
		Protocol:     rules.ProtocolTCP,
	}

	rule, err := rules.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return rule, err
	}

	t.Logf("Created security group rule: %s", rule.ID)

	th.AssertEquals(t, rule.SecGroupID, secGroupID)
	th.AssertEquals(t, rule.Description, description)

	return rule, nil
}

// CreateSecurityGroupRulesBulk will create security group rules with a random name
// and random port between 80 and 99.
// An error will be returned if one was failed to be created.
func CreateSecurityGroupRulesBulk(t *testing.T, client *gophercloud.ServiceClient, secGroupID string) ([]rules.SecGroupRule, error) {
	t.Logf("Attempting to bulk create security group rules in group: %s", secGroupID)

	sgRulesCreateOpts := make([]rules.CreateOpts, 3)
	for i := range 3 {
		description := "Rule description"
		fromPort := tools.RandomInt(1080, 1089)
		toPort := tools.RandomInt(1090, 1099)

		sgRulesCreateOpts[i] = rules.CreateOpts{
			Description:  description,
			Direction:    "ingress",
			EtherType:    "IPv4",
			SecGroupID:   secGroupID,
			PortRangeMin: fromPort,
			PortRangeMax: toPort,
			Protocol:     rules.ProtocolTCP,
		}
	}

	rules, err := rules.CreateBulk(context.TODO(), client, sgRulesCreateOpts).Extract()
	if err != nil {
		return rules, err
	}

	for i, rule := range rules {
		t.Logf("Created security group rule: %s", rule.ID)

		th.AssertEquals(t, sgRulesCreateOpts[i].SecGroupID, rule.SecGroupID)
		th.AssertEquals(t, sgRulesCreateOpts[i].Description, rule.Description)
	}

	return rules, nil
}

// DeleteSecurityGroup will delete a security group of a specified ID.
// A fatal error will occur if the deletion failed. This works best as a
// deferred function
func DeleteSecurityGroup(t *testing.T, client *gophercloud.ServiceClient, secGroupID string) {
	t.Logf("Attempting to delete security group: %s", secGroupID)

	err := groups.Delete(context.TODO(), client, secGroupID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete security group: %v", err)
	}
}

// DeleteSecurityGroupRule will delete a security group rule of a specified ID.
// A fatal error will occur if the deletion failed. This works best as a
// deferred function
func DeleteSecurityGroupRule(t *testing.T, client *gophercloud.ServiceClient, ruleID string) {
	t.Logf("Attempting to delete security group rule: %s", ruleID)

	err := rules.Delete(context.TODO(), client, ruleID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete security group rule: %v", err)
	}
}
