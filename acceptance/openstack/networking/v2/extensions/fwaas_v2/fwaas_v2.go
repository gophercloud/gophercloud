package fwaas_v2

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateRule will create a Firewall Rule with a random source address and
//source port, destination address and port. An error will be returned if
// the rule could not be created.
func CreateRule(t *testing.T, client *gophercloud.ServiceClient) (*rules.Rule, error) {
	ruleName := tools.RandomString("TESTACC-", 8)
	sourceAddress := fmt.Sprintf("192.168.1.%d", tools.RandomInt(1, 100))
	sourcePort := strconv.Itoa(tools.RandomInt(1, 100))
	destinationAddress := fmt.Sprintf("192.168.2.%d", tools.RandomInt(1, 100))
	destinationPort := strconv.Itoa(tools.RandomInt(1, 100))

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

	rule, err := rules.Create(client, createOpts).Extract()
	if err != nil {
		return rule, err
	}

	t.Logf("Rule %s successfully created", ruleName)

	th.AssertEquals(t, rule.Name, ruleName)
	th.AssertEquals(t, rule.Protocol, string(rules.ProtocolTCP))
	th.AssertEquals(t, rule.Action, string(rules.ActionAllow))
	th.AssertEquals(t, rule.SourceIPAddress, sourceAddress)
	th.AssertEquals(t, rule.SourcePort, sourcePort)
	th.AssertEquals(t, rule.DestinationIPAddress, destinationAddress)
	th.AssertEquals(t, rule.DestinationPort, destinationPort)

	return rule, nil
}

// DeleteRule will delete a rule with a specified ID. A fatal error will occur
// if the delete was not successful. This works best when used as a deferred
// function.
func DeleteRule(t *testing.T, client *gophercloud.ServiceClient, ruleID string) {
	t.Logf("Attempting to delete rule: %s", ruleID)

	err := rules.Delete(client, ruleID).ExtractErr()
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

	group, err := groups.Create(client, createOpts).Extract()
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

	err := groups.Delete(client, groupId).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete firewall group %s: %v", groupId, err)
	}

	t.Logf("Deleted firewall group %s", groupId)
}
