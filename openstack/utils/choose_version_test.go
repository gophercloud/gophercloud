package utils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud/testhelper"
)

func TestChooseVersion(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
      {
        "versions": {
          "values": [
            {
              "status": "stable",
              "id": "v3.0",
              "links": [
                { "href": "https://example.com:1000/", "rel": "self" }
              ]
            },
            {
              "status": "stable",
              "id": "v2.0",
              "links": [
                { "href": "https://example.com:2000/", "rel": "self" }
              ]
            }
          ]
        }
      }
    `)
	})

	v2 := &Version{ID: "v2.0", Priority: 2}
	v3 := &Version{ID: "v3.0", Priority: 3}

	v, endpoint, err := ChooseVersion(testhelper.Endpoint(), []*Version{v2, v3})

	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v3 {
		t.Errorf("Expected %#v to win, but %#v did instead", v3, v)
	}

	expected := "https://example.com:1000/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}
