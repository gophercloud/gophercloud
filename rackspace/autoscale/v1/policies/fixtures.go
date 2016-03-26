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
      "changePercent": 3,
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
      "change": 2
    }
  ]
}
`

var (
	// WebhookPolicy is a Policy corresponding to the first result in PolicyListBody.
	WebhookPolicy = Policy{
		ID:            "2b48d247-0282-4b9d-8775-5c4b67e8e649",
		Name:          "webhook policy",
		Type:          Webhook,
		Cooldown:      300,
		ChangePercent: 3,
	}

	// OneTimePolicy is a Policy corresponding to the second result in PolicyListBody.
	OneTimePolicy = Policy{
		ID:     "c175c31e-65f9-41de-8b15-918420d3253e",
		Name:   "one time",
		Type:   Schedule,
		Change: -1,
		Args: map[string]interface{}{
			"at": "2020-04-01T23:00:00.000Z",
		},
	}

	// SundayAfternoonPolicy is a Policy corresponding to the third result in PolicyListBody.
	SundayAfternoonPolicy = Policy{
		ID:     "e785e3e7-af9e-4f3c-99ae-b80a532e1663",
		Name:   "sunday afternoon",
		Type:   Schedule,
		Change: 2,
		Args: map[string]interface{}{
			"cron": "59 15 * * 0",
		},
	}
)

// HandlePolicyListSuccessfully sets up the test server to respond to a policies List request.
func HandlePolicyListSuccessfully(t *testing.T) {
	path := "/groups/10eb3219-1b12-4b34-b1e4-e10ee4f24c65/policies"

	th.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, PolicyListBody)
	})
}
