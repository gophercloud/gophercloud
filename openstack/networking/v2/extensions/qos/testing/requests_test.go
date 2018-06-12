package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
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

	versions, err := qos.ListRuleTypes(fake.ServiceClient()).Extract()
	if err != nil {
		t.Errorf("Failed to list rule types: %s", err.Error())
		return
	}

	th.AssertDeepEquals(t, []string{"bandwidth_limit", "dscp_marking", "minimum_bandwidth"}, versions)
}
