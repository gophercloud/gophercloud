package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"firewall_rule": {
		"protocol": "tcp",
		"description": "ssh rule",
		"destination_ip_address": "192.168.1.0/24",
		"destination_port": "22",
		"name": "ssh_form_any",
		"action": "allow",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61"
	}
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
	"firewall_rule":{
		"protocol": "tcp",
		"description": "ssh rule",
		"source_port": null,
		"source_ip_address": null,
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": "e2a5fb51-698c-4898-87e8-f1eee6b50919",
		"position": 2,
		"destination_port": "22",
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "ssh_form_any",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": true,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
        `)
	})

	options := rules.CreateOpts{
		TenantID:             "80cf934d6ffb4ef5b244f1c512ad1e61",
		Protocol:             rules.ProtocolTCP,
		Description:          "ssh rule",
		DestinationIPAddress: "192.168.1.0/24",
		DestinationPort:      "22",
		Name:                 "ssh_form_any",
		Action:               "allow",
	}

	_, err := rules.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestCreateAnyProtocol(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"firewall_rule": {
		"protocol": null,
		"description": "any to 192.168.1.0/24",
		"destination_ip_address": "192.168.1.0/24",
		"name": "any_to_192.168.1.0/24",
		"action": "allow",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61"
	}
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
	"firewall_rule":{
		"protocol": null,
		"description": "any to 192.168.1.0/24",
		"source_port": null,
		"source_ip_address": null,
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": "e2a5fb51-698c-4898-87e8-f1eee6b50919",
		"position": 2,
		"destination_port": null,
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "any_to_192.168.1.0/24",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": true,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
        `)
	})

	options := rules.CreateOpts{
		TenantID:             "80cf934d6ffb4ef5b244f1c512ad1e61",
		Protocol:             rules.ProtocolAny,
		Description:          "any to 192.168.1.0/24",
		DestinationIPAddress: "192.168.1.0/24",
		Name:                 "any_to_192.168.1.0/24",
		Action:               "allow",
	}

	_, err := rules.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules/f03bd950-6c56-4f5e-a307-45967078f507", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
	"firewall_rule":{
		"protocol": "tcp",
		"description": "ssh rule",
		"source_port": null,
		"source_ip_address": null,
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": "e2a5fb51-698c-4898-87e8-f1eee6b50919",
		"position": 2,
		"destination_port": "22",
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "ssh_form_any",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": true,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
        `)
	})

	rule, err := rules.Get(fake.ServiceClient(), "f03bd950-6c56-4f5e-a307-45967078f507").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "tcp", rule.Protocol)
	th.AssertEquals(t, "ssh rule", rule.Description)
	th.AssertEquals(t, "192.168.1.0/24", rule.DestinationIPAddress)
	th.AssertEquals(t, "e2a5fb51-698c-4898-87e8-f1eee6b50919", rule.FirewallPolicyID)
	th.AssertEquals(t, "22", rule.DestinationPort)
	th.AssertEquals(t, "f03bd950-6c56-4f5e-a307-45967078f507", rule.ID)
	th.AssertEquals(t, "ssh_form_any", rule.Name)
	th.AssertEquals(t, "80cf934d6ffb4ef5b244f1c512ad1e61", rule.TenantID)
	th.AssertEquals(t, true, rule.Enabled)
	th.AssertEquals(t, "allow", rule.Action)
	th.AssertEquals(t, 4, rule.IPVersion)
	th.AssertEquals(t, false, rule.Shared)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules/4ec89077-d057-4a2b-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := rules.Delete(fake.ServiceClient(), "4ec89077-d057-4a2b-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}
