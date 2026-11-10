package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/runbooks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const runbookID = "fc6b1709-8dd5-86b0-2d34-5203d1c29127"

// RunbookListBody is the first page of a runbook list, with a next link.
const RunbookListBody = `
{
  "runbooks": [
    {
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/runbooks/fc6b1709-8dd5-86b0-2d34-5203d1c29127",
          "rel": "self"
        }
      ],
      "name": "CUSTOM_AWESOME",
      "uuid": "fc6b1709-8dd5-86b0-2d34-5203d1c29127"
    }
  ],
  "next": "%s/runbooks?limit=1&marker=fc6b1709-8dd5-86b0-2d34-5203d1c29127"
}
`

// RunbookListSecondPageBody is the last page of a runbook list.
const RunbookListSecondPageBody = `
{
  "runbooks": [
    {
      "name": "CUSTOM_OTHER",
      "uuid": "8bbbbbc6-7c37-4a32-9c5d-6b1f5e7c2a11"
    }
  ]
}
`

// SingleRunbookBody is a runbook returned with every field at microversion 1.112.
const SingleRunbookBody = `
{
  "created_at": "2024-08-18T22:28:48+00:00",
  "description": "Enable logical processors",
  "disable_ramdisk": false,
  "extra": {"team": "metal"},
  "links": [
    {
      "href": "http://ironic.example.com:6385/v1/runbooks/fc6b1709-8dd5-86b0-2d34-5203d1c29127",
      "rel": "self"
    }
  ],
  "name": "CUSTOM_AWESOME",
  "owner": "f5a7c3a8b9e94d7f8a1b2c3d4e5f6a7b",
  "public": false,
  "steps": [
    {
      "args": {
        "settings": [
          {
            "name": "LogicalProc",
            "value": "Enabled"
          }
        ]
      },
      "interface": "bios",
      "order": 1,
      "step": "apply_configuration"
    }
  ],
  "traits": ["CUSTOM_GPU"],
  "updated_at": "2024-08-18T22:28:49+00:00",
  "uuid": "fc6b1709-8dd5-86b0-2d34-5203d1c29127"
}
`

// RunbookTraitsBody is the response to a runbook traits request.
const RunbookTraitsBody = `
{
  "traits": ["CUSTOM_GPU", "CUSTOM_FPGA"]
}
`

var (
	createdAt, _ = time.Parse(time.RFC3339, "2024-08-18T22:28:48+00:00")
	updatedAt, _ = time.Parse(time.RFC3339, "2024-08-18T22:28:49+00:00")
	description  = "Enable logical processors"

	RunbookAwesome = runbooks.Runbook{
		UUID:           runbookID,
		Name:           "CUSTOM_AWESOME",
		DisableRamdisk: false,
		Extra:          map[string]any{"team": "metal"},
		Public:         false,
		Owner:          "f5a7c3a8b9e94d7f8a1b2c3d4e5f6a7b",
		Description:    &description,
		Traits:         []string{"CUSTOM_GPU"},
		Steps: []runbooks.Step{
			{
				Interface: nodes.InterfaceBIOS,
				Step:      "apply_configuration",
				Args: map[string]any{
					"settings": []any{
						map[string]any{"name": "LogicalProc", "value": "Enabled"},
					},
				},
				Order: 1,
			},
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Links: []any{
			map[string]any{
				"href": "http://ironic.example.com:6385/v1/runbooks/fc6b1709-8dd5-86b0-2d34-5203d1c29127",
				"rel":  "self",
			},
		},
	}
)

// HandleRunbookListSuccessfully sets up the test server to respond to a
// paginated runbook List request.
func HandleRunbookListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}

		w.Header().Add("Content-Type", "application/json")
		switch r.Form.Get("marker") {
		case "":
			th.TestFormValues(t, r, map[string]string{"limit": "1"})
			fmt.Fprintf(w, RunbookListBody, fakeServer.Server.URL)
		case runbookID:
			th.TestFormValues(t, r, map[string]string{"limit": "1", "marker": runbookID})
			fmt.Fprint(w, RunbookListSecondPageBody)
		default:
			t.Fatalf("/runbooks invoked with unexpected marker=[%s]", r.Form.Get("marker"))
		}
	})
}

// HandleRunbookListDetailSuccessfully sets up the test server to respond to
// a runbook List request with every list option set.
func HandleRunbookListDetailSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestFormValues(t, r, map[string]string{
			"detail":   "true",
			"fields":   "uuid,name",
			"limit":    "10",
			"marker":   runbookID,
			"sort_dir": "desc",
			"sort_key": "name",
		})

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"runbooks": [%s]}`, SingleRunbookBody)
	})
}

// HandleRunbookCreationSuccessfully sets up the test server to respond to a
// runbook Create request with every option set.
func HandleRunbookCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"name": "CUSTOM_AWESOME",
			"uuid": "fc6b1709-8dd5-86b0-2d34-5203d1c29127",
			"description": "Enable logical processors",
			"disable_ramdisk": false,
			"extra": {"team": "metal"},
			"owner": "f5a7c3a8b9e94d7f8a1b2c3d4e5f6a7b",
			"public": false,
			"steps": [
				{
					"interface": "bios",
					"step": "apply_configuration",
					"args": {"settings": [{"name": "LogicalProc", "value": "Enabled"}]},
					"order": 0
				}
			]
		}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, SingleRunbookBody)
	})
}

// HandleRunbookGetSuccessfully sets up the test server to respond to a
// runbook Get request.
func HandleRunbookGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks/"+runbookID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, SingleRunbookBody)
	})
}

// HandleRunbookUpdateSuccessfully sets up the test server to respond to a
// runbook Update request.
func HandleRunbookUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks/"+runbookID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `[
			{"op": "replace", "path": "/description", "value": "Enable logical processors"},
			{"op": "remove", "path": "/extra/team", "value": null}
		]`)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, SingleRunbookBody)
	})
}

// HandleRunbookDeletionSuccessfully sets up the test server to respond to a
// runbook Delete request.
func HandleRunbookDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks/"+runbookID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleRunbookTraitsSuccessfully sets up the test server to respond to
// listing, replacing and removing every runbook trait.
func HandleRunbookTraitsSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks/"+runbookID+"/traits", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, RunbookTraitsBody)
		case "PUT":
			th.TestJSONRequest(t, r, `{"traits": ["CUSTOM_GPU", "CUSTOM_FPGA"]}`)
			w.WriteHeader(http.StatusNoContent)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})
}

// HandleRunbookTraitSuccessfully sets up the test server to respond to
// adding and removing a single runbook trait.
func HandleRunbookTraitSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/runbooks/"+runbookID+"/traits/CUSTOM_GPU", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		switch r.Method {
		case "PUT":
			th.TestBody(t, r, "")
			w.WriteHeader(http.StatusNoContent)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})
}
