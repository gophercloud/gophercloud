// +build acceptance compute secgroups

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
)

func TestSecGroupsList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := secgroups.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve security groups: %v", err)
	}

	allSecGroups, err := secgroups.ExtractSecurityGroups(allPages)
	if err != nil {
		t.Fatalf("Unable to extract security groups: %v", err)
	}

	for _, secgroup := range allSecGroups {
		printSecurityGroup(t, &secgroup)
	}
}

func TestSecGroupsCreate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	securityGroup, err := createSecurityGroup(t, client)
	if err != nil {
		t.Fatalf("Unable to create security group: %v", err)
	}
	defer deleteSecurityGroup(t, client, securityGroup)
}

func TestSecGroupsUpdate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	securityGroup, err := createSecurityGroup(t, client)
	if err != nil {
		t.Fatalf("Unable to create security group: %v", err)
	}
	defer deleteSecurityGroup(t, client, securityGroup)

	updateOpts := secgroups.UpdateOpts{
		Name:        tools.RandomString("secgroup_", 4),
		Description: tools.RandomString("dec_", 10),
	}
	updatedSecurityGroup, err := secgroups.Update(client, securityGroup.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update security group: %v", err)
	}

	t.Logf("Updated %s's name to %s", updatedSecurityGroup.ID, updatedSecurityGroup.Name)
}

func TestSecGroupsRuleCreate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	securityGroup, err := createSecurityGroup(t, client)
	if err != nil {
		t.Fatalf("Unable to create security group: %v", err)
	}
	defer deleteSecurityGroup(t, client, securityGroup)

	rule, err := createSecurityGroupRule(t, client, securityGroup.ID)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer deleteSecurityGroupRule(t, client, rule)
}

func TestSecGroupsAddGroupToServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	securityGroup, err := createSecurityGroup(t, client)
	if err != nil {
		t.Fatalf("Unable to create security group: %v", err)
	}
	defer deleteSecurityGroup(t, client, securityGroup)

	rule, err := createSecurityGroupRule(t, client, securityGroup.ID)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer deleteSecurityGroupRule(t, client, rule)

	t.Logf("Adding group %s to server %s", securityGroup.ID, server.ID)
	err = secgroups.AddServer(client, server.ID, securityGroup.Name).ExtractErr()
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Unable to add group %s to server %s: %s", securityGroup.ID, server.ID, err)
	}

	t.Logf("Removing group %s from server %s", securityGroup.ID, server.ID)
	err = secgroups.RemoveServer(client, server.ID, securityGroup.Name).ExtractErr()
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Unable to remove group %s from server %s: %s", securityGroup.ID, server.ID, err)
	}
}

func createSecurityGroup(t *testing.T, client *gophercloud.ServiceClient) (secgroups.SecurityGroup, error) {
	createOpts := secgroups.CreateOpts{
		Name:        tools.RandomString("secgroup_", 5),
		Description: "something",
	}

	securityGroup, err := secgroups.Create(client, createOpts).Extract()
	if err != nil {
		return *securityGroup, err
	}

	t.Logf("Created security group: %s", securityGroup.ID)
	return *securityGroup, nil
}

func deleteSecurityGroup(t *testing.T, client *gophercloud.ServiceClient, securityGroup secgroups.SecurityGroup) {
	err := secgroups.Delete(client, securityGroup.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete security group %s: %s", securityGroup.ID, err)
	}

	t.Logf("Deleted security group: %s", securityGroup.ID)
}

func createSecurityGroupRule(t *testing.T, client *gophercloud.ServiceClient, securityGroupID string) (secgroups.Rule, error) {
	createOpts := secgroups.CreateRuleOpts{
		ParentGroupID: securityGroupID,
		FromPort:      tools.RandomInt(80, 89),
		ToPort:        tools.RandomInt(90, 99),
		IPProtocol:    "TCP",
		CIDR:          "0.0.0.0/0",
	}

	rule, err := secgroups.CreateRule(client, createOpts).Extract()
	if err != nil {
		return *rule, err
	}

	t.Logf("Created security group rule: %s", rule.ID)
	return *rule, nil
}

func deleteSecurityGroupRule(t *testing.T, client *gophercloud.ServiceClient, rule secgroups.Rule) {
	err := secgroups.DeleteRule(client, rule.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete rule: %v", err)
	}

	t.Logf("Deleted security group rule: %s", rule.ID)
}

func printSecurityGroup(t *testing.T, securityGroup *secgroups.SecurityGroup) {
	t.Logf("ID: %s", securityGroup.ID)
	t.Logf("Name: %s", securityGroup.Name)
	t.Logf("Description: %s", securityGroup.Description)
	t.Logf("Tenant ID: %s", securityGroup.TenantID)
	t.Logf("Rules:")

	for _, rule := range securityGroup.Rules {
		t.Logf("\tID: %s", rule.ID)
		t.Logf("\tFrom Port: %d", rule.FromPort)
		t.Logf("\tTo Port: %d", rule.ToPort)
		t.Logf("\tIP Protocol: %s", rule.IPProtocol)
		t.Logf("\tIP Range: %s", rule.IPRange.CIDR)
		t.Logf("\tParent Group ID: %s", rule.ParentGroupID)
		t.Logf("\tGroup Tenant ID: %s", rule.Group.TenantID)
		t.Logf("\tGroup Name: %s", rule.Group.Name)
	}
}
