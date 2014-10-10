package accounts

import (
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var metadata = map[string]string{"gophercloud-test": "accounts"}

func TestUpdateAccount(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "X-Account-Meta-Gophercloud-Test", "accounts")

		w.Header().Set("X-Account-Container-Count", "2")
		w.Header().Set("X-Account-Bytes-Used", "14")
		w.Header().Set("X-Account-Meta-Subject", "books")

		w.WriteHeader(http.StatusNoContent)
	})

	res := Update(fake.ServiceClient(), UpdateOpts{Metadata: metadata})

	metadata := res.ExtractMetadata()
	expected := map[string]string{"Subject": "books"}

	testhelper.AssertDeepEquals(t, expected, metadata)
	testhelper.AssertNoErr(t, res.Err)
}

func TestGetAccount(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "HEAD")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("X-Account-Container-Count", "2")
		w.Header().Set("X-Account-Bytes-Used", "14")
		w.Header().Set("X-Account-Meta-Subject", "books")

		w.WriteHeader(http.StatusNoContent)
	})

	res := Get(fake.ServiceClient(), GetOpts{})

	metadata := res.ExtractMetadata()
	expected := map[string]string{"Subject": "books"}

	testhelper.AssertDeepEquals(t, expected, metadata)
	testhelper.AssertNoErr(t, res.Err)
}
