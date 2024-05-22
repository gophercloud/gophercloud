package testing

import (
	"fmt"
	"net/http"
	s "strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	transferAccepts "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/accept"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput is a sample response to a List call.
const ListOutput = `
{
    "transfer_accepts": [
        {
            "id": "92236f39-0fad-4f8f-bf25-fbdf027de34d",
            "status": "COMPLETE",
            "project_id": "9f3cfb08bf52469abe598e127676cd57",
            "zone_id": "cd046f4b-f4dc-4e41-b946-1a2d32e1be40",
            "key": "M2KA0Y20",
            "zone_transfer_request_id": "fc46bb1f-bdf0-4e67-96e0-f8c04f26261c",
            "created_at": "2020-10-12T08:38:58.000000",
            "links": {
                "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/92236f39-0fad-4f8f-bf25-fbdf027de34d",
                "zone": "https://127.0.0.1:9001/v2/zones/cd046f4b-f4dc-4e41-b946-1a2d32e1be40"
            }
        },
        {
            "id": "f785ef12-7ee0-4c30-bd67-a2b9edba0dff",
            "status": "ACTIVE",
            "project_id": "9f3cfb08bf52469abe598e127676cd57",
            "zone_id": "30d67a9a-d6df-4ba7-9b55-fb49e7987f84",
            "key": "SDF32HJ1",
            "zone_transfer_request_id": "c5d11193-72ea-4d9f-ba04-7f80e99627fa",
            "created_at": "2020-10-12T09:38:58.000000",
            "updated_at": "2020-10-12T09:38:58.000000",
            "links": {
                "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/f785ef12-7ee0-4c30-bd67-a2b9edba0dff",
                "zone": "https://127.0.0.1:9001/v2/zones/30d67a9a-d6df-4ba7-9b55-fb49e7987f84"
            }
        }
    ]
}
`

// FilteredListOutput is a sample response to a List call with Opts.
const FilteredListOutput = `
{
    "transfer_accepts": [
        {
            "id": "f785ef12-7ee0-4c30-bd67-a2b9edba0dff",
            "status": "ACTIVE",
            "project_id": "9f3cfb08bf52469abe598e127676cd57",
            "zone_id": "30d67a9a-d6df-4ba7-9b55-fb49e7987f84",
            "key": "SDF32HJ1",
            "zone_transfer_request_id": "c5d11193-72ea-4d9f-ba04-7f80e99627fa",
            "created_at": "2020-10-12T09:38:58.000000",
            "updated_at": "2020-10-12T09:38:58.000000",
            "links": {
                "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/f785ef12-7ee0-4c30-bd67-a2b9edba0dff",
                "zone": "https://127.0.0.1:9001/v2/zones/30d67a9a-d6df-4ba7-9b55-fb49e7987f84"
            }
        }
    ]
}
`

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
    "id": "92236f39-0fad-4f8f-bf25-fbdf027de34d",
    "status": "COMPLETE",
    "project_id": "9f3cfb08bf52469abe598e127676cd57",
    "zone_id": "cd046f4b-f4dc-4e41-b946-1a2d32e1be40",
    "key": "M2KA0Y20",
    "zone_transfer_request_id": "fc46bb1f-bdf0-4e67-96e0-f8c04f26261c",
    "created_at": "2020-10-12T08:38:58.000000",
    "links": {
        "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/92236f39-0fad-4f8f-bf25-fbdf027de34d",
        "zone": "https://127.0.0.1:9001/v2/zones/cd046f4b-f4dc-4e41-b946-1a2d32e1be40"
    }
}
`

// FirstTransferAccept is the first result in ListOutput
var FirstTransferAcceptCreatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-10-12T08:38:58.000000")
var FirstTransferAccept = transferAccepts.TransferAccept{
	ID:                    "92236f39-0fad-4f8f-bf25-fbdf027de34d",
	ZoneID:                "cd046f4b-f4dc-4e41-b946-1a2d32e1be40",
	ProjectID:             "9f3cfb08bf52469abe598e127676cd57",
	ZoneTransferRequestID: "fc46bb1f-bdf0-4e67-96e0-f8c04f26261c",
	Key:                   "M2KA0Y20",
	Status:                "COMPLETE",
	CreatedAt:             FirstTransferAcceptCreatedAt,
	Links: map[string]any{
		"self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/92236f39-0fad-4f8f-bf25-fbdf027de34d",
		"zone": "https://127.0.0.1:9001/v2/zones/cd046f4b-f4dc-4e41-b946-1a2d32e1be40",
	},
}

// SecondTransferRequest is the second result in ListOutput
var SecondTransferAcceptCreatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-10-12T09:38:58.000000")
var SecondTransferAcceptUpdatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-10-12T09:38:58.000000")
var SecondTransferAccept = transferAccepts.TransferAccept{
	ID:                    "f785ef12-7ee0-4c30-bd67-a2b9edba0dff",
	Status:                "ACTIVE",
	ProjectID:             "9f3cfb08bf52469abe598e127676cd57",
	ZoneID:                "30d67a9a-d6df-4ba7-9b55-fb49e7987f84",
	ZoneTransferRequestID: "c5d11193-72ea-4d9f-ba04-7f80e99627fa",
	Key:                   "SDF32HJ1",
	CreatedAt:             SecondTransferAcceptCreatedAt,
	UpdatedAt:             SecondTransferAcceptUpdatedAt,
	Links: map[string]any{
		"self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/f785ef12-7ee0-4c30-bd67-a2b9edba0dff",
		"zone": "https://127.0.0.1:9001/v2/zones/30d67a9a-d6df-4ba7-9b55-fb49e7987f84",
	},
}

// ExpectedTransferAcceptSlice is the slice of results that should be parsed
// from ListOutput, in the expected order.
var ExpectedTransferAcceptSlice = []transferAccepts.TransferAccept{FirstTransferAccept, SecondTransferAccept}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_accepts"
	th.Mux.HandleFunc(baseURL,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, ListOutput)
		})
}

// HandleFilteredListSuccessfully configures the test server to respond to a List request with Opts.
func HandleFilteredListSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_accepts"
	th.Mux.HandleFunc(baseURL,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, FilteredListOutput)
		})
}

// HandleGetSuccessfully configures the test server to respond to a List request.
func HandleGetSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_accepts"
	th.Mux.HandleFunc(s.Join([]string{baseURL, FirstTransferAccept.ID}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, GetOutput)
		})
}

// CreateTransferAccept is a sample request to create a zone.
const CreateTransferAccept = `
{
    "key": "M2KA0Y20",
    "zone_transfer_request_id": "fc46bb1f-bdf0-4e67-96e0-f8c04f26261c"
}
`

// CreateTransferAcceptResponse is a sample response to a create request.
const CreateTransferAcceptResponse = `
{
    "id": "92236f39-0fad-4f8f-bf25-fbdf027de34d",
    "zone_transfer_request_id": "fc46bb1f-bdf0-4e67-96e0-f8c04f26261c",
    "project_id": "9f3cfb08bf52469abe598e127676cd57",
    "key": "M2KA0Y20",
    "status": "COMPLETE",
    "zone_id": "cd046f4b-f4dc-4e41-b946-1a2d32e1be40",
    "created_at": "2020-10-12T08:38:58.000000",
    "links": {
        "self": "https://127.0.0.1:9001/v2/zones/tasks/transfer_accepts/92236f39-0fad-4f8f-bf25-fbdf027de34d",
        "zone": "https://127.0.0.1:9001/v2/zones/cd046f4b-f4dc-4e41-b946-1a2d32e1be40"
    }
}
`

// CreatedTransferRequest is the expected created zone transfer request.
var CreatedTransferAccept = FirstTransferAccept

// HandleTransferRequestCreationSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T) {
	baseURL := "/zones/tasks/transfer_accepts"
	th.Mux.HandleFunc(baseURL,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, CreateTransferAccept)

			w.WriteHeader(http.StatusCreated)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, CreateTransferAcceptResponse)
		})
}
