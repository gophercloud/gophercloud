package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/receivers"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestUpdateReceiver(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/receivers/6dc6d336e3fc4c0a951b5698cd1236ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"receiver": {
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
		}`)
	})

	createdAt, _ := time.Parse(time.RFC3339, "2015-06-27T05:09:43Z")
	updatedAt := time.Time{}
	expected := receivers.Receiver{
		Action: "CLUSTER_SCALE_OUT",
		Actor: map[string]interface{}{
			"trust_id": []string{
				"6dc6d336e3fc4c0a951b5698cd1236d9",
			},
		},
		Channel: map[string]interface{}{
			"alarm_url": "http://node1:8778/v1/webhooks/e03dd2e5-8f2e-4ec1-8c6a-74ba891e5422/trigger?V=1&count=1",
		},
		Cluster:   "ae63a10b-4a90-452c-aef1-113a0b255ee3",
		CreatedAt: createdAt,
		Domain:    "Default",
		ID:        "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3",
		Name:      "cluster_inflate",
		Params: map[string]interface{}{
			"count": "1",
		},
		Project:   "6e18cc2bdbeb48a5b3cad2dc499f6804",
		Type:      "webhook",
		UpdatedAt: updatedAt,
		User:      "b4ad2d6e18cc2b9c48049f6dbe8a5b3c",
	}

	opts := receivers.UpdateOpts{
		Name:   "cluster_inflate",
		Action: "CLUSTER_SCALE_OUT",
		Params: map[string]interface{}{
			"count": "2",
		},
	}
	actual, err := receivers.Update(fake.ServiceClient(), "6dc6d336e3fc4c0a951b5698cd1236ee", opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, *actual)
}

func TestUpdateReceiversInvalidTimeFloat(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/receivers/6dc6d336e3fc4c0a951b5698cd1236ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"receiver": {
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
				"created_at": 123456789.0,
				"domain": "Default",
				"id": "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3",
				"name": "cluster_inflate",
				"params": {
					"count": "1"
				},
				"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
				"type": "webhook",
				"updated_at": 123456789.0,
				"user": "b4ad2d6e18cc2b9c48049f6dbe8a5b3c"
			}
		}`)
	})
	opts := receivers.UpdateOpts{
		Name: "cluster_inflate",
	}
	_, err := receivers.Update(fake.ServiceClient(), "6dc6d336e3fc4c0a951b5698cd1236ee", opts).Extract()
	th.AssertEquals(t, false, err == nil)
}

func TestUpdateReceiverInvalidTimeString(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/receivers/6dc6d336e3fc4c0a951b5698cd1236ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"receiver": {
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
				"created_at": null,
				"domain": "Default",
				"id": "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3",
				"name": "cluster_inflate",
				"params": {
					"count": "1"
				},
				"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
				"type": "webhook",
				"updated_at": "invalid",
				"user": "b4ad2d6e18cc2b9c48049f6dbe8a5b3c"
			}
		}`)
	})
	opts := receivers.UpdateOpts{
		Name: "cluster_inflate",
	}
	_, err := receivers.Update(fake.ServiceClient(), "6dc6d336e3fc4c0a951b5698cd1236ee", opts).Extract()
	th.AssertEquals(t, false, err == nil)
}
