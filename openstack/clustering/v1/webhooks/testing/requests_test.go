package testing

import (
	"fmt"
	"net/http"
	"testing"

	"encoding/json"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/webhooks"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestWebhooks(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"action": "290c44fa-c60f-4d75-a0eb-87433ba982a3"
			}`)
	})

	result, err := webhooks.Trigger(fake.ServiceClient(), "f93f83f6-762b-41b6-b757-80507834d394", webhooks.TriggerOpts{}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, result, "290c44fa-c60f-4d75-a0eb-87433ba982a3")
}

// Return an invalid type of integer as action id
func TestWebhooksInvalidAction(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"action": 123
			}`)
	})

	_, err := webhooks.Trigger(fake.ServiceClient(), "f93f83f6-762b-41b6-b757-80507834d394", webhooks.TriggerOpts{}).Extract()
	isValid := err.(*json.UnmarshalTypeError) == nil
	th.AssertEquals(t, false, isValid)
}
