package images

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "aaaaaa"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestListImages(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/images/detail", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `
				{
					"images": [
						{
							"status": "ACTIVE",
							"updated": "2014-09-23T12:54:56Z",
							"id": "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
							"OS-EXT-IMG-SIZE:size": 476704768,
							"name": "F17-x86_64-cfntools",
							"created": "2014-09-23T12:54:52Z",
							"minDisk": 0,
							"progress": 100,
							"minRam": 0,
							"metadata": {}
						},
						{
							"status": "ACTIVE",
							"updated": "2014-09-23T12:51:43Z",
							"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
							"OS-EXT-IMG-SIZE:size": 13167616,
							"name": "cirros-0.3.2-x86_64-disk",
							"created": "2014-09-23T12:51:42Z",
							"minDisk": 0,
							"progress": 100,
							"minRam": 0,
							"metadata": {}
						}
					]
				}
			`)
		case "2":
			fmt.Fprintf(w, `{ "images": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	client := serviceClient()
	pages := 0
	err := List(client).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractImages(page)
		if err != nil {
			return false, err
		}

		expected := []Image{
			Image{
				ID:       "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
				Name:     "F17-x86_64-cfntools",
				Created:  "2014-09-23T12:54:52Z",
				Updated:  "2014-09-23T12:54:56Z",
				MinDisk:  0,
				MinRAM:   0,
				Progress: 100,
				Status:   "ACTIVE",
			},
			Image{
				ID:       "f90f6034-2570-4974-8351-6b49732ef2eb",
				Name:     "cirros-0.3.2-x86_64-disk",
				Created:  "2014-09-23T12:51:42Z",
				Updated:  "2014-09-23T12:51:43Z",
				MinDisk:  0,
				MinRAM:   0,
				Progress: 100,
				Status:   "ACTIVE",
			},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Unexpected page contents: expected %#v, got %#v", expected, actual)
		}

		return false, nil
	})

	if err != nil {
		t.Fatalf("EachPage error: %v", err)
	}
	if pages != 1 {
		t.Errorf("Expected one page, got %d", pages)
	}
}

func TestGetImage(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/images/12345678", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"image": {
					"status": "ACTIVE",
					"updated": "2014-09-23T12:54:56Z",
					"id": "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
					"OS-EXT-IMG-SIZE:size": 476704768,
					"name": "F17-x86_64-cfntools",
					"created": "2014-09-23T12:54:52Z",
					"minDisk": 0,
					"progress": 100,
					"minRam": 0,
					"metadata": {}
				}
			}
		`)
	})

	client := serviceClient()
	actual, err := Get(client, "12345678").Extract()
	if err != nil {
		t.Fatalf("Unexpected error from Get: %v", err)
	}

	expected := &Image{
		Status:   "ACTIVE",
		Updated:  "2014-09-23T12:54:56Z",
		ID:       "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
		Name:     "F17-x86_64-cfntools",
		Created:  "2014-09-23T12:54:52Z",
		MinDisk:  0,
		Progress: 100,
		MinRAM:   0,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but got %#v", expected, actual)
	}
}
