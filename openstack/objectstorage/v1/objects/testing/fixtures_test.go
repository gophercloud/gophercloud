package testing

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
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

// HandleDownloadObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux that
// responds with a `Download` response.
func HandleDownloadObjectSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer/testObject",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Set("Date", date.Format(time.RFC1123))
		w.Header().Set("X-Static-Large-Object", "True")

		unModifiedSince := r.Header.Get("If-Unmodified-Since")
		modifiedSince := r.Header.Get("If-Modified-Since")
		if unModifiedSince != "" {
			ums, _ := time.Parse(time.RFC1123, unModifiedSince)
			if ums.Before(date) || ums.Equal(date) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
		}
		if modifiedSince != "" {
			ms, _ := time.Parse(time.RFC1123, modifiedSince)
			if ms.After(date) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Header().Set("Last-Modified", date.Format(time.RFC1123))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successful download with Gophercloud")
	})
}

// ExpectedListInfo is the result expected from a call to `List` when full
// info is requested.
var ExpectedListInfo = []objects.Object{
	{
		Hash:         "451e372e48e0f6b1114fa0724aa79fa1",
		LastModified: time.Date(2016, time.August, 17, 22, 11, 58, 602650000, time.UTC),
		Bytes:        14,
		Name:         "goodbye",
		ContentType:  "application/octet-stream",
	},
	{
		Hash:         "451e372e48e0f6b1114fa0724aa79fa1",
		LastModified: time.Date(2016, time.August, 17, 22, 11, 58, 602650000, time.UTC),
		Bytes:        14,
		Name:         "hello",
		ContentType:  "application/octet-stream",
	},
}

// ExpectedListSubdir is the result expected from a call to `List` when full
// info is requested.
var ExpectedListSubdir = []objects.Object{
	{
		Subdir: "directory/",
	},
}

// ExpectedListNames is the result expected from a call to `List` when just
// object names are requested.
var ExpectedListNames = []string{"goodbye", "hello"}

// HandleListObjectsInfoSuccessfully creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with a `List` response when full info is requested.
func HandleListObjectsInfoSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
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
        "hash": "451e372e48e0f6b1114fa0724aa79fa1",
        "last_modified": "2016-08-17T22:11:58.602650",
        "bytes": 14,
        "name": "goodbye",
        "content_type": "application/octet-stream"
      },
      {
        "hash": "451e372e48e0f6b1114fa0724aa79fa1",
        "last_modified": "2016-08-17T22:11:58.602650",
        "bytes": 14,
        "name": "hello",
        "content_type": "application/octet-stream"
      }
    ]`)
		case "hello":
			fmt.Fprintf(w, `[]`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

// HandleListSubdirSuccessfully creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with a `List` response when full info is requested.
func HandleListSubdirSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
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
        "subdir": "directory/"
      }
    ]`)
		case "directory/":
			fmt.Fprintf(w, `[]`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

// HandleListZeroObjectNames204 creates an HTTP handler at `/testContainer` on the test handler mux that
// responds with "204 No Content" when object names are requested. This happens on some, but not all, objectstorage
// instances. This case is peculiar in that the server sends no `content-type` header.
func HandleListZeroObjectNames204(t *testing.T) {
	th.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleCreateTextObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux
// that responds with a `Create` response. A Content-Type of "text/plain" is expected.
func HandleCreateTextObjectSuccessfully(t *testing.T, content string, options ...option) {
	ho := handlerOptions{
		path: "/testContainer/testObject",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestBody(t, r, `Did gyre and gimble in the wabe`)

		hash := md5.New()
		_, err := io.WriteString(hash, content)
		th.AssertNoErr(t, err)
		localChecksum := hash.Sum(nil)

		w.Header().Set("ETag", fmt.Sprintf("%x", localChecksum))
		w.WriteHeader(http.StatusCreated)
	})
}

// HandleCreateTextWithCacheControlSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler
// mux that responds with a `Create` response. A Cache-Control of `max-age="3600", public` is expected.
func HandleCreateTextWithCacheControlSuccessfully(t *testing.T, content string) {
	th.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Cache-Control", `max-age="3600", public`)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestBody(t, r, `All mimsy were the borogoves`)

		hash := md5.New()
		_, err := io.WriteString(hash, content)
		th.AssertNoErr(t, err)
		localChecksum := hash.Sum(nil)

		w.Header().Set("ETag", fmt.Sprintf("%x", localChecksum))
		w.WriteHeader(http.StatusCreated)
	})
}

// HandleCreateTypelessObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler
// mux that responds with a `Create` response. No Content-Type header may be present in the request, so that server-
// side content-type detection will be triggered properly.
func HandleCreateTypelessObjectSuccessfully(t *testing.T, content string) {
	th.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestBody(t, r, `The sky was the color of television, tuned to a dead channel.`)

		if contentType, present := r.Header["Content-Type"]; present {
			t.Errorf("Expected Content-Type header to be omitted, but was %#v", contentType)
		}

		hash := md5.New()
		_, err := io.WriteString(hash, content)
		th.AssertNoErr(t, err)
		localChecksum := hash.Sum(nil)

		w.Header().Set("ETag", fmt.Sprintf("%x", localChecksum))
		w.WriteHeader(http.StatusCreated)
	})
}

// HandleCopyObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux that
// responds with a `Copy` response.
func HandleCopyObjectSuccessfully(t *testing.T, destination string) {
	th.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "COPY")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Destination", destination)
		w.WriteHeader(http.StatusCreated)
	})
}

// HandleCopyObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux that
// responds with a `Copy` response.
func HandleCopyObjectVersionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "COPY")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Destination", "/newTestContainer/newTestObject")
		th.TestFormValues(t, r, map[string]string{"version-id": "123456788"})
		w.Header().Set("X-Object-Version-Id", "123456789")
		w.WriteHeader(http.StatusCreated)
	})
}

// HandleDeleteObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux that
// responds with a `Delete` response.
func HandleDeleteObjectSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer/testObject",
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
// handler mux that responds with a `BulkDelete` response.
func HandleBulkDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestFormValues(t, r, map[string]string{
			"bulk-delete": "true",
		})
		th.TestBody(t, r, "testContainer/testObject1\ntestContainer/testObject2\n")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, bulkDeleteResponse)
	})
}

// HandleUpdateObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux that
// responds with a `Update` response.
func HandleUpdateObjectSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer/testObject",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Object-Meta-Gophercloud-Test", "objects")
		th.TestHeader(t, r, "X-Remove-Object-Meta-Gophercloud-Test-Remove", "remove")
		th.TestHeader(t, r, "Content-Disposition", "")
		th.TestHeader(t, r, "Content-Encoding", "")
		th.TestHeader(t, r, "Content-Type", "")
		th.TestHeaderUnset(t, r, "X-Delete-After")
		th.TestHeader(t, r, "X-Delete-At", "0")
		th.TestHeader(t, r, "X-Detect-Content-Type", "false")
		w.WriteHeader(http.StatusAccepted)
	})
}

// HandleGetObjectSuccessfully creates an HTTP handler at `/testContainer/testObject` on the test handler mux that
// responds with a `Get` response.
func HandleGetObjectSuccessfully(t *testing.T, options ...option) {
	ho := handlerOptions{
		path: "/testContainer/testObject",
	}
	for _, apply := range options {
		apply(&ho)
	}

	th.Mux.HandleFunc(ho.path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("X-Object-Meta-Gophercloud-Test", "objects")
		w.Header().Add("X-Static-Large-Object", "true")
		w.WriteHeader(http.StatusNoContent)
	})
}
