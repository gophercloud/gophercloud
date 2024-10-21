package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const (
	shareEndpoint = "/share-replicas"
	replicaID     = "3b9c33e8-b136-45c6-84a6-019c8db1d550"
)

var createRequest = `{
  "share_replica": {
    "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
    "availability_zone": "zone-1"
  }
}
`

var createResponse = `{
  "share_replica": {
    "id": "3b9c33e8-b136-45c6-84a6-019c8db1d550",
    "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
    "availability_zone": "zone-1",
    "created_at": "2023-05-26T12:32:56.391337",
    "status": "creating",
    "share_network_id": "ca0163c8-3941-4420-8b01-41517e19e366",
    "share_server_id": null,
    "replica_state": null,
    "updated_at": null
  }
}
`

// MockCreateResponse creates a mock response
func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		th.TestJSONRequest(t, r, createRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, createResponse)
	})
}

// MockDeleteResponse creates a mock delete response
func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		w.WriteHeader(http.StatusAccepted)
	})
}

var promoteRequest = `{
  "promote": {
    "quiesce_wait_time": 30
  }
}
`

func MockPromoteResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		th.TestJSONRequest(t, r, promoteRequest)
		w.WriteHeader(http.StatusAccepted)
	})
}

var resyncRequest = `{
  "resync": null
}
`

func MockResyncResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		th.TestJSONRequest(t, r, resyncRequest)
		w.WriteHeader(http.StatusAccepted)
	})
}

var resetStatusRequest = `{
  "reset_status": {
    "status": "available"
  }
}
`

func MockResetStatusResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		th.TestJSONRequest(t, r, resetStatusRequest)
		w.WriteHeader(http.StatusAccepted)
	})
}

var resetStateRequest = `{
  "reset_replica_state": {
    "replica_state": "active"
  }
}
`

func MockResetStateResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		th.TestJSONRequest(t, r, resetStateRequest)
		w.WriteHeader(http.StatusAccepted)
	})
}

var deleteRequest = `{
  "force_delete": null
}
`

func MockForceDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		th.TestJSONRequest(t, r, deleteRequest)
		w.WriteHeader(http.StatusAccepted)
	})
}

var getResponse = `{
  "share_replica": {
    "id": "3b9c33e8-b136-45c6-84a6-019c8db1d550",
    "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
    "availability_zone": "zone-1",
    "created_at": "2023-05-26T12:32:56.391337",
    "status": "available",
    "share_network_id": "ca0163c8-3941-4420-8b01-41517e19e366",
    "share_server_id": "5ccc1b0c-334a-4e46-81e6-b52e03223060",
    "replica_state": "active",
    "updated_at": "2023-05-26T12:33:28.265716"
  }
}
`

// MockGetResponse creates a mock get response
func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getResponse)
	})
}

var listResponse = `{
  "share_replicas": [
    {
      "id": "3b9c33e8-b136-45c6-84a6-019c8db1d550",
      "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
      "status": "available",
      "replica_state": "active"
    },
    {
      "id": "4b70c2e2-eec7-4699-880d-4da9051ca162",
      "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
      "status": "available",
      "replica_state": "out_of_sync"
    },
    {
      "id": "920bb037-bdd7-48a1-98f0-1aa1787ca3eb",
      "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
      "status": "available",
      "replica_state": "in_sync"
    }
  ]
}
`

var listEmptyResponse = `{"share_replicas": []}`

// MockListResponse creates a mock detailed-list response
func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("offset")
		shareID := r.Form.Get("share_id")
		if shareID != "65a34695-f9e5-4eea-b48d-a0b261d82943" {
			th.AssertNoErr(t, fmt.Errorf("unexpected share_id"))
		}

		switch marker {
		case "":
			fmt.Fprint(w, listResponse)
		default:
			fmt.Fprint(w, listEmptyResponse)
		}
	})
}

var listDetailResponse = `{
  "share_replicas": [
    {
      "id": "3b9c33e8-b136-45c6-84a6-019c8db1d550",
      "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
      "availability_zone": "zone-1",
      "created_at": "2023-05-26T12:32:56.391337",
      "status": "available",
      "share_network_id": "ca0163c8-3941-4420-8b01-41517e19e366",
      "share_server_id": "5ccc1b0c-334a-4e46-81e6-b52e03223060",
      "replica_state": "active",
      "updated_at": "2023-05-26T12:33:28.265716"
    },
    {
      "id": "4b70c2e2-eec7-4699-880d-4da9051ca162",
      "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
      "availability_zone": "zone-2",
      "created_at": "2023-05-26T11:59:38.313089",
      "status": "available",
      "share_network_id": "ca0163c8-3941-4420-8b01-41517e19e366",
      "share_server_id": "81aa586e-3a03-4f92-98bd-807d87a61c1a",
      "replica_state": "out_of_sync",
      "updated_at": "2023-05-26T12:00:04.321081"
    },
    {
      "id": "920bb037-bdd7-48a1-98f0-1aa1787ca3eb",
      "share_id": "65a34695-f9e5-4eea-b48d-a0b261d82943",
      "availability_zone": "zone-1",
      "created_at": "2023-05-26T12:32:45.751834",
      "status": "available",
      "share_network_id": "ca0163c8-3941-4420-8b01-41517e19e366",
      "share_server_id": "b87ea601-7d4c-47f3-8956-6876b7a6b6db",
      "replica_state": "in_sync",
      "updated_at": "2023-05-26T12:36:04.110328"
    }
  ]
}
`

var listDetailEmptyResponse = `{"share_replicas": []}`

// MockListDetailResponse creates a mock detailed-list response
func MockListDetailResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.11")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("offset")
		shareID := r.Form.Get("share_id")
		if shareID != "65a34695-f9e5-4eea-b48d-a0b261d82943" {
			th.AssertNoErr(t, fmt.Errorf("unexpected share_id"))
		}

		switch marker {
		case "":
			fmt.Fprint(w, listDetailResponse)
		default:
			fmt.Fprint(w, listDetailEmptyResponse)
		}
	})
}

var listExportLocationsResponse = `{
  "export_locations": [
    {
      "id": "3fc02d3c-da47-42a2-88b8-2d48f8c276bd",
      "path": "192.168.1.123:/var/lib/manila/mnt/share-3b9c33e8-b136-45c6-84a6-019c8db1d550",
      "preferred": true,
      "replica_state": "active",
      "availability_zone": "zone-1"
    },
    {
      "id": "ae73e762-e8b9-4aad-aad3-23afb7cd6825",
      "path": "192.168.1.124:/var/lib/manila/mnt/share-3b9c33e8-b136-45c6-84a6-019c8db1d550",
      "preferred": false,
      "replica_state": "active",
      "availability_zone": "zone-1"
    }
  ]
}
`

// MockListExportLocationsResponse creates a mock get export locations response
func MockListExportLocationsResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/export-locations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.47")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, listExportLocationsResponse)
	})
}

var getExportLocationResponse = `{
  "export_location": {
    "id": "ae73e762-e8b9-4aad-aad3-23afb7cd6825",
    "path": "192.168.1.124:/var/lib/manila/mnt/share-3b9c33e8-b136-45c6-84a6-019c8db1d550",
    "preferred": false,
    "created_at": "2023-05-26T12:44:33.987960",
    "updated_at": "2023-05-26T12:44:33.958363",
    "replica_state": "active",
    "availability_zone": "zone-1"
  }
}
`

// MockGetExportLocationResponse creates a mock get export location response
func MockGetExportLocationResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+replicaID+"/export-locations/ae73e762-e8b9-4aad-aad3-23afb7cd6825", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		th.TestHeader(t, r, "X-OpenStack-Manila-API-Version", "2.47")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getExportLocationResponse)
	})
}
