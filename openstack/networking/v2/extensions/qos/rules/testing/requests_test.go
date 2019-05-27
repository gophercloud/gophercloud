package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/qos/policies/501005fa-3b56-4061-aaca-3f24995112e1/bandwidth_limit_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, BandwidthLimitRulesListResult)
	})

	count := 0

	err := rules.BandwidthLimitRulesList(
		fake.ServiceClient(),
		"501005fa-3b56-4061-aaca-3f24995112e1",
		rules.BandwidthLimitRulesListOpts{},
	).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := rules.ExtractBandwidthLimitRules(page)
		if err != nil {
			t.Errorf("Failed to extract bandwith limit rules: %v", err)
			return false, nil
		}

		expected := []rules.BandwidthLimitRule{
			{
				ID:           "30a57f4a-336b-4382-8275-d708babd2241",
				MaxKBps:      3000,
				MaxBurstKBps: 300,
				Direction:    "egress",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
