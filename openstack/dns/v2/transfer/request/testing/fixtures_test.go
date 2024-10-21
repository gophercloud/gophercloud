package testing

import (
	"fmt"
	"net/http"
	s "strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	transferRequests "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/request"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput is a sample response to a List call.
const ListOutput = `
{
    "transfer_requests": [
        {
            "id": "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3",
            "zone_id": "a6a8515c-5d80-48c0-955b-fde631b59791",
            "project_id": "4335d1f0-f793-11e2-b778-0800200c9a66",
            "target_project_id": "05d98711-b3a1-4264-a395-f46383671ee6",
            "description": "This is a first example zone transfer request.",
            "key": "KJSDH23Z",
            "status": "ACTIVE",
            "zone_name": "example1.org.",
            "created_at": "2020-10-12T08:38:58.000000",
            "links": {
                "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/a86dba58-0043-4cc6-a1bb-69d5e86f3ca3"
            }
        },
        {
            "id": "34c4561c-9205-4386-9df5-167436f5a222",
            "zone_id": "572ba08c-d929-4c70-8e42-03824bb24ca2",
            "project_id": "4335d1f0-f793-11e2-b778-0800200c9a66",
            "target_project_id": "05d98711-b3a1-4264-a395-f46383671ee6",
            "description": "This is second example zone transfer request.",
            "key": "KSDFJ22H",
            "status": "ACTIVE",
            "zone_name": "example2.org.",
            "created_at": "2020-10-12T09:38:58.000000",
            "updated_at": "2020-10-12T10:38:58.000000",
            "links": {
                "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/34c4561c-9205-4386-9df5-167436f5a222"
            }
        }
    ],
    "links": {
        "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests"
    }
}
`

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
    "id": "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3",
    "zone_id": "a6a8515c-5d80-48c0-955b-fde631b59791",
    "project_id": "4335d1f0-f793-11e2-b778-0800200c9a66",
    "target_project_id": "05d98711-b3a1-4264-a395-f46383671ee6",
    "description": "This is a first example zone transfer request.",
    "key": "KJSDH23Z",
    "status": "ACTIVE",
    "zone_name": "example1.org.",
    "created_at": "2020-10-12T08:38:58.000000",
    "links": {
        "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/a86dba58-0043-4cc6-a1bb-69d5e86f3ca3"
    }
}
`

// FirstTransferRequest is the first result in ListOutput
var FirstTransferRequestCreatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-10-12T08:38:58.000000")
var FirstTransferRequest = transferRequests.TransferRequest{
	ID:              "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3",
	ZoneID:          "a6a8515c-5d80-48c0-955b-fde631b59791",
	ProjectID:       "4335d1f0-f793-11e2-b778-0800200c9a66",
	TargetProjectID: "05d98711-b3a1-4264-a395-f46383671ee6",
	ZoneName:        "example1.org.",
	Key:             "KJSDH23Z",
	Description:     "This is a first example zone transfer request.",
	Status:          "ACTIVE",
	CreatedAt:       FirstTransferRequestCreatedAt,
	Links: map[string]any{
		"self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/a86dba58-0043-4cc6-a1bb-69d5e86f3ca3",
	},
}

// SecondTransferRequest is the second result in ListOutput
var SecondTransferRequestCreatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-10-12T09:38:58.000000")
var SecondTransferRequestUpdatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-10-12T10:38:58.000000")
var SecondTransferRequest = transferRequests.TransferRequest{
	ID:              "34c4561c-9205-4386-9df5-167436f5a222",
	ZoneID:          "572ba08c-d929-4c70-8e42-03824bb24ca2",
	ProjectID:       "4335d1f0-f793-11e2-b778-0800200c9a66",
	TargetProjectID: "05d98711-b3a1-4264-a395-f46383671ee6",
	ZoneName:        "example2.org.",
	Key:             "KSDFJ22H",
	Description:     "This is second example zone transfer request.",
	Status:          "ACTIVE",
	CreatedAt:       SecondTransferRequestCreatedAt,
	UpdatedAt:       SecondTransferRequestUpdatedAt,
	Links: map[string]any{
		"self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/34c4561c-9205-4386-9df5-167436f5a222",
	},
}

// ExpectedTransferRequestsSlice is the slice of results that should be parsed
// from ListOutput, in the expected order.
var ExpectedTransferRequestsSlice = []transferRequests.TransferRequest{FirstTransferRequest, SecondTransferRequest}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_requests"
	th.Mux.HandleFunc(baseURL,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, ListOutput)
		})
}

// HandleGetSuccessfully configures the test server to respond to a List request.
func HandleGetSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_requests"
	th.Mux.HandleFunc(s.Join([]string{baseURL, FirstTransferRequest.ID}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, GetOutput)
		})
}

// CreateTransferRequest is a sample request to create a zone.
const CreateTransferRequest = `
{
    "target_project_id": "05d98711-b3a1-4264-a395-f46383671ee6",
    "description": "This is a first example zone transfer request."
}
`

// CreateZoneResponse is a sample response to a create request.
const CreateTransferRequestResponse = `
{
    "id": "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3",
    "zone_id": "a6a8515c-5d80-48c0-955b-fde631b59791",
    "project_id": "4335d1f0-f793-11e2-b778-0800200c9a66",
    "target_project_id": "05d98711-b3a1-4264-a395-f46383671ee6",
    "description": "This is a first example zone transfer request.",
    "key": "KJSDH23Z",
    "status": "ACTIVE",
    "zone_name": "example1.org.",
    "created_at": "2020-10-12T08:38:58.000000",
    "links": {
        "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/a86dba58-0043-4cc6-a1bb-69d5e86f3ca3"
    }
}
`

// CreatedTransferRequest is the expected created zone transfer request.
var CreatedTransferRequest = FirstTransferRequest

// HandleTransferRequestCreationSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T) {
	createURL := "/zones/a6a8515c-5d80-48c0-955b-fde631b59791/tasks/transfer_requests"
	th.Mux.HandleFunc(createURL,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, CreateTransferRequest)

			w.WriteHeader(http.StatusCreated)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, CreateTransferRequestResponse)
		})
}

// UpdateTransferRequest is a sample request to update a zone transfer request.
const UpdateTransferRequest = `
{
    "description": "Updated Description"
}
`

// UpdatedTransferRequestResponse is a sample response to update a zone transfer request.
const UpdatedTransferRequestResponse = `
{
    "id": "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3",
    "zone_id": "a6a8515c-5d80-48c0-955b-fde631b59791",
    "project_id": "4335d1f0-f793-11e2-b778-0800200c9a66",
    "target_project_id": "05d98711-b3a1-4264-a395-f46383671ee6",
    "description": "Updated Description",
    "key": "KJSDH23Z",
    "status": "ACTIVE",
    "zone_name": "example1.org.",
    "created_at": "2020-10-12T08:38:58.000000",
    "links": {
        "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_requests/a86dba58-0043-4cc6-a1bb-69d5e86f3ca3"
    }
}
`

// HandleTransferRequestUpdateSuccessfully configures the test server to respond to an Update request.
func HandleUpdateSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_requests"
	th.Mux.HandleFunc(s.Join([]string{baseURL, FirstTransferRequest.ID}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, UpdateTransferRequest)

			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, UpdatedTransferRequestResponse)
		})
}

// HandleTransferRequestDeleteSuccessfully configures the test server to respond to an Delete request.
func HandleDeleteSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_requests"
	th.Mux.HandleFunc(s.Join([]string{baseURL, FirstTransferRequest.ID}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}
