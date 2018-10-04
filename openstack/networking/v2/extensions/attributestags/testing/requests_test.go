package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestReplaceAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks/fakeid/tags", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, attributestagsReplaceAllRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, attributestagsReplaceAllResult)
	})

	opts := attributestags.ReplaceAllOpts{
		Tags: []string{"abc", "xyz"},
	}
	res, err := attributestags.ReplaceAll(fake.ServiceClient(), "networks", "fakeid", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, res, []string{"abc", "xyz"})
}
