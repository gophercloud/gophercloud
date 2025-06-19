package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/secgroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const (
	serverID = "{serverID}"
	groupID  = "{groupID}"
	ruleID   = "{ruleID}"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockListGroupsResponse(t, fakeServer)

	count := 0

	err := secgroups.List(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := secgroups.ExtractSecurityGroups(page)
		if err != nil {
			t.Errorf("Failed to extract users: %v", err)
			return false, err
		}

		expected := []secgroups.SecurityGroup{
			{
				ID:          groupID,
				Description: "default",
				Name:        "default",
				Rules:       []secgroups.Rule{},
				TenantID:    "openstack",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestListByServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockListGroupsByServerResponse(t, fakeServer, serverID)

	count := 0

	err := secgroups.ListByServer(client.ServiceClient(fakeServer), serverID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := secgroups.ExtractSecurityGroups(page)
		if err != nil {
			t.Errorf("Failed to extract users: %v", err)
			return false, err
		}

		expected := []secgroups.SecurityGroup{
			{
				ID:          groupID,
				Description: "default",
				Name:        "default",
				Rules:       []secgroups.Rule{},
				TenantID:    "openstack",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockCreateGroupResponse(t, fakeServer)

	opts := secgroups.CreateOpts{
		Name:        "test",
		Description: "something",
	}

	group, err := secgroups.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.SecurityGroup{
		ID:          groupID,
		Name:        "test",
		Description: "something",
		TenantID:    "openstack",
		Rules:       []secgroups.Rule{},
	}
	th.AssertDeepEquals(t, expected, group)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockUpdateGroupResponse(t, fakeServer, groupID)

	description := "new_desc"
	opts := secgroups.UpdateOpts{
		Name:        "new_name",
		Description: &description,
	}

	group, err := secgroups.Update(context.TODO(), client.ServiceClient(fakeServer), groupID, opts).Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.SecurityGroup{
		ID:          groupID,
		Name:        "new_name",
		Description: "something",
		TenantID:    "openstack",
		Rules:       []secgroups.Rule{},
	}
	th.AssertDeepEquals(t, expected, group)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockGetGroupsResponse(t, fakeServer, groupID)

	group, err := secgroups.Get(context.TODO(), client.ServiceClient(fakeServer), groupID).Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.SecurityGroup{
		ID:          groupID,
		Description: "default",
		Name:        "default",
		TenantID:    "openstack",
		Rules: []secgroups.Rule{
			{
				FromPort:      80,
				ToPort:        85,
				IPProtocol:    "TCP",
				IPRange:       secgroups.IPRange{CIDR: "0.0.0.0"},
				Group:         secgroups.Group{TenantID: "openstack", Name: "default"},
				ParentGroupID: groupID,
				ID:            ruleID,
			},
		},
	}

	th.AssertDeepEquals(t, expected, group)
}

func TestGetNumericID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	numericGroupID := 12345

	mockGetNumericIDGroupResponse(t, fakeServer, numericGroupID)

	group, err := secgroups.Get(context.TODO(), client.ServiceClient(fakeServer), "12345").Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.SecurityGroup{ID: "12345"}
	th.AssertDeepEquals(t, expected, group)
}

func TestGetNumericRuleID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	numericGroupID := 12345

	mockGetNumericIDGroupRuleResponse(t, fakeServer, numericGroupID)

	group, err := secgroups.Get(context.TODO(), client.ServiceClient(fakeServer), "12345").Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.SecurityGroup{
		ID: "12345",
		Rules: []secgroups.Rule{
			{
				ParentGroupID: "12345",
				ID:            "12345",
			},
		},
	}
	th.AssertDeepEquals(t, expected, group)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockDeleteGroupResponse(t, fakeServer, groupID)

	err := secgroups.Delete(context.TODO(), client.ServiceClient(fakeServer), groupID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAddRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockAddRuleResponse(t, fakeServer)

	opts := secgroups.CreateRuleOpts{
		ParentGroupID: groupID,
		FromPort:      22,
		ToPort:        22,
		IPProtocol:    "TCP",
		CIDR:          "0.0.0.0/0",
	}

	rule, err := secgroups.CreateRule(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.Rule{
		FromPort:      22,
		ToPort:        22,
		Group:         secgroups.Group{},
		IPProtocol:    "TCP",
		ParentGroupID: groupID,
		IPRange:       secgroups.IPRange{CIDR: "0.0.0.0/0"},
		ID:            ruleID,
	}

	th.AssertDeepEquals(t, expected, rule)
}

func TestAddRuleICMPZero(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockAddRuleResponseICMPZero(t, fakeServer)

	opts := secgroups.CreateRuleOpts{
		ParentGroupID: groupID,
		FromPort:      0,
		ToPort:        0,
		IPProtocol:    "ICMP",
		CIDR:          "0.0.0.0/0",
	}

	rule, err := secgroups.CreateRule(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	expected := &secgroups.Rule{
		FromPort:      0,
		ToPort:        0,
		Group:         secgroups.Group{},
		IPProtocol:    "ICMP",
		ParentGroupID: groupID,
		IPRange:       secgroups.IPRange{CIDR: "0.0.0.0/0"},
		ID:            ruleID,
	}

	th.AssertDeepEquals(t, expected, rule)
}

func TestDeleteRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockDeleteRuleResponse(t, fakeServer, ruleID)

	err := secgroups.DeleteRule(context.TODO(), client.ServiceClient(fakeServer), ruleID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAddServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockAddServerToGroupResponse(t, fakeServer, serverID)

	err := secgroups.AddServer(context.TODO(), client.ServiceClient(fakeServer), serverID, "test").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRemoveServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockRemoveServerFromGroupResponse(t, fakeServer, serverID)

	err := secgroups.RemoveServer(context.TODO(), client.ServiceClient(fakeServer), serverID, "test").ExtractErr()
	th.AssertNoErr(t, err)
}
