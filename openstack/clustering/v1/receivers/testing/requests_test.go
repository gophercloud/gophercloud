package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/receivers"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListReceivers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "receivers": [
        {
					"action": "CLUSTER_SCALE_OUT",
					"actor": {
							"trust_id": [
									"6dc6d336e3fc4c0a951b5698cd1236d9"
							]
					},
					"channel": {
							"alarm_url": "http://node1:8778/v1/webhooks/e03dd2e5-8f2e-4ec1-8c6a-74ba891e5422/trigger?V=1&count=1"
					},
					"cluster_id": "ae63a10b-4a90-452c-aef1-113a0b255ee3",
					"created_at": "2015-06-27T05:09:43Z",
					"domain": "Default",
					"id": "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3",
					"name": "cluster_inflate",
					"params": {
							"count": "1"
					},
					"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
					"type": "webhook",
					"updated_at": null,
					"user": "b4ad2d6e18cc2b9c48049f6dbe8a5b3c"
				}
    ]
		}`)
	})

	count := 0

	receivers.ListDetail(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := receivers.ExtractReceivers(page)
		if err != nil {
			t.Errorf("Failed to extract receivers: %v", err)
			return false, err
		}

		createdAt, _ := time.Parse(time.RFC3339, "2015-06-27T05:09:43Z")
		updatedAt := time.Time{}

		expected := []receivers.Receiver{
			{
				Action: "CLUSTER_SCALE_OUT",
				Actor: map[string]interface{}{
					"trust_id": []string{
						"6dc6d336e3fc4c0a951b5698cd1236d9",
					},
				},
				Channel: map[string]interface{}{
					"alarm_url": "http://node1:8778/v1/webhooks/e03dd2e5-8f2e-4ec1-8c6a-74ba891e5422/trigger?V=1&count=1",
				},
				ClusterUUID: "ae63a10b-4a90-452c-aef1-113a0b255ee3",
				CreatedAt:   createdAt,
				DomainUUID:  "Default",
				ID:          "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3",
				Name:        "cluster_inflate",
				Params: map[string]interface{}{
					"count": "1",
				},
				ProjectUUID: "6e18cc2bdbeb48a5b3cad2dc499f6804",
				Type:        "webhook",
				UpdatedAt:   updatedAt,
				UserUUID:    "b4ad2d6e18cc2b9c48049f6dbe8a5b3c",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoReceivers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	receivers.ListDetail(fake.ServiceClient(), receivers.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := receivers.ExtractReceivers(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
