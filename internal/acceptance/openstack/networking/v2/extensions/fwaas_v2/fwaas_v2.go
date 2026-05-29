package fwaas_v2

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/rules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// RemoveRule will remove a rule from the  policy.
func RemoveRule(t *testing.T, client *gophercloud.ServiceClient, policyID string, ruleID string) {
	t.Logf("Attempting to remove rule %s from policy %s", ruleID, policyID)

	_, err := policies.RemoveRule(context.TODO(), client, policyID, ruleID).Extract()
	if err != nil {
		t.Fatalf("Unable to remove rule %s from policy %s: %v", ruleID, policyID, err)
	}
}

// AddRule will add a rule to to a policy.
func AddRule(t *testing.T, client *gophercloud.ServiceClient, policyID string, ruleID string, beforeRuleID string) {
	t.Logf("Attempting to insert rule %s in to policy %s", ruleID, policyID)

	addOpts := policies.InsertRuleOpts{
		ID:           ruleID,
		InsertBefore: beforeRuleID,
	}

	_, err := policies.InsertRule(context.TODO(), client, policyID, addOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to insert rule %s before rule %s in policy %s: %v", ruleID, beforeRuleID, policyID, err)
	}
}

// CreatePolicy will create a Firewall Policy with a random name and given
// rule. An error will be returned if the rule could not be created.
func CreatePolicy(t *testing.T, client *gophercloud.ServiceClient, ruleID string) (*policies.Policy, error) {
	policyName := tools.RandomString("TESTACC-", 8)
	policyDescription := tools.RandomString("TESTACC-DESC-", 8)

	t.Logf("Attempting to create policy %s", policyName)

	createOpts := policies.CreateOpts{
		Name:        policyName,
		Description: policyDescription,
		FirewallRules: []string{
			ruleID,
		},
	}

	policy, err := policies.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return policy, err
	}

	t.Logf("Successfully created policy %s", policyName)

	th.AssertEquals(t, policy.Name, policyName)
	th.AssertEquals(t, policy.Description, policyDescription)
	th.AssertEquals(t, len(policy.Rules), 1)

	return policy, nil
}

// CreateRule will create a Firewall Rule with a random source address and
// source port, destination address and port. An error will be returned if
// the rule could not be created.
func CreateRule(t *testing.T, client *gophercloud.ServiceClient) (*rules.Rule, error) {
	ruleName := tools.RandomString("TESTACC-", 8)
	sourceAddress := fmt.Sprintf("192.168.1.%d", tools.RandomInt(1, 100))
	sourcePortInt := strconv.Itoa(tools.RandomInt(1, 100))
	sourcePort := fmt.Sprintf("%s:%s", sourcePortInt, sourcePortInt)
	destinationAddress := fmt.Sprintf("192.168.2.%d", tools.RandomInt(1, 100))
	destinationPortInt := strconv.Itoa(tools.RandomInt(1, 100))
	destinationPort := fmt.Sprintf("%s:%s", destinationPortInt, destinationPortInt)

	t.Logf("Attempting to create rule %s with source %s:%s and destination %s:%s",
		ruleName, sourceAddress, sourcePort, destinationAddress, destinationPort)

	createOpts := rules.CreateOpts{
		Name:                 ruleName,
		Protocol:             rules.ProtocolTCP,
		Action:               rules.ActionAllow,
		SourceIPAddress:      sourceAddress,
		SourcePort:           sourcePort,
		DestinationIPAddress: destinationAddress,
		DestinationPort:      destinationPort,
	}

	rule, err := rules.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return rule, err
	}

	t.Logf("Rule %s successfully created", ruleName)

	th.AssertEquals(t, rule.Name, ruleName)
	th.AssertEquals(t, rule.Protocol, string(rules.ProtocolTCP))
	th.AssertEquals(t, rule.Action, string(rules.ActionAllow))
	th.AssertEquals(t, rule.SourceIPAddress, sourceAddress)
	th.AssertEquals(t, rule.SourcePort, sourcePortInt)
	th.AssertEquals(t, rule.DestinationIPAddress, destinationAddress)
	th.AssertEquals(t, rule.DestinationPort, destinationPortInt)

	return rule, nil
}

// DeletePolicy will delete a policy with a specified ID. A fatal error will
// occur if the delete was not successful. This works best when used as a
// deferred function.
func DeletePolicy(t *testing.T, client *gophercloud.ServiceClient, policyID string) {
	t.Logf("Attempting to delete policy: %s", policyID)

	err := policies.Delete(context.TODO(), client, policyID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete policy %s: %v", policyID, err)
	}

	t.Logf("Deleted policy: %s", policyID)
}

// DeleteRule will delete a rule with a specified ID. A fatal error will occur
// if the delete was not successful. This works best when used as a deferred
// function.
func DeleteRule(t *testing.T, client *gophercloud.ServiceClient, ruleID string) {
	t.Logf("Attempting to delete rule: %s", ruleID)

	err := rules.Delete(context.TODO(), client, ruleID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete rule %s: %v", ruleID, err)
	}

	t.Logf("Deleted rule: %s", ruleID)
}

// CreateGroup will create a Firewall Group. An error will be returned if the
// firewall group could not be created.
func CreateGroup(t *testing.T, client *gophercloud.ServiceClient) (*groups.Group, error) {

	groupName := tools.RandomString("TESTACC-", 8)
	description := tools.RandomString("TESTACC-", 8)
	adminStateUp := true
	shared := false

	createOpts := groups.CreateOpts{
		Name:         groupName,
		Description:  description,
		AdminStateUp: &adminStateUp,
		Shared:       &shared,
	}

	t.Logf("Attempting to create firewall group %s",
		groupName)

	group, err := groups.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return group, err
	}

	t.Logf("firewall group %s successfully created", groupName)

	th.AssertEquals(t, group.Name, groupName)
	return group, nil
}

// DeleteGroup will delete a group with a specified ID. A fatal error will occur
// if the delete was not successful. This works best when used as a deferred
// function.
func DeleteGroup(t *testing.T, client *gophercloud.ServiceClient, groupId string) {
	t.Logf("Attempting to delete firewall group %s", groupId)

	err := groups.Delete(context.TODO(), client, groupId).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete firewall group %s: %v", groupId, err)
	}

	t.Logf("Deleted firewall group %s", groupId)
}
