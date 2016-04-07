// +build fixtures

package policies

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

// PolicyListBody contains the canned body of a policies.List response.
const PolicyListBody = `
{
  "policies_links": [],
  "policies": [
    {
      "name": "webhook policy",
      "links": [
        {
          "href": "https://dfw.autoscale.api.rackspacecloud.com/v1.0/123456/groups/60b15dad-5ea1-43fa-9a12-a1d737b4da07/policies/2b48d247-0282-4b9d-8775-5c4b67e8e649/",
          "rel": "self"
        }
      ],
      "changePercent": 3.3,
      "cooldown": 300,
      "type": "webhook",
      "id": "2b48d247-0282-4b9d-8775-5c4b67e8e649"
    },
    {
      "cooldown": 0,
      "name": "one time",
      "links": [
        {
          "href": "https://dfw.autoscale.api.rackspacecloud.com/v1.0/123456/groups/60b15dad-5ea1-43fa-9a12-a1d737b4da07/policies/c175c31e-65f9-41de-8b15-918420d3253e/",
          "rel": "self"
        }
      ],
      "args": {
        "at": "2020-04-01T23:00:00.000Z"
      },
      "type": "schedule",
      "id": "c175c31e-65f9-41de-8b15-918420d3253e",
      "change": -1
    },
    {
      "cooldown": 0,
      "name": "sunday afternoon",
      "links": [
        {
          "href": "https://dfw.autoscale.api.rackspacecloud.com/v1.0/123456/groups/60b15dad-5ea1-43fa-9a12-a1d737b4da07/policies/e785e3e7-af9e-4f3c-99ae-b80a532e1663/",
          "rel": "self"
        }
      ],
      "args": {
        "cron": "59 15 * * 0"
      },
      "type": "schedule",
      "id": "e785e3e7-af9e-4f3c-99ae-b80a532e1663",
      "desiredCapacity": 2
    }
  ]
}
`

// PolicyCreateBody contains the canned body of a policies.Create response.
const PolicyCreateBody = PolicyListBody

// PolicyCreateRequest contains the canned body of a policies.Create request.
const PolicyCreateRequest = `
[
  {
    "name": "webhook policy",
    "changePercent": 3.3,
    "cooldown": 300,
    "type": "webhook"
  },
  {
    "cooldown": 0,
    "name": "one time",
    "args": {
      "at": "2020-04-01T23:00:00.000Z"
    },
    "type": "schedule",
    "change": -1
  },
  {
    "cooldown": 0,
    "name": "sunday afternoon",
    "args": {
      "cron": "59 15 * * 0"
    },
    "type": "schedule",
    "desiredCapacity": 2
  }
]
`

// PolicyGetBody contains the canned body of a policies.Get response.
const PolicyGetBody = `
{
  "policy": {
    "name": "webhook policy",
    "links": [
      {
        "href": "https://dfw.autoscale.api.rackspacecloud.com/v1.0/123456/groups/60b15dad-5ea1-43fa-9a12-a1d737b4da07/policies/2b48d247-0282-4b9d-8775-5c4b67e8e649/",
        "rel": "self"
      }
    ],
    "changePercent": 3.3,
    "cooldown": 300,
    "type": "webhook",
    "id": "2b48d247-0282-4b9d-8775-5c4b67e8e649"
  }
}
`

// PolicyUpdateRequest contains the canned body of a policies.Update request.
const PolicyUpdateRequest = `
{
  "name": "updated webhook policy",
  "type": "webhook",
  "cooldown": 600,
  "changePercent": 6.6
}
`

var (
	// WebhookPolicy is a Policy corresponding to the first result in PolicyListBody.
	WebhookPolicy = Policy{
		ID:            "2b48d247-0282-4b9d-8775-5c4b67e8e649",
		Name:          "webhook policy",
		Type:          Webhook,
		Cooldown:      300,
		ChangePercent: 3.3,
	}

	// OneTimePolicy is a Policy corresponding to the second result in PolicyListBody.
	OneTimePolicy = Policy{
		ID:     "c175c31e-65f9-41de-8b15-918420d3253e",
		Name:   "one time",
		Type:   Schedule,
		Change: float64(-1),
		Args: map[string]interface{}{
			"at": "2020-04-01T23:00:00.000Z",
		},
	}

	// SundayAfternoonPolicy is a Policy corresponding to the third result in PolicyListBody.
	SundayAfternoonPolicy = Policy{
		ID:              "e785e3e7-af9e-4f3c-99ae-b80a532e1663",
		Name:            "sunday afternoon",
		Type:            Schedule,
		DesiredCapacity: float64(2),
		Args: map[string]interface{}{
			"cron": "59 15 * * 0",
		},
	}
)

// HandlePolicyListSuccessfully sets up the test server to respond to a policies List request.
func HandlePolicyListSuccessfully(t *testing.T) {
	path := "/groups/60b15dad-5ea1-43fa-9a12-a1d737b4da07/policies"

	th.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, PolicyListBody)
	})
}

// HandlePolicyCreateSuccessfully sets up the test server to respond to a policies Create request.
func HandlePolicyCreateSuccessfully(t *testing.T) {
	path := "/groups/60b15dad-5ea1-43fa-9a12-a1d737b4da07/policies"

	th.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")

		th.TestJSONRequest(t, r, PolicyCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, PolicyCreateBody)
	})
}

// HandlePolicyGetSuccessfully sets up the test server to respond to a policies Get request.
func HandlePolicyGetSuccessfully(t *testing.T) {
	groupID := "60b15dad-5ea1-43fa-9a12-a1d737b4da07"
	policyID := "2b48d247-0282-4b9d-8775-5c4b67e8e649"

	path := fmt.Sprintf("/groups/%s/policies/%s", groupID, policyID)

	th.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, PolicyGetBody)
	})
}

// HandlePolicyUpdateSuccessfully sets up the test server to respond to a policies Update request.
func HandlePolicyUpdateSuccessfully(t *testing.T) {
	groupID := "60b15dad-5ea1-43fa-9a12-a1d737b4da07"
	policyID := "2b48d247-0282-4b9d-8775-5c4b67e8e649"

	path := fmt.Sprintf("/groups/%s/policies/%s", groupID, policyID)

	th.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		th.TestJSONRequest(t, r, PolicyUpdateRequest)

		w.WriteHeader(http.StatusNoContent)
	})
}
