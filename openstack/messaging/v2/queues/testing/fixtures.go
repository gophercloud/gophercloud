package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/queues"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// QueueName is the name of the queue
var QueueName = "FakeTestQueue"

// ClientID is a required parameter used the the Header.
var ClientID = "1234567890"

// CreateQueueRequest is a sample request to create a queue.
const CreateQueueRequest = `
{
    "_max_messages_post_size": 262144,
    "_default_message_ttl": 3600,
    "_default_message_delay": 30,
    "_dead_letter_queue": "dead_letter",
    "_dead_letter_queue_messages_ttl": 3600,
    "_max_claim_count": 10,
    "description": "Queue for unit testing."
}`

// ListQueuesResponse is a sample response to a List queues.
const ListQueuesResponse = `
{
   "queues":[
      {
         "href": "/v2/queues/beijing",
         "name": "beijing"
      },
      {
         "href": "/v2/queues/london",
         "name": "london"
      }
   ]
}`

// FirstQueue is the first result in a List.
var FirstQueue = queues.Queue{
	Href: "/v2/queues/beijing",
	Name: "beijing",
}

// SecondQueue is the second result in a List.
var SecondQueue = queues.Queue{
	Href: "/v2/queues/london",
	Name: "london",
}

// ExpectedQueueSlice is the expected result in a List.
var ExpectedQueueSlice = []queues.Queue{FirstQueue, SecondQueue}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2/queues",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, ListQueuesResponse)
		})
}

// HandleCreateSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc(fmt.Sprintf("/v2/queues/%s", QueueName),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestJSONRequest(t, r, CreateQueueRequest)

			w.WriteHeader(http.StatusNoContent)
		})
}
