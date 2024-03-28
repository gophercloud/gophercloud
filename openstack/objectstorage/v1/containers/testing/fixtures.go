package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

type handlerOptions struct {
	path string
}

type option func(*handlerOptions)

func WithPath(s string) option {
	return func(h *handlerOptions) {
		h.path = s
	}
}

// ExpectedListInfo is the result expected from a call to `List` when full
// info is requested.
var ExpectedListInfo = []containers.Container{
	{
		Count: 0,
		Bytes: 0,
		Name:  "janeausten",
	},
	{
		Count: 1,
		Bytes: 14,
		Name:  "marktwain",
	},
}

// ExpectedListNames is the result expected from a call to `List` when just
// container names are requested.
var ExpectedListNames = []string{"janeausten", "marktwain"}

// HandleListContainerInfoSuccessfully creates an HTTP handler at `/` on the test handler mux that
// responds with a `List` response when full info is requested.
func HandleListContainerInfoSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `[
        {
          "count": 0,
          "bytes": 0,
          "name": "janeausten"
        },
        {
          "count": 1,
          "bytes": 14,
          "name": "marktwain"
        }
      ]`)
		case "janeausten":
			fmt.Fprintf(w, `[
				{
					"count": 1,
					"bytes": 14,
					"name": "marktwain"
				}
			]`)
		case "marktwain":
			fmt.Fprintf(w, `[]`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

// HandleListZeroContainerNames204 creates an HTTP handler at `/` on the test handler mux that
// responds with "204 No Content" when container names are requested. This happens on some, but not all,
// objectstorage instances. This case is peculiar in that the server sends no `content-type` header.
func HandleListZeroContainerNames204(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleCreateContainerSuccessfully creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with a `Create` response.
func HandleCreateContainerSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("X-Container-Meta-Foo", "bar")
		w.Header().Set("Content-Length", "0")
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Header().Set("Date", "Wed, 17 Aug 2016 19:25:43 UTC")
		w.Header().Set("X-Trans-Id", "tx554ed59667a64c61866f1-0058b4ba37")
		w.Header().Set("X-Storage-Policy", "multi-region-three-replicas")
		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleDeleteContainerSuccessfully creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with a `Delete` response.
func HandleDeleteContainerSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})
}

const bulkDeleteResponse = `
{
    "Response Status": "foo",
    "Response Body": "bar",
    "Errors": [],
    "Number Deleted": 2,
    "Number Not Found": 0
}
`

// HandleBulkDeleteSuccessfully creates an HTTP handler at `/` on the test
// handler mux that responds with a `Delete` response.
func HandleBulkDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestFormValues(t, r, map[string]string{
			"bulk-delete": "true",
		})
		th.TestBody(t, r, "testContainer1\ntestContainer2\n")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, bulkDeleteResponse)
	})
}

// HandleUpdateContainerSuccessfully creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with a `Update` response.
func HandleUpdateContainerSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Container-Write", "")
		th.TestHeader(t, r, "X-Container-Read", "")
		th.TestHeader(t, r, "X-Container-Sync-To", "")
		th.TestHeader(t, r, "X-Container-Sync-Key", "")
		th.TestHeader(t, r, "Content-Type", "text/plain")
		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateContainerVersioningOn creates an HTTP handler at `/testVersioning` on the test handler mux that
// responds with a `Update` response.
func HandleUpdateContainerVersioningOn(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testVersioning",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Container-Write", "")
		th.TestHeader(t, r, "X-Container-Read", "")
		th.TestHeader(t, r, "X-Container-Sync-To", "")
		th.TestHeader(t, r, "X-Container-Sync-Key", "")
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestHeader(t, r, "X-Versions-Enabled", "true")
		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateContainerVersioningOff creates an HTTP handler at `/testVersioning` on the test handler mux that
// responds with a `Update` response.
func HandleUpdateContainerVersioningOff(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testVersioning",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Container-Write", "")
		th.TestHeader(t, r, "X-Container-Read", "")
		th.TestHeader(t, r, "X-Container-Sync-To", "")
		th.TestHeader(t, r, "X-Container-Sync-Key", "")
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestHeader(t, r, "X-Versions-Enabled", "false")
		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGetContainerSuccessfully creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with a `Get` response.
func HandleGetContainerSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Date", "Wed, 17 Aug 2016 19:25:43 UTC")
		w.Header().Set("X-Container-Bytes-Used", "100")
		w.Header().Set("X-Container-Object-Count", "4")
		w.Header().Set("X-Container-Read", "test")
		w.Header().Set("X-Container-Write", "test2,user4")
		w.Header().Set("X-Timestamp", "1471298837.95721")
		w.Header().Set("X-Trans-Id", "tx554ed59667a64c61866f1-0057b4ba37")
		w.Header().Set("X-Storage-Policy", "test_policy")
		w.Header().Set("X-Versions-Enabled", "True")
		w.WriteHeader(http.StatusNoContent)
	})
}
