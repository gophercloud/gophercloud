package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/messages"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// QueueName is the name of the queue
var QueueName = "FakeTestQueue"

// CreateMessageResponse is a sample response to a Create message.
const CreateMessageResponse = `
{
  "resources": [
    "/v2/queues/demoqueue/messages/51db6f78c508f17ddc924357",
    "/v2/queues/demoqueue/messages/51db6f78c508f17ddc924358"
  ]
}`

// CreateMessageRequest is a sample request to create a message.
const CreateMessageRequest = `
{
  "messages": [
	{
	  "body": {
		"backup_id": "c378813c-3f0b-11e2-ad92-7823d2b0f3ce",
		"event": "BackupStarted"
	  },
	  "delay": 20,
	  "ttl": 300
	},
	{
	  "body": {
		"current_bytes": "0",
		"event": "BackupProgress",
		"total_bytes": "99614720"
	  }
	}
  ]
}`

// ExpectedResources is the expected result in Create
var ExpectedResources = messages.ResourceList{
	Resources: []string{
		"/v2/queues/demoqueue/messages/51db6f78c508f17ddc924357",
		"/v2/queues/demoqueue/messages/51db6f78c508f17ddc924358",
	},
}

// HandleCreateSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc(fmt.Sprintf("/v2/queues/%s/messages", QueueName),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestJSONRequest(t, r, CreateMessageRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, CreateMessageResponse)
		})
}
