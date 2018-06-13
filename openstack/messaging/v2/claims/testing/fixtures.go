package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/claims"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// QueueName is the name of the queue
var QueueName = "FakeTestQueue"

var ClaimID = "51db7067821e727dc24df754"

var ClientID = "1234567890"

// CreateClaimResponse is a sample response to a create claim
const CreateClaimResponse = `
{
	"messages": [
		{
			"body": {"event": "BackupStarted"},
			"href": "/v2/queues/FakeTestQueue/messages/51db6f78c508f17ddc924357?claim_id=51db7067821e727dc24df754",
			"age": 57,
			"ttl": 300
		}
	]
}`

// CreateClaimRequest is a sample request to create a claim.
const CreateClaimRequest = `
{
	"ttl": 3600,
	"grace": 3600
}
`

// CreatedClaim is the result of a create request.
var CreatedClaim = []claims.Messages{
	{
		Age:  57,
		Href: fmt.Sprintf("/v2/queues/%s/messages/51db6f78c508f17ddc924357?claim_id=%s", QueueName, ClaimID),
		TTL:  300,
		Body: map[string]interface{}{"event": "BackupStarted"},
	},
}

// HandleCreateSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc(fmt.Sprintf("/v2/queues/%s/claims", QueueName),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestJSONRequest(t, r, CreateClaimRequest)

			w.WriteHeader(http.StatusCreated)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, CreateClaimResponse)
		})
}
