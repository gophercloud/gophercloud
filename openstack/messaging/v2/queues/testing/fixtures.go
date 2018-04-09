package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// QueueName is the name of the queue
var QueueName = "FakeTestQueue"

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
