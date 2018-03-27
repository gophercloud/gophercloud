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

// ListQueuesResponse1 is a sample response to a List queues.
const ListQueuesResponse1 = `
{
    "queues":[
        {
            "href":"/v2/queues/london",
            "name":"london",
            "metadata":{
                "_dead_letter_queue":"fake_queue",
                "_dead_letter_queue_messages_ttl":3500,
                "_default_message_delay":25,
                "_default_message_ttl":3700,
                "_max_claim_count":10,
                "_max_messages_post_size":262143,
                "description":"Test queue."
            }
        }
    ],
    "links":[
        {
            "href":"/v2/queues?marker=london",
            "rel":"next"
        }
    ]
}`

// ListQueuesResponse2 is a sample response to a List queues.
const ListQueuesResponse2 = `
{
    "queues":[
		{
            "href":"/v2/queues/beijing",
            "name":"beijing",
            "metadata":{
                "_dead_letter_queue":"fake_queue",
                "_dead_letter_queue_messages_ttl":3500,
                "_default_message_delay":25,
                "_default_message_ttl":3700,
                "_max_claim_count":10,
                "_max_messages_post_size":262143,
                "description":"Test queue."
            }
        }
    ],
    "links":[
        {
            "href":"/v2/queues?marker=beijing",
            "rel":"next"
        }
    ]
}`

// UpdateQueueRequest is a sample request to update a queue.
const UpdateQueueRequest = `
[
    {
        "op": "replace",
        "path": "/metadata/_max_claim_count",
        "value": 10
    }
]`

// UpdateQueueResponse is a sample response to a update queue.
const UpdateQueueResponse = `
{
    "_max_claim_count": 10
}`

// FirstQueue is the first result in a List.
var FirstQueue = queues.Queue{
	Href: "/v2/queues/london",
	Name: "london",
	Metadata: queues.QueueDetails{
		DeadLetterQueue:           "fake_queue",
		DeadLetterQueueMessageTTL: 3500,
		DefaultMessageDelay:       25,
		DefaultMessageTTL:         3700,
		MaxClaimCount:             10,
		MaxMessagesPostSize:       262143,
		Extra:                     map[string]interface{}{"description": "Test queue."},
	},
}

// SecondQueue is the second result in a List.
var SecondQueue = queues.Queue{
	Href: "/v2/queues/beijing",
	Name: "beijing",
	Metadata: queues.QueueDetails{
		DeadLetterQueue:           "fake_queue",
		DeadLetterQueueMessageTTL: 3500,
		DefaultMessageDelay:       25,
		DefaultMessageTTL:         3700,
		MaxClaimCount:             10,
		MaxMessagesPostSize:       262143,
		Extra:                     map[string]interface{}{"description": "Test queue."},
	},
}

// ExpectedQueueSlice is the expected result in a List.
var ExpectedQueueSlice = [][]queues.Queue{{FirstQueue}, {SecondQueue}}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2/queues",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			next := r.RequestURI

			switch next {
			case "/v2/queues?limit=1":
				fmt.Fprintf(w, ListQueuesResponse1)
			case "/v2/queues?marker=london":
				fmt.Fprint(w, ListQueuesResponse2)
			case "/v2/queues?marker=beijing":
				fmt.Fprint(w, `{ "queues": [] }`)
			}
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

// HandleUpdateSuccessfully configures the test server to respond to an Update request.
func HandleUpdateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc(fmt.Sprintf("/v2/queues/%s", QueueName),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestJSONRequest(t, r, UpdateQueueRequest)

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, UpdateQueueResponse)
		})
}
