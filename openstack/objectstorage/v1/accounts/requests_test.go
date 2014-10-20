package accounts

import (
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var metadata = map[string]string{"gophercloud-test": "accounts"}

func TestUpdateAccount(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "X-Account-Meta-Gophercloud-Test", "accounts")

		w.Header().Set("X-Account-Container-Count", "2")
		w.Header().Set("X-Account-Bytes-Used", "14")
		w.Header().Set("X-Account-Meta-Subject", "books")

		w.WriteHeader(http.StatusNoContent)
	})

	options := &UpdateOpts{Metadata: map[string]string{"gophercloud-test": "accounts"}}
	_, err := Update(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestGetAccount(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Set("X-Account-Meta-Foo", "bar")
		w.WriteHeader(http.StatusNoContent)
	})

	expected := map[string]string{"Foo": "bar"}
	actual, err := Get(fake.ServiceClient(), &GetOpts{}).ExtractMetadata()
	if err != nil {
		t.Fatalf("Unable to get account metadata: %v", err)
	}
	th.CheckDeepEquals(t, expected, actual)
}
