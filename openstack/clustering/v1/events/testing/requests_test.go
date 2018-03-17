package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/events"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListEvents(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    	"events": [
        {
					"action": "create",
					"cluster_id": "",
					"id": "2d255b9c-8f36-41a2-a137-c0175ccc29c3",
					"level": "20",
					"oid": "0df0931b-e251-4f2e-8719-4ebfda3627ba",
					"oname": "node009",
					"otype": "NODE",
					"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
					"status": "CREATING",
					"status_reason": "Initializing",
					"timestamp": "2015-03-05T08:53:15Z",
					"user": "a21ded6060534d99840658a777c2af5a"
				}
    	]
		}`)
	})

	count := 0
	timeStamp, _ := time.Parse(time.RFC3339, "2015-03-05T08:53:15Z")
	events.ListDetail(fake.ServiceClient(), events.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := events.ExtractEvents(page)
		if err != nil {
			t.Errorf("Failed to extract events: %v", err)
			return false, err
		}

		expected := []events.Event{
			{
				Action:       "create",
				ClusterUUID:  "",
				ID:           "2d255b9c-8f36-41a2-a137-c0175ccc29c3",
				Level:        "20",
				OidUUID:      "0df0931b-e251-4f2e-8719-4ebfda3627ba",
				OName:        "node009",
				OType:        "NODE",
				ProjectUUID:  "6e18cc2bdbeb48a5b3cad2dc499f6804",
				Status:       "CREATING",
				StatusReason: "Initializing",
				Timestamp:    timeStamp,
				UserUUID:     "a21ded6060534d99840658a777c2af5a",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoEvents(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	events.ListDetail(fake.ServiceClient(), events.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := events.ExtractEvents(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
