//go:build acceptance || networking || loadbalancer || loadbalancers

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestLoadbalancersList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	allPages, err := loadbalancers.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)

	for _, lb := range allLoadbalancers {
		tools.PrintResource(t, lb)
	}
}

func TestLoadbalancersListByTags(t *testing.T) {
	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnet(t, netClient, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	// Add "test" tag intentionally to test the "not-tags" parameter. Because "test" tag is also used in other test
	// cases, we use "test" tag to exclude load balancers created by other test case.
	tags := []string{"tag1", "tag2", "test"}
	lb, err := CreateLoadBalancer(t, lbClient, subnet.ID, tags, "", nil)
	th.AssertNoErr(t, err)
	defer DeleteLoadBalancer(t, lbClient, lb.ID)

	tags = []string{"tag1"}
	listOpts := loadbalancers.ListOpts{
		Tags: tags,
	}
	allPages, err := loadbalancers.List(lbClient, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allLoadbalancers))

	tags = []string{"test"}
	listOpts = loadbalancers.ListOpts{
		TagsNot: tags,
	}
	allPages, err = loadbalancers.List(lbClient, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allLoadbalancers, err = loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(allLoadbalancers))

	tags = []string{"tag1", "tag3"}
	listOpts = loadbalancers.ListOpts{
		TagsAny: tags,
	}
	allPages, err = loadbalancers.List(lbClient, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allLoadbalancers, err = loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allLoadbalancers))

	tags = []string{"tag1", "test"}
	listOpts = loadbalancers.ListOpts{
		TagsNotAny: tags,
	}
	allPages, err = loadbalancers.List(lbClient, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allLoadbalancers, err = loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(allLoadbalancers))
}

func TestLoadbalancerHTTPCRUD(t *testing.T) {
	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnet(t, netClient, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	lb, err := CreateLoadBalancer(t, lbClient, subnet.ID, nil, "", nil)
	th.AssertNoErr(t, err)
	defer DeleteLoadBalancer(t, lbClient, lb.ID)

	// Listener
	listener, err := CreateListenerHTTP(t, lbClient, lb)
	th.AssertNoErr(t, err)
	defer DeleteListener(t, lbClient, lb.ID, listener.ID)

	// L7 policy
	tags := []string{"test"}
	policy, err := CreateL7Policy(t, lbClient, listener, lb, tags)
	th.AssertNoErr(t, err)
	defer DeleteL7Policy(t, lbClient, lb.ID, policy.ID)

	tags = []string{"test", "test1"}
	newDescription := ""
	updateL7policyOpts := l7policies.UpdateOpts{
		Description: &newDescription,
		Tags:        &tags,
	}
	_, err = l7policies.Update(context.TODO(), lbClient, policy.ID, updateL7policyOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newPolicy, err := l7policies.Get(context.TODO(), lbClient, policy.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPolicy)

	th.AssertEquals(t, newPolicy.Description, newDescription)
	th.AssertDeepEquals(t, newPolicy.Tags, tags)

	// L7 rule
	tags = []string{"test"}
	rule, err := CreateL7Rule(t, lbClient, newPolicy.ID, lb, tags)
	th.AssertNoErr(t, err)
	defer DeleteL7Rule(t, lbClient, lb.ID, policy.ID, rule.ID)

	allPages, err := l7policies.ListRules(lbClient, policy.ID, l7policies.ListRulesOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allRules, err := l7policies.ExtractRules(allPages)
	th.AssertNoErr(t, err)
	for _, rule := range allRules {
		tools.PrintResource(t, rule)
	}

	tags = []string{"test", "test1"}
	updateL7ruleOpts := l7policies.UpdateRuleOpts{
		RuleType:    l7policies.TypePath,
		CompareType: l7policies.CompareTypeRegex,
		Value:       "/images/special*",
		Tags:        &tags,
	}
	_, err = l7policies.UpdateRule(context.TODO(), lbClient, policy.ID, rule.ID, updateL7ruleOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newRule, err := l7policies.GetRule(context.TODO(), lbClient, newPolicy.ID, rule.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)

	th.AssertDeepEquals(t, newRule.Tags, tags)

	// Pool
	pool, err := CreatePoolHTTP(t, lbClient, lb)
	th.AssertNoErr(t, err)
	defer DeletePool(t, lbClient, lb.ID, pool.ID)

	poolName := ""
	poolDescription := ""
	updatePoolOpts := pools.UpdateOpts{
		Name:        &poolName,
		Description: &poolDescription,
	}
	_, err = pools.Update(context.TODO(), lbClient, pool.ID, updatePoolOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newPool, err := pools.Get(context.TODO(), lbClient, pool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPool)
	th.AssertEquals(t, newPool.Name, poolName)
	th.AssertEquals(t, newPool.Description, poolDescription)

	// Update L7policy to redirect to pool
	newRedirectURL := ""
	updateL7policyOpts = l7policies.UpdateOpts{
		Action:         l7policies.ActionRedirectToPool,
		RedirectPoolID: &newPool.ID,
		RedirectURL:    &newRedirectURL,
	}
	_, err = l7policies.Update(context.TODO(), lbClient, policy.ID, updateL7policyOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newPolicy, err = l7policies.Get(context.TODO(), lbClient, policy.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPolicy)

	th.AssertEquals(t, newPolicy.Description, newDescription)
	th.AssertEquals(t, newPolicy.Action, string(l7policies.ActionRedirectToPool))
	th.AssertEquals(t, newPolicy.RedirectPoolID, newPool.ID)
	th.AssertEquals(t, newPolicy.RedirectURL, newRedirectURL)

	// Workaround for proper delete order
	defer DeleteL7Policy(t, lbClient, lb.ID, policy.ID)
	defer DeleteL7Rule(t, lbClient, lb.ID, policy.ID, rule.ID)

	// Member
	member, err := CreateMember(t, lbClient, lb, pool, subnet.ID, subnet.CIDR)
	th.AssertNoErr(t, err)
	defer DeleteMember(t, lbClient, lb.ID, pool.ID, member.ID)

	monitor, err := CreateMonitor(t, lbClient, lb, pool)
	th.AssertNoErr(t, err)
	defer DeleteMonitor(t, lbClient, lb.ID, monitor.ID)
}

func TestLoadBalancerWithAdditionalVips(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/zed")

	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnetWithCIDR(t, netClient, network.ID, "192.168.1.0/24", "192.168.1.1")
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	additionalSubnet, err := networking.CreateSubnetWithCIDR(t, netClient, network.ID, "192.168.2.0/24", "192.168.2.1")
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, additionalSubnet.ID)

	tags := []string{"test"}
	// Octavia takes care of creating the port for the loadbalancer
	additionalSubnetIP := "192.168.2.207"
	lb, err := CreateLoadBalancer(t, lbClient, subnet.ID, tags, "", []loadbalancers.AdditionalVip{{SubnetID: additionalSubnet.ID, IPAddress: additionalSubnetIP}})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(lb.AdditionalVips))
	th.AssertEquals(t, additionalSubnetIP, lb.AdditionalVips[0].IPAddress)
	defer DeleteLoadBalancer(t, lbClient, lb.ID)
}

func TestLoadbalancersCRUD(t *testing.T) {
	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create QoS policy first as the loadbalancer and its port
	// needs to be deleted before the QoS policy can be deleted
	policy2, err := policies.CreateQoSPolicy(t, netClient)
	th.AssertNoErr(t, err)
	defer policies.DeleteQoSPolicy(t, netClient, policy2.ID)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnetWithCIDR(t, netClient, network.ID, "192.168.1.0/24", "192.168.1.1")
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	policy1, err := policies.CreateQoSPolicy(t, netClient)
	th.AssertNoErr(t, err)
	defer policies.DeleteQoSPolicy(t, netClient, policy1.ID)

	tags := []string{"test"}
	lb, err := CreateLoadBalancer(t, lbClient, subnet.ID, tags, policy1.ID, nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, lb.VipQosPolicyID, policy1.ID)
	defer DeleteLoadBalancer(t, lbClient, lb.ID)

	lbDescription := ""
	updateLoadBalancerOpts := loadbalancers.UpdateOpts{
		Description:    &lbDescription,
		VipQosPolicyID: &policy2.ID,
	}
	_, err = loadbalancers.Update(context.TODO(), lbClient, lb.ID, updateLoadBalancerOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newLB, err := loadbalancers.Get(context.TODO(), lbClient, lb.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newLB)

	th.AssertEquals(t, newLB.Description, lbDescription)
	th.AssertEquals(t, newLB.VipQosPolicyID, policy2.ID)

	lbStats, err := loadbalancers.GetStats(context.TODO(), lbClient, lb.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, lbStats)

	// Because of the time it takes to create a loadbalancer,
	// this test will include some other resources.

	// Listener
	listener, err := CreateListener(t, lbClient, lb)
	th.AssertNoErr(t, err)
	defer DeleteListener(t, lbClient, lb.ID, listener.ID)

	listenerName := ""
	listenerDescription := ""
	updateListenerOpts := listeners.UpdateOpts{
		Name:        &listenerName,
		Description: &listenerDescription,
	}
	_, err = listeners.Update(context.TODO(), lbClient, listener.ID, updateListenerOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newListener, err := listeners.Get(context.TODO(), lbClient, listener.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newListener)

	th.AssertEquals(t, newListener.Name, listenerName)
	th.AssertEquals(t, newListener.Description, listenerDescription)

	listenerStats, err := listeners.GetStats(context.TODO(), lbClient, listener.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, listenerStats)

	// Pool
	pool, err := CreatePool(t, lbClient, lb)
	th.AssertNoErr(t, err)
	defer DeletePool(t, lbClient, lb.ID, pool.ID)

	// Update listener's default pool ID.
	updateListenerOpts = listeners.UpdateOpts{
		DefaultPoolID: &pool.ID,
	}
	_, err = listeners.Update(context.TODO(), lbClient, listener.ID, updateListenerOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newListener, err = listeners.Get(context.TODO(), lbClient, listener.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newListener)

	th.AssertEquals(t, newListener.DefaultPoolID, pool.ID)

	// Remove listener's default pool ID
	emptyPoolID := ""
	updateListenerOpts = listeners.UpdateOpts{
		DefaultPoolID: &emptyPoolID,
	}
	_, err = listeners.Update(context.TODO(), lbClient, listener.ID, updateListenerOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newListener, err = listeners.Get(context.TODO(), lbClient, listener.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newListener)

	th.AssertEquals(t, newListener.DefaultPoolID, "")

	// Member
	member, err := CreateMember(t, lbClient, lb, pool, subnet.ID, subnet.CIDR)
	th.AssertNoErr(t, err)
	defer DeleteMember(t, lbClient, lb.ID, pool.ID, member.ID)

	memberName := ""
	newWeight := tools.RandomInt(11, 100)
	updateMemberOpts := pools.UpdateMemberOpts{
		Name:   &memberName,
		Weight: &newWeight,
	}
	_, err = pools.UpdateMember(context.TODO(), lbClient, pool.ID, member.ID, updateMemberOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMember, err := pools.GetMember(context.TODO(), lbClient, pool.ID, member.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMember)
	th.AssertEquals(t, newMember.Name, memberName)

	newWeight = tools.RandomInt(11, 100)
	memberOpts := pools.BatchUpdateMemberOpts{
		Address:      member.Address,
		ProtocolPort: member.ProtocolPort,
		Weight:       &newWeight,
	}
	batchMembers := []pools.BatchUpdateMemberOpts{memberOpts}
	if err = pools.BatchUpdateMembers(context.TODO(), lbClient, pool.ID, batchMembers).ExtractErr(); err != nil {
		t.Fatalf("Unable to batch update members")
	}

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMember, err = pools.GetMember(context.TODO(), lbClient, pool.ID, member.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMember)

	pool, err = pools.Get(context.TODO(), lbClient, pool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, pool)

	// Monitor
	monitor, err := CreateMonitor(t, lbClient, lb, pool)
	th.AssertNoErr(t, err)
	defer DeleteMonitor(t, lbClient, lb.ID, monitor.ID)

	monName := ""
	newDelay := tools.RandomInt(20, 30)
	newMaxRetriesDown := tools.RandomInt(4, 10)
	updateMonitorOpts := monitors.UpdateOpts{
		Name:           &monName,
		Delay:          newDelay,
		MaxRetriesDown: newMaxRetriesDown,
	}
	_, err = monitors.Update(context.TODO(), lbClient, monitor.ID, updateMonitorOpts).Extract()
	th.AssertNoErr(t, err)

	if err = WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMonitor, err := monitors.Get(context.TODO(), lbClient, monitor.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMonitor)

	th.AssertEquals(t, newMonitor.Name, monName)
	th.AssertEquals(t, newMonitor.Delay, newDelay)
	th.AssertEquals(t, newMonitor.MaxRetriesDown, newMaxRetriesDown)
}

func TestLoadbalancersCascadeCRUD(t *testing.T) {
	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnet(t, netClient, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	tags := []string{"test"}
	lb, err := CreateLoadBalancer(t, lbClient, subnet.ID, tags, "", nil)
	th.AssertNoErr(t, err)
	defer CascadeDeleteLoadBalancer(t, lbClient, lb.ID)

	newLB, err := loadbalancers.Get(context.TODO(), lbClient, lb.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newLB)

	// Because of the time it takes to create a loadbalancer,
	// this test will include some other resources.

	// Listener
	listener, err := CreateListener(t, lbClient, lb)
	th.AssertNoErr(t, err)

	listenerDescription := "Some listener description"
	updateListenerOpts := listeners.UpdateOpts{
		Description: &listenerDescription,
	}
	_, err = listeners.Update(context.TODO(), lbClient, listener.ID, updateListenerOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newListener, err := listeners.Get(context.TODO(), lbClient, listener.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newListener)

	// Pool
	pool, err := CreatePool(t, lbClient, lb)
	th.AssertNoErr(t, err)

	poolDescription := "Some pool description"
	updatePoolOpts := pools.UpdateOpts{
		Description: &poolDescription,
	}
	_, err = pools.Update(context.TODO(), lbClient, pool.ID, updatePoolOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newPool, err := pools.Get(context.TODO(), lbClient, pool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPool)

	// Member
	member, err := CreateMember(t, lbClient, lb, newPool, subnet.ID, subnet.CIDR)
	th.AssertNoErr(t, err)

	newWeight := tools.RandomInt(11, 100)
	updateMemberOpts := pools.UpdateMemberOpts{
		Weight: &newWeight,
	}
	_, err = pools.UpdateMember(context.TODO(), lbClient, pool.ID, member.ID, updateMemberOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMember, err := pools.GetMember(context.TODO(), lbClient, pool.ID, member.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMember)

	// Monitor
	monitor, err := CreateMonitor(t, lbClient, lb, newPool)
	th.AssertNoErr(t, err)

	newDelay := tools.RandomInt(20, 30)
	newMaxRetriesDown := tools.RandomInt(4, 10)
	updateMonitorOpts := monitors.UpdateOpts{
		Delay:          newDelay,
		MaxRetriesDown: newMaxRetriesDown,
	}
	_, err = monitors.Update(context.TODO(), lbClient, monitor.ID, updateMonitorOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(lbClient, lb.ID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMonitor, err := monitors.Get(context.TODO(), lbClient, monitor.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMonitor)

	th.AssertEquals(t, newMonitor.Delay, newDelay)
	th.AssertEquals(t, newMonitor.MaxRetriesDown, newMaxRetriesDown)
}

func TestLoadbalancersFullyPopulatedCRUD(t *testing.T) {
	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnet(t, netClient, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	tags := []string{"test"}
	lb, err := CreateLoadBalancerFullyPopulated(t, lbClient, subnet.ID, tags)
	th.AssertNoErr(t, err)
	defer CascadeDeleteLoadBalancer(t, lbClient, lb.ID)

	newLB, err := loadbalancers.Get(context.TODO(), lbClient, lb.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newLB)
}
