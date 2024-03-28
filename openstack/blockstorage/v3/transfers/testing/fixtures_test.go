package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/transfers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const ListOutput = `
{
  "transfers": [
    {
      "created_at": "2020-02-28T12:44:28.051989",
      "volume_id": "2f6f1684-1ded-40db-8a49-7c87dedbc758",
      "id": "b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
      "links": [
        {
          "href": "https://volume/v3/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
          "rel": "self"
        },
        {
          "href": "https://volume/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
          "rel": "bookmark"
        }
      ],
      "name": null
    }
  ]
}
`

const GetOutput = `
{
  "transfer": {
    "created_at": "2020-02-28T12:44:28.051989",
    "volume_id": "2f6f1684-1ded-40db-8a49-7c87dedbc758",
    "id": "b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
    "links": [
      {
        "href": "https://volume/v3/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
        "rel": "self"
      },
      {
        "href": "https://volume/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
        "rel": "bookmark"
      }
    ],
    "name": null
  }
}
`

const CreateRequest = `
{
  "transfer": {
    "volume_id": "2f6f1684-1ded-40db-8a49-7c87dedbc758"
  }
}
`

const CreateResponse = `
{
  "transfer": {
    "auth_key": "cb67e0e7387d9eac",
    "created_at": "2020-02-28T12:44:28.051989",
    "id": "b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
    "links": [
      {
        "href": "https://volume/v3/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
        "rel": "self"
      },
      {
        "href": "https://volume/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
        "rel": "bookmark"
      }
    ],
    "name": null,
    "volume_id": "2f6f1684-1ded-40db-8a49-7c87dedbc758"
  }
}
`

const AcceptTransferRequest = `
{
    "accept": {
        "auth_key": "9266c59563c84664"
    }
}
`

const AcceptTransferResponse = `
{
  "transfer": {
    "id": "b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
    "links": [
      {
        "href": "https://volume/v3/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
        "rel": "self"
      },
      {
        "href": "https://volume/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
        "rel": "bookmark"
      }
    ],
    "name": null,
    "volume_id": "2f6f1684-1ded-40db-8a49-7c87dedbc758"
  }
}
`

var TransferRequest = transfers.CreateOpts{
	VolumeID: "2f6f1684-1ded-40db-8a49-7c87dedbc758",
}

var createdAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2020-02-28T12:44:28.051989")
var TransferResponse = transfers.Transfer{
	ID:        "b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
	AuthKey:   "cb67e0e7387d9eac",
	Name:      "",
	VolumeID:  "2f6f1684-1ded-40db-8a49-7c87dedbc758",
	CreatedAt: createdAt,
	Links: []map[string]string{
		{
			"href": "https://volume/v3/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
			"rel":  "self",
		},
		{
			"href": "https://volume/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
			"rel":  "bookmark",
		},
	},
}

var TransferListResponse = []transfers.Transfer{TransferResponse}

var AcceptRequest = transfers.AcceptOpts{
	AuthKey: "9266c59563c84664",
}

var AcceptResponse = transfers.Transfer{
	ID:       "b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
	Name:     "",
	VolumeID: "2f6f1684-1ded-40db-8a49-7c87dedbc758",
	Links: []map[string]string{
		{
			"href": "https://volume/v3/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
			"rel":  "self",
		},
		{
			"href": "https://volume/53c2b94f63fb4f43a21b92d119ce549f/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f",
			"rel":  "bookmark",
		},
	},
}

func HandleCreateTransfer(t *testing.T) {
	th.Mux.HandleFunc("/os-volume-transfer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, CreateResponse)
	})
}

func HandleAcceptTransfer(t *testing.T) {
	th.Mux.HandleFunc("/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f/accept", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		th.TestJSONRequest(t, r, AcceptTransferRequest)

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, AcceptTransferResponse)
	})
}

func HandleDeleteTransfer(t *testing.T) {
	th.Mux.HandleFunc("/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleListTransfers(t *testing.T) {
	th.Mux.HandleFunc("/os-volume-transfer/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		th.TestFormValues(t, r, map[string]string{"all_tenants": "true"})

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

func HandleGetTransfer(t *testing.T) {
	th.Mux.HandleFunc("/os-volume-transfer/b8913bfd-a4d3-4ec5-bd8b-fe2dbeef9f4f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}
