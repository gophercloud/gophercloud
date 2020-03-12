package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/backups"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const ListResponse = `
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
      "href": "%s/backups?marker=1",
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
    "backup_url": "eyJpZCI6ImQzMjAxOWQzLWJjNmUtNDMxOS05YzFkLTY3MjJmYzEzNmEyMiIsInZvbHVtZV9pZCI6ImNmOWJjNmZhLWM1YmMtNDFmNi1iYzRlLTZlNzZjMGJlYTk1OSIsInNuYXBzaG90X2lkIjoiIiwic3RhdHVzIjoiYXZhaWxhYmxlIiwic2l6ZSI6MSwib2JqZWN0X2NvdW50IjoyLCJjb250YWluZXIiOiJteS10ZXN0LWJhY2t1cCIsInNlcnZpY2VfbWV0YWRhdGEiOiJ2b2x1bWVfY2Y5YmM2ZmEtYzViYy00MWY2LWJjNGUtNmU3NmMwYmVhOTU5LzIwMjAwMzExMTkyODU1L2F6X3JlZ2lvbmJfYmFja3VwX2I4N2JiMWU1LTBkNGUtNDQ1ZS1hNTQ4LTVhZTc0MjU2MmJhYyIsInNlcnZpY2UiOiJjaW5kZXIuYmFja3VwLmRyaXZlcnMuc3dpZnQuU3dpZnRCYWNrdXBEcml2ZXIiLCJob3N0IjoiY2luZGVyLWJhY2t1cC1ob3N0MSIsInVzZXJfaWQiOiI5MzUxNGUwNC1hMDI2LTRmNjAtODE3Ni0zOTVjODU5NTAxZGQiLCJ0ZW1wX3NuYXBzaG90X2lkIjoiIiwidGVtcF92b2x1bWVfaWQiOiIiLCJyZXN0b3JlX3ZvbHVtZV9pZCI6IiIsIm51bV9kZXBlbmRlbnRfYmFja3VwcyI6MCwiZW5jcnlwdGlvbl9rZXlfaWQiOiIiLCJwYXJlbnRfaWQiOiIiLCJkZWxldGVkIjpmYWxzZSwiZGlzcGxheV9uYW1lIjoiIiwiZGlzcGxheV9kZXNjcmlwdGlvbiI6IiIsImRyaXZlcl9pbmZvIjpudWxsLCJmYWlsX3JlYXNvbiI6IiIsInByb2plY3RfaWQiOiIxNGYxYzFmNWQxMmI0NzU1Yjk0ZWRlZjc4ZmY4YjMyNSIsIm1ldGFkYXRhIjpudWxsLCJhdmFpbGFiaWxpdHlfem9uZSI6InJlZ2lvbjFiIiwiY3JlYXRlZF9hdCI6IjIwMjAtMDMtMTFUMTk6MjU6MjRaIiwidXBkYXRlZF9hdCI6IjIwMjAtMDMtMTFUMTk6Mjk6MDhaIiwiZGVsZXRlZF9hdCI6bnVsbCwiZGF0YV90aW1lc3RhbXAiOiIyMDIwLTAzLTExVDE5OjI1OjI0WiJ9"
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

var backupImport = backups.ImportBackup{
	ID:               "d32019d3-bc6e-4319-9c1d-6722fc136a22",
	Status:           "available",
	AvailabilityZone: &availabilityZone,
	VolumeID:         "cf9bc6fa-c5bc-41f6-bc4e-6e76c0bea959",
	UpdatedAt:        time.Date(2020, 3, 11, 19, 29, 8, 0, time.UTC),
	Host:             "cinder-backup-host1",
	UserID:           "93514e04-a026-4f60-8176-395c859501dd",
	ServiceMetadata:  "volume_cf9bc6fa-c5bc-41f6-bc4e-6e76c0bea959/20200311192855/az_regionb_backup_b87bb1e5-0d4e-445e-a548-5ae742562bac",
	Size:             1,
	ObjectCount:      2,
	Container:        "my-test-backup",
	Service:          "cinder.backup.drivers.swift.SwiftBackupDriver",
	CreatedAt:        time.Date(2020, 3, 11, 19, 25, 24, 0, time.UTC),
	DataTimestamp:    time.Date(2020, 3, 11, 19, 25, 24, 0, time.UTC),
	ProjectID:        "14f1c1f5d12b4755b94edef78ff8b325",
}
var availabilityZone = "region1b"
var backupURL, _ = json.Marshal(backupImport)

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		r.ParseForm()
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
