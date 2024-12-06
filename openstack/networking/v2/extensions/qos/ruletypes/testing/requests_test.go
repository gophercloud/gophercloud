package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/ruletypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListRuleTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListRuleTypesResponse)
	})

	page, err := ruletypes.ListRuleTypes(fake.ServiceClient()).AllPages(context.TODO())
	if err != nil {
		t.Errorf("Failed to list rule types pages: %v", err)
		return
	}

	rules, err := ruletypes.ExtractRuleTypes(page)
	if err != nil {
		t.Errorf("Failed to list rule types: %v", err)
		return
	}

	expected := []ruletypes.RuleType{{Type: "bandwidth_limit"}, {Type: "dscp_marking"}, {Type: "minimum_bandwidth"}}
	th.AssertDeepEquals(t, expected, rules)
}

func TestGetRuleType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/qos/rule-types/bandwidth_limit", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetRuleTypeResponse)
		th.AssertNoErr(t, err)
	})

	r, err := ruletypes.GetRuleType(context.TODO(), fake.ServiceClient(), "bandwidth_limit").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "bandwidth_limit", r.Type)

	th.AssertEquals(t, 2, len(r.Drivers))

	th.AssertEquals(t, "linuxbridge", r.Drivers[0].Name)
	th.AssertEquals(t, 3, len(r.Drivers[0].SupportedParameters))

	th.AssertEquals(t, "openvswitch", r.Drivers[1].Name)
	th.AssertEquals(t, 3, len(r.Drivers[1].SupportedParameters))
}
