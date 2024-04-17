package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateListener will create a listener for a given load balancer on a random
// port with a random name. An error will be returned if the listener could not
// be created.
func CreateListener(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancers.LoadBalancer) (*listeners.Listener, error) {
	listenerName := tools.RandomString("TESTACCT-", 8)
	listenerDescription := tools.RandomString("TESTACCT-DESC-", 8)
	listenerPort := tools.RandomInt(1, 100)

	t.Logf("Attempting to create listener %s on port %d", listenerName, listenerPort)

	createOpts := listeners.CreateOpts{
		Name:           listenerName,
		Description:    listenerDescription,
		LoadbalancerID: lb.ID,
		Protocol:       listeners.ProtocolTCP,
		ProtocolPort:   listenerPort,
	}

	listener, err := listeners.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return listener, err
	}

	t.Logf("Successfully created listener %s", listenerName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return listener, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, listener.Name, listenerName)
	th.AssertEquals(t, listener.Description, listenerDescription)
	th.AssertEquals(t, listener.Loadbalancers[0].ID, lb.ID)
	th.AssertEquals(t, listener.Protocol, string(listeners.ProtocolTCP))
	th.AssertEquals(t, listener.ProtocolPort, listenerPort)

	return listener, nil
}

// CreateListenerHTTP will create an HTTP-based listener for a given load
// balancer on a random port with a random name. An error will be returned
// if the listener could not be created.
func CreateListenerHTTP(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancers.LoadBalancer) (*listeners.Listener, error) {
	tlsVersions := []listeners.TLSVersion{}
	tlsVersionsExp := []string(nil)
	listenerName := tools.RandomString("TESTACCT-", 8)
	listenerDescription := tools.RandomString("TESTACCT-DESC-", 8)
	listenerPort := tools.RandomInt(1, 100)

	t.Logf("Attempting to create listener %s on port %d", listenerName, listenerPort)

	headers := map[string]string{
		"X-Forwarded-For": "true",
	}

	// tls_version is only supported in microversion v2.17 introduced in victoria
	if clients.IsCurrentAbove(t, "stable/ussuri") {
		tlsVersions = []listeners.TLSVersion{"TLSv1.2", "TLSv1.3"}
		tlsVersionsExp = []string{"TLSv1.2", "TLSv1.3"}
	}

	createOpts := listeners.CreateOpts{
		Name:           listenerName,
		Description:    listenerDescription,
		LoadbalancerID: lb.ID,
		InsertHeaders:  headers,
		Protocol:       listeners.ProtocolHTTP,
		ProtocolPort:   listenerPort,
		TLSVersions:    tlsVersions,
	}

	listener, err := listeners.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return listener, err
	}

	t.Logf("Successfully created listener %s", listenerName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return listener, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, listener.Name, listenerName)
	th.AssertEquals(t, listener.Description, listenerDescription)
	th.AssertEquals(t, listener.Loadbalancers[0].ID, lb.ID)
	th.AssertEquals(t, listener.Protocol, string(listeners.ProtocolHTTP))
	th.AssertEquals(t, listener.ProtocolPort, listenerPort)
	th.AssertDeepEquals(t, listener.InsertHeaders, headers)
	th.AssertDeepEquals(t, listener.TLSVersions, tlsVersionsExp)

	return listener, nil
}

// CreateLoadBalancer will create a load balancer with a random name on a given
// subnet. An error will be returned if the loadbalancer could not be created.
func CreateLoadBalancer(t *testing.T, client *gophercloud.ServiceClient, subnetID string, tags []string, policyID string, additionalVips []loadbalancers.AdditionalVip) (*loadbalancers.LoadBalancer, error) {
	lbName := tools.RandomString("TESTACCT-", 8)
	lbDescription := tools.RandomString("TESTACCT-DESC-", 8)

	t.Logf("Attempting to create loadbalancer %s on subnet %s", lbName, subnetID)

	createOpts := loadbalancers.CreateOpts{
		Name:           lbName,
		Description:    lbDescription,
		VipSubnetID:    subnetID,
		AdminStateUp:   gophercloud.Enabled,
		AdditionalVips: additionalVips,
	}
	if len(tags) > 0 {
		createOpts.Tags = tags
	}

	if len(policyID) > 0 {
		createOpts.VipQosPolicyID = policyID
	}

	lb, err := loadbalancers.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return lb, err
	}

	t.Logf("Successfully created loadbalancer %s on subnet %s", lbName, subnetID)
	t.Logf("Waiting for loadbalancer %s to become active", lbName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return lb, err
	}

	t.Logf("LoadBalancer %s is active", lbName)

	th.AssertEquals(t, lb.Name, lbName)
	th.AssertEquals(t, lb.Description, lbDescription)
	th.AssertEquals(t, lb.VipSubnetID, subnetID)
	th.AssertEquals(t, lb.AdminStateUp, true)

	if len(tags) > 0 {
		th.AssertDeepEquals(t, lb.Tags, tags)
	}

	if len(policyID) > 0 {
		th.AssertEquals(t, lb.VipQosPolicyID, policyID)
	}

	return lb, nil
}

// CreateLoadBalancerFullyPopulated will create a  fully populated load balancer with a random name on a given
// subnet. It will contain a listener, l7policy, l7rule, pool, member and health monitor.
// An error will be returned if the loadbalancer could not be created.
func CreateLoadBalancerFullyPopulated(t *testing.T, client *gophercloud.ServiceClient, subnetID string, tags []string) (*loadbalancers.LoadBalancer, error) {
	lbName := tools.RandomString("TESTACCT-", 8)
	lbDescription := tools.RandomString("TESTACCT-DESC-", 8)
	listenerName := tools.RandomString("TESTACCT-", 8)
	listenerDescription := tools.RandomString("TESTACCT-DESC-", 8)
	listenerPort := tools.RandomInt(1, 100)
	policyName := tools.RandomString("TESTACCT-", 8)
	policyDescription := tools.RandomString("TESTACCT-DESC-", 8)
	poolName := tools.RandomString("TESTACCT-", 8)
	poolDescription := tools.RandomString("TESTACCT-DESC-", 8)
	memberName := tools.RandomString("TESTACCT-", 8)
	memberPort := tools.RandomInt(100, 1000)
	memberWeight := tools.RandomInt(1, 10)

	t.Logf("Attempting to create fully populated loadbalancer %s on subnet %s which contains listener: %s, l7Policy: %s, pool %s, member %s",
		lbName, subnetID, listenerName, policyName, poolName, memberName)

	createOpts := loadbalancers.CreateOpts{
		Name:         lbName,
		Description:  lbDescription,
		VipSubnetID:  subnetID,
		AdminStateUp: gophercloud.Enabled,
		Listeners: []listeners.CreateOpts{{
			Name:         listenerName,
			Description:  listenerDescription,
			Protocol:     listeners.ProtocolHTTP,
			ProtocolPort: listenerPort,
			DefaultPool: &pools.CreateOpts{
				Name:        poolName,
				Description: poolDescription,
				Protocol:    pools.ProtocolHTTP,
				LBMethod:    pools.LBMethodLeastConnections,
				Members: []pools.CreateMemberOpts{{
					Name:         memberName,
					ProtocolPort: memberPort,
					Weight:       &memberWeight,
					Address:      "1.2.3.4",
					SubnetID:     subnetID,
				}},
				Monitor: &monitors.CreateOpts{
					Delay:          10,
					Timeout:        5,
					MaxRetries:     5,
					MaxRetriesDown: 4,
					Type:           monitors.TypeHTTP,
				},
			},
			L7Policies: []l7policies.CreateOpts{{
				Name:        policyName,
				Description: policyDescription,
				Action:      l7policies.ActionRedirectToURL,
				RedirectURL: "http://www.example.com",
				Rules: []l7policies.CreateRuleOpts{{
					RuleType:    l7policies.TypePath,
					CompareType: l7policies.CompareTypeStartWith,
					Value:       "/api",
				}},
			}},
		}},
	}
	if len(tags) > 0 {
		createOpts.Tags = tags
	}

	lb, err := loadbalancers.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return lb, err
	}

	t.Logf("Successfully created loadbalancer %s on subnet %s", lbName, subnetID)
	t.Logf("Waiting for loadbalancer %s to become active", lbName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return lb, err
	}

	t.Logf("LoadBalancer %s is active", lbName)

	th.AssertEquals(t, lb.Name, lbName)
	th.AssertEquals(t, lb.Description, lbDescription)
	th.AssertEquals(t, lb.VipSubnetID, subnetID)
	th.AssertEquals(t, lb.AdminStateUp, true)

	th.AssertEquals(t, len(lb.Listeners), 1)
	th.AssertEquals(t, lb.Listeners[0].Name, listenerName)
	th.AssertEquals(t, lb.Listeners[0].Description, listenerDescription)
	th.AssertEquals(t, lb.Listeners[0].ProtocolPort, listenerPort)

	th.AssertEquals(t, len(lb.Listeners[0].L7Policies), 1)
	th.AssertEquals(t, lb.Listeners[0].L7Policies[0].Name, policyName)
	th.AssertEquals(t, lb.Listeners[0].L7Policies[0].Description, policyDescription)
	th.AssertEquals(t, lb.Listeners[0].L7Policies[0].Description, policyDescription)
	th.AssertEquals(t, len(lb.Listeners[0].L7Policies[0].Rules), 1)

	th.AssertEquals(t, len(lb.Pools), 1)
	th.AssertEquals(t, lb.Pools[0].Name, poolName)
	th.AssertEquals(t, lb.Pools[0].Description, poolDescription)

	th.AssertEquals(t, len(lb.Pools[0].Members), 1)
	th.AssertEquals(t, lb.Pools[0].Members[0].Name, memberName)
	th.AssertEquals(t, lb.Pools[0].Members[0].ProtocolPort, memberPort)
	th.AssertEquals(t, lb.Pools[0].Members[0].Weight, memberWeight)

	if len(tags) > 0 {
		th.AssertDeepEquals(t, lb.Tags, tags)
	}

	return lb, nil
}

// CreateMember will create a member with a random name, port, address, and
// weight. An error will be returned if the member could not be created.
func CreateMember(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancers.LoadBalancer, pool *pools.Pool, subnetID, subnetCIDR string) (*pools.Member, error) {
	memberName := tools.RandomString("TESTACCT-", 8)
	memberPort := tools.RandomInt(100, 1000)
	memberWeight := tools.RandomInt(1, 10)

	cidrParts := strings.Split(subnetCIDR, "/")
	subnetParts := strings.Split(cidrParts[0], ".")
	memberAddress := fmt.Sprintf("%s.%s.%s.%d", subnetParts[0], subnetParts[1], subnetParts[2], tools.RandomInt(10, 100))

	t.Logf("Attempting to create member %s", memberName)

	createOpts := pools.CreateMemberOpts{
		Name:         memberName,
		ProtocolPort: memberPort,
		Weight:       &memberWeight,
		Address:      memberAddress,
		SubnetID:     subnetID,
	}

	t.Logf("Member create opts: %#v", createOpts)

	member, err := pools.CreateMember(context.TODO(), client, pool.ID, createOpts).Extract()
	if err != nil {
		return member, err
	}

	t.Logf("Successfully created member %s", memberName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return member, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, member.Name, memberName)

	return member, nil
}

// CreateMonitor will create a monitor with a random name for a specific pool.
// An error will be returned if the monitor could not be created.
func CreateMonitor(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancers.LoadBalancer, pool *pools.Pool) (*monitors.Monitor, error) {
	monitorName := tools.RandomString("TESTACCT-", 8)

	t.Logf("Attempting to create monitor %s", monitorName)

	createOpts := monitors.CreateOpts{
		PoolID:         pool.ID,
		Name:           monitorName,
		Delay:          10,
		Timeout:        5,
		MaxRetries:     5,
		MaxRetriesDown: 4,
		Type:           monitors.TypePING,
	}

	monitor, err := monitors.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return monitor, err
	}

	t.Logf("Successfully created monitor: %s", monitorName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return monitor, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, monitor.Name, monitorName)
	th.AssertEquals(t, monitor.Type, monitors.TypePING)
	th.AssertEquals(t, monitor.Delay, 10)
	th.AssertEquals(t, monitor.Timeout, 5)
	th.AssertEquals(t, monitor.MaxRetries, 5)
	th.AssertEquals(t, monitor.MaxRetriesDown, 4)

	return monitor, nil
}

// CreatePool will create a pool with a random name with a specified listener
// and loadbalancer. An error will be returned if the pool could not be
// created.
func CreatePool(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancers.LoadBalancer) (*pools.Pool, error) {
	poolName := tools.RandomString("TESTACCT-", 8)
	poolDescription := tools.RandomString("TESTACCT-DESC-", 8)

	t.Logf("Attempting to create pool %s", poolName)

	createOpts := pools.CreateOpts{
		Name:           poolName,
		Description:    poolDescription,
		Protocol:       pools.ProtocolTCP,
		LoadbalancerID: lb.ID,
		LBMethod:       pools.LBMethodLeastConnections,
	}

	pool, err := pools.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return pool, err
	}

	t.Logf("Successfully created pool %s", poolName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return pool, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, pool.Name, poolName)
	th.AssertEquals(t, pool.Description, poolDescription)
	th.AssertEquals(t, pool.Protocol, string(pools.ProtocolTCP))
	th.AssertEquals(t, pool.Loadbalancers[0].ID, lb.ID)
	th.AssertEquals(t, pool.LBMethod, string(pools.LBMethodLeastConnections))

	return pool, nil
}

// CreatePoolHTTP will create an HTTP-based pool with a random name with a
// specified listener and loadbalancer. An error will be returned if the pool
// could not be created.
func CreatePoolHTTP(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancers.LoadBalancer) (*pools.Pool, error) {
	poolName := tools.RandomString("TESTACCT-", 8)
	poolDescription := tools.RandomString("TESTACCT-DESC-", 8)

	t.Logf("Attempting to create pool %s", poolName)

	createOpts := pools.CreateOpts{
		Name:           poolName,
		Description:    poolDescription,
		Protocol:       pools.ProtocolHTTP,
		LoadbalancerID: lb.ID,
		LBMethod:       pools.LBMethodLeastConnections,
	}

	pool, err := pools.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return pool, err
	}

	t.Logf("Successfully created pool %s", poolName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return pool, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, pool.Name, poolName)
	th.AssertEquals(t, pool.Description, poolDescription)
	th.AssertEquals(t, pool.Protocol, string(pools.ProtocolHTTP))
	th.AssertEquals(t, pool.Loadbalancers[0].ID, lb.ID)
	th.AssertEquals(t, pool.LBMethod, string(pools.LBMethodLeastConnections))

	return pool, nil
}

// CreateL7Policy will create a l7 policy with a random name with a specified listener
// and loadbalancer. An error will be returned if the l7 policy could not be
// created.
func CreateL7Policy(t *testing.T, client *gophercloud.ServiceClient, listener *listeners.Listener, lb *loadbalancers.LoadBalancer, tags []string) (*l7policies.L7Policy, error) {
	policyName := tools.RandomString("TESTACCT-", 8)
	policyDescription := tools.RandomString("TESTACCT-DESC-", 8)

	t.Logf("Attempting to create l7 policy %s", policyName)

	createOpts := l7policies.CreateOpts{
		Name:        policyName,
		Description: policyDescription,
		ListenerID:  listener.ID,
		Action:      l7policies.ActionRedirectToURL,
		RedirectURL: "http://www.example.com",
		Tags:        tags,
	}

	policy, err := l7policies.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return policy, err
	}

	t.Logf("Successfully created l7 policy %s", policyName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return policy, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, policy.Name, policyName)
	th.AssertEquals(t, policy.Description, policyDescription)
	th.AssertEquals(t, policy.ListenerID, listener.ID)
	th.AssertEquals(t, policy.Action, string(l7policies.ActionRedirectToURL))
	th.AssertEquals(t, policy.RedirectURL, "http://www.example.com")
	th.AssertDeepEquals(t, policy.Tags, tags)

	return policy, nil
}

// CreateL7Rule creates a l7 rule for specified l7 policy.
func CreateL7Rule(t *testing.T, client *gophercloud.ServiceClient, policyID string, lb *loadbalancers.LoadBalancer, tags []string) (*l7policies.Rule, error) {
	t.Logf("Attempting to create l7 rule for policy %s", policyID)

	createOpts := l7policies.CreateRuleOpts{
		RuleType:    l7policies.TypePath,
		CompareType: l7policies.CompareTypeStartWith,
		Value:       "/api",
		Tags:        tags,
	}

	rule, err := l7policies.CreateRule(context.TODO(), client, policyID, createOpts).Extract()
	if err != nil {
		return rule, err
	}

	t.Logf("Successfully created l7 rule for policy %s", policyID)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE"); err != nil {
		return rule, fmt.Errorf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	th.AssertEquals(t, rule.RuleType, string(l7policies.TypePath))
	th.AssertEquals(t, rule.CompareType, string(l7policies.CompareTypeStartWith))
	th.AssertEquals(t, rule.Value, "/api")
	th.AssertDeepEquals(t, rule.Tags, tags)

	return rule, nil
}

// DeleteL7Policy will delete a specified l7 policy. A fatal error will occur if
// the l7 policy could not be deleted. This works best when used as a deferred
// function.
func DeleteL7Policy(t *testing.T, client *gophercloud.ServiceClient, lbID, policyID string) {
	t.Logf("Attempting to delete l7 policy %s", policyID)

	if err := l7policies.Delete(context.TODO(), client, policyID).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete l7 policy: %v", err)
		}
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	t.Logf("Successfully deleted l7 policy %s", policyID)
}

// DeleteL7Rule will delete a specified l7 rule. A fatal error will occur if
// the l7 rule could not be deleted. This works best when used as a deferred
// function.
func DeleteL7Rule(t *testing.T, client *gophercloud.ServiceClient, lbID, policyID, ruleID string) {
	t.Logf("Attempting to delete l7 rule %s", ruleID)

	if err := l7policies.DeleteRule(context.TODO(), client, policyID, ruleID).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete l7 rule: %v", err)
		}
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	t.Logf("Successfully deleted l7 rule %s", ruleID)
}

// DeleteListener will delete a specified listener. A fatal error will occur if
// the listener could not be deleted. This works best when used as a deferred
// function.
func DeleteListener(t *testing.T, client *gophercloud.ServiceClient, lbID, listenerID string) {
	t.Logf("Attempting to delete listener %s", listenerID)

	if err := listeners.Delete(context.TODO(), client, listenerID).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete listener: %v", err)
		}
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	t.Logf("Successfully deleted listener %s", listenerID)
}

// DeleteMember will delete a specified member. A fatal error will occur if the
// member could not be deleted. This works best when used as a deferred
// function.
func DeleteMember(t *testing.T, client *gophercloud.ServiceClient, lbID, poolID, memberID string) {
	t.Logf("Attempting to delete member %s", memberID)

	if err := pools.DeleteMember(context.TODO(), client, poolID, memberID).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete member: %s", memberID)
		}
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	t.Logf("Successfully deleted member %s", memberID)
}

// DeleteLoadBalancer will delete a specified loadbalancer. A fatal error will
// occur if the loadbalancer could not be deleted. This works best when used
// as a deferred function.
func DeleteLoadBalancer(t *testing.T, client *gophercloud.ServiceClient, lbID string) {
	t.Logf("Attempting to delete loadbalancer %s", lbID)

	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: false,
	}

	if err := loadbalancers.Delete(context.TODO(), client, lbID, deleteOpts).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete loadbalancer: %v", err)
		}
	}

	t.Logf("Waiting for loadbalancer %s to delete", lbID)

	if err := WaitForLoadBalancerState(client, lbID, "DELETED"); err != nil {
		t.Fatalf("Loadbalancer did not delete in time: %s", err)
	}

	t.Logf("Successfully deleted loadbalancer %s", lbID)
}

// CascadeDeleteLoadBalancer will perform a cascading delete on a loadbalancer.
// A fatal error will occur if the loadbalancer could not be deleted. This works
// best when used as a deferred function.
func CascadeDeleteLoadBalancer(t *testing.T, client *gophercloud.ServiceClient, lbID string) {
	t.Logf("Attempting to cascade delete loadbalancer %s", lbID)

	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: true,
	}

	if err := loadbalancers.Delete(context.TODO(), client, lbID, deleteOpts).ExtractErr(); err != nil {
		t.Fatalf("Unable to cascade delete loadbalancer: %v", err)
	}

	t.Logf("Waiting for loadbalancer %s to cascade delete", lbID)

	if err := WaitForLoadBalancerState(client, lbID, "DELETED"); err != nil {
		t.Fatalf("Loadbalancer did not delete in time.")
	}

	t.Logf("Successfully deleted loadbalancer %s", lbID)
}

// DeleteMonitor will delete a specified monitor. A fatal error will occur if
// the monitor could not be deleted. This works best when used as a deferred
// function.
func DeleteMonitor(t *testing.T, client *gophercloud.ServiceClient, lbID, monitorID string) {
	t.Logf("Attempting to delete monitor %s", monitorID)

	if err := monitors.Delete(context.TODO(), client, monitorID).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete monitor: %v", err)
		}
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	t.Logf("Successfully deleted monitor %s", monitorID)
}

// DeletePool will delete a specified pool. A fatal error will occur if the
// pool could not be deleted. This works best when used as a deferred function.
func DeletePool(t *testing.T, client *gophercloud.ServiceClient, lbID, poolID string) {
	t.Logf("Attempting to delete pool %s", poolID)

	if err := pools.Delete(context.TODO(), client, poolID).ExtractErr(); err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatalf("Unable to delete pool: %v", err)
		}
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE"); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active: %s", err)
	}

	t.Logf("Successfully deleted pool %s", poolID)
}

// WaitForLoadBalancerState will wait until a loadbalancer reaches a given state.
func WaitForLoadBalancerState(client *gophercloud.ServiceClient, lbID, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		current, err := loadbalancers.Get(ctx, client, lbID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) && status == "DELETED" {
				return true, nil
			}
			return false, err
		}

		if current.ProvisioningStatus == status {
			return true, nil
		}

		if current.ProvisioningStatus == "ERROR" {
			return false, fmt.Errorf("Load balancer is in ERROR state")
		}

		return false, nil
	})
}

func CreateFlavorProfile(t *testing.T, client *gophercloud.ServiceClient) (*flavorprofiles.FlavorProfile, error) {
	flavorProfileName := tools.RandomString("TESTACCT-", 8)
	flavorProfileDriver := "amphora"
	flavorProfileData := "{\"loadbalancer_topology\": \"SINGLE\"}"

	createOpts := flavorprofiles.CreateOpts{
		Name:         flavorProfileName,
		ProviderName: flavorProfileDriver,
		FlavorData:   flavorProfileData,
	}

	flavorProfile, err := flavorprofiles.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return flavorProfile, err
	}

	t.Logf("Successfully created flavorprofile %s", flavorProfileName)

	th.AssertEquals(t, flavorProfileName, flavorProfile.Name)
	th.AssertEquals(t, flavorProfileDriver, flavorProfile.ProviderName)
	th.AssertEquals(t, flavorProfileData, flavorProfile.FlavorData)

	return flavorProfile, nil
}

func DeleteFlavorProfile(t *testing.T, client *gophercloud.ServiceClient, flavorProfile *flavorprofiles.FlavorProfile) {
	err := flavorprofiles.Delete(context.TODO(), client, flavorProfile.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete flavorprofile: %v", err)
	}

	t.Logf("Successfully deleted flavorprofile %s", flavorProfile.Name)
}

func CreateFlavor(t *testing.T, client *gophercloud.ServiceClient, flavorProfile *flavorprofiles.FlavorProfile) (*flavors.Flavor, error) {
	flavorName := tools.RandomString("TESTACCT-", 8)
	description := tools.RandomString("TESTACCT-desc-", 32)

	createOpts := flavors.CreateOpts{
		Name:            flavorName,
		Description:     description,
		FlavorProfileId: flavorProfile.ID,
		Enabled:         true,
	}

	flavor, err := flavors.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return flavor, err
	}

	t.Logf("Successfully created flavor %s with flavorprofile %s", flavor.Name, flavorProfile.Name)

	th.AssertEquals(t, flavorName, flavor.Name)
	th.AssertEquals(t, description, flavor.Description)
	th.AssertEquals(t, flavorProfile.ID, flavor.FlavorProfileId)
	th.AssertEquals(t, true, flavor.Enabled)

	return flavor, nil
}

func DeleteFlavor(t *testing.T, client *gophercloud.ServiceClient, flavor *flavors.Flavor) {
	err := flavors.Delete(context.TODO(), client, flavor.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete flavor: %v", err)
	}

	t.Logf("Successfully deleted flavor %s", flavor.Name)
}
