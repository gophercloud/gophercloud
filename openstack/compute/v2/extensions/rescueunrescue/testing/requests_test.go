package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/rescueunrescue"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestRescue(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/3f54d05f-3430-4d80-aa07-63e6af9e2488/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, RescueRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, RescueResult)
	})

	s, err := rescueunrescue.Rescue(context.TODO(), fake.ServiceClient(), "3f54d05f-3430-4d80-aa07-63e6af9e2488", rescueunrescue.RescueOpts{
		AdminPass:      "aUPtawPzE9NU",
		RescueImageRef: "115e5c5b-72f0-4a0a-9067-60706545248c",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "aUPtawPzE9NU", s)
}

func TestUnrescue(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/3f54d05f-3430-4d80-aa07-63e6af9e2488/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, UnrescueRequest)

		w.WriteHeader(http.StatusAccepted)
	})

	err := rescueunrescue.Unrescue(context.TODO(), fake.ServiceClient(), "3f54d05f-3430-4d80-aa07-63e6af9e2488").ExtractErr()
	th.AssertNoErr(t, err)
}
