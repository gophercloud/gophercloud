package cdncontainers

import (
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestEnableCDNContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("X-Ttl", "259200")
		w.Header().Add("X-Cdn-Enabled", "True")
		w.WriteHeader(http.StatusNoContent)
	})

	options := &EnableOpts{CDNEnabled: true, TTL: 259200}
	actual, err := Enable(fake.ServiceClient(), "testContainer", options).ExtractHeaders()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, actual["X-Ttl"][0], "259200")
	th.CheckEquals(t, actual["X-Cdn-Enabled"][0], "True")
}
