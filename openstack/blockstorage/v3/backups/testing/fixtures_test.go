package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/backups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const ListResponse = `
{
  "backups": [
    {
      "id": "289da7f8-6440-407c-9fb4-7db01ec49164",
      "name": "backup-001"
    },
    {
      "id": "96c3bda7-c82a-4f50-be73-ca7621794835",
      "name": "backup-002"
    }
  ],
  "backups_links": [
    {
      "href": "%s/backups?marker=1",
      "rel": "next"
    }
  ]
}
`

const ListDetailResponse = `
{
  "backups": [
    {
      "id": "289da7f8-6440-407c-9fb4-7db01ec49164",
      "name": "backup-001",
      "volume_id": "521752a6-acf6-4b2d-bc7a-119f9148cd8c",
      "description": "Daily Backup",
      "status": "available",
      "size": 30,
      "created_at": "2017-05-30T03:35:03.000000"
    },
    {
      "id": "96c3bda7-c82a-4f50-be73-ca7621794835",
      "name": "backup-002",
      "volume_id": "76b8950a-8594-4e5b-8dce-0dfa9c696358",
      "description": "Weekly Backup",
      "status": "available",
      "size": 25,
      "created_at": "2017-05-30T03:35:03.000000"
    }
  ],
  "backups_links": [
    {
      "href": "%s/backups/detail?marker=1",
      "rel": "next"
    }
  ]
}
`

const GetResponse = `
{
  "backup": {
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "name": "backup-001",
    "description": "Daily backup",
    "volume_id": "521752a6-acf6-4b2d-bc7a-119f9148cd8c",
    "status": "available",
    "size": 30,
    "created_at": "2017-05-30T03:35:03.000000"
  }
}
`
const CreateRequest = `
{
  "backup": {
    "volume_id": "1234",
    "name": "backup-001"
  }
}
`

const CreateResponse = `
{
  "backup": {
    "volume_id": "1234",
    "name": "backup-001",
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "description": "Daily backup",
    "volume_id": "1234",
    "status": "available",
    "size": 30,
    "created_at": "2017-05-30T03:35:03.000000"
  }
}
`

const RestoreRequest = `
{
  "restore": {
    "name": "vol-001",
    "volume_id": "1234"
  }
}
`

const RestoreResponse = `
{
  "restore": {
    "backup_id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "volume_id": "1234",
    "volume_name": "vol-001"
  }
}
`

const ExportResponse = `
{
  "backup-record": {
    "backup_service": "cinder.backup.drivers.swift.SwiftBackupDriver",
    "backup_url": "eyJpZCI6ImQzMjAxOWQzLWJjNmUtNDMxOS05YzFkLTY3MjJmYzEzNmEyMiIsInZvbHVtZV9pZCI6ImNmOWJjNmZhLWM1YmMtNDFmNi1iYzRlLTZlNzZjMGJlYTk1OSIsInNuYXBzaG90X2lkIjpudWxsLCJzdGF0dXMiOiJhdmFpbGFibGUiLCJzaXplIjoxLCJvYmplY3RfY291bnQiOjIsImNvbnRhaW5lciI6Im15LXRlc3QtYmFja3VwIiwic2VydmljZV9tZXRhZGF0YSI6InZvbHVtZV9jZjliYzZmYS1jNWJjLTQxZjYtYmM0ZS02ZTc2YzBiZWE5NTkvMjAyMDAzMTExOTI4NTUvYXpfcmVnaW9uYl9iYWNrdXBfYjg3YmIxZTUtMGQ0ZS00NDVlLWE1NDgtNWFlNzQyNTYyYmFjIiwic2VydmljZSI6ImNpbmRlci5iYWNrdXAuZHJpdmVycy5zd2lmdC5Td2lmdEJhY2t1cERyaXZlciIsImhvc3QiOiJjaW5kZXItYmFja3VwLWhvc3QxIiwidXNlcl9pZCI6IjkzNTE0ZTA0LWEwMjYtNGY2MC04MTc2LTM5NWM4NTk1MDFkZCIsInRlbXBfc25hcHNob3RfaWQiOm51bGwsInRlbXBfdm9sdW1lX2lkIjpudWxsLCJyZXN0b3JlX3ZvbHVtZV9pZCI6bnVsbCwibnVtX2RlcGVuZGVudF9iYWNrdXBzIjpudWxsLCJlbmNyeXB0aW9uX2tleV9pZCI6bnVsbCwicGFyZW50X2lkIjpudWxsLCJkZWxldGVkIjpmYWxzZSwiZGlzcGxheV9uYW1lIjpudWxsLCJkaXNwbGF5X2Rlc2NyaXB0aW9uIjpudWxsLCJkcml2ZXJfaW5mbyI6bnVsbCwiZmFpbF9yZWFzb24iOm51bGwsInByb2plY3RfaWQiOiIxNGYxYzFmNWQxMmI0NzU1Yjk0ZWRlZjc4ZmY4YjMyNSIsIm1ldGFkYXRhIjp7fSwiYXZhaWxhYmlsaXR5X3pvbmUiOiJyZWdpb24xYiIsImNyZWF0ZWRfYXQiOiIyMDIwLTAzLTExVDE5OjI1OjI0WiIsInVwZGF0ZWRfYXQiOiIyMDIwLTAzLTExVDE5OjI5OjA4WiIsImRlbGV0ZWRfYXQiOm51bGwsImRhdGFfdGltZXN0YW1wIjoiMjAyMC0wMy0xMVQxOToyNToyNFoifQ=="
  }
}
`

const ImportRequest = ExportResponse

const ImportResponse = `
{
  "backup": {
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "links": [
      {
        "href": "https://volume/v2/14f1c1f5d12b4755b94edef78ff8b325/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "rel": "self"
      },
      {
        "href": "https://volume/14f1c1f5d12b4755b94edef78ff8b325/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "rel": "bookmark"
      }
    ],
    "name": null
  }
}
`

const ResetRequest = `
{
  "os-reset_status": {
    "status": "error"
  }
}
`

const ForceDeleteRequest = `
{
  "os-force_delete": {}
}
`

var (
	status           = "available"
	availabilityZone = "region1b"
	host             = "cinder-backup-host1"
	serviceMetadata  = "volume_cf9bc6fa-c5bc-41f6-bc4e-6e76c0bea959/20200311192855/az_regionb_backup_b87bb1e5-0d4e-445e-a548-5ae742562bac"
	size             = 1
	objectCount      = 2
	container        = "my-test-backup"
	service          = "cinder.backup.drivers.swift.SwiftBackupDriver"
	backupImport     = backups.ImportBackup{
		ID:               "d32019d3-bc6e-4319-9c1d-6722fc136a22",
		Status:           &status,
		AvailabilityZone: &availabilityZone,
		VolumeID:         "cf9bc6fa-c5bc-41f6-bc4e-6e76c0bea959",
		UpdatedAt:        time.Date(2020, 3, 11, 19, 29, 8, 0, time.UTC),
		Host:             &host,
		UserID:           "93514e04-a026-4f60-8176-395c859501dd",
		ServiceMetadata:  &serviceMetadata,
		Size:             &size,
		ObjectCount:      &objectCount,
		Container:        &container,
		Service:          &service,
		CreatedAt:        time.Date(2020, 3, 11, 19, 25, 24, 0, time.UTC),
		DataTimestamp:    time.Date(2020, 3, 11, 19, 25, 24, 0, time.UTC),
		ProjectID:        "14f1c1f5d12b4755b94edef78ff8b325",
		Metadata:         make(map[string]string),
	}
	backupURL, _ = json.Marshal(backupImport)
)

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, ListResponse, th.Server.URL)
		case "1":
			fmt.Fprintf(w, `{"backups": []}`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockListDetailResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, ListDetailResponse, th.Server.URL)
		case "1":
			fmt.Fprintf(w, `{"backups": []}`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetResponse)
	})
}

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, CreateResponse)
	})
}

func MockRestoreResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22/restore", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RestoreRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, RestoreResponse)
	})
}

func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

func MockExportResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22/export_record", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ExportResponse)
	})
}

func MockImportResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/import_record", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ImportRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, ImportResponse)
	})
}

// MockResetStatusResponse provides mock response for reset backup status API call
func MockResetStatusResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, ResetRequest)

		w.WriteHeader(http.StatusAccepted)
	})
}

// MockForceDeleteResponse provides mock response for force delete backup API call
func MockForceDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups/d32019d3-bc6e-4319-9c1d-6722fc136a22/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, ForceDeleteRequest)

		w.WriteHeader(http.StatusAccepted)
	})
}
