package flavors

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "blerb"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestListFlavors(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/flavors/detail", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `
					{
						"flavors": [
							{
								"id": "1",
								"name": "m1.tiny",
								"disk": 1,
								"ram": 512,
								"vcpus": 1
							},
							{
								"id": "2",
								"name": "m2.small",
								"disk": 10,
								"ram": 1024,
								"vcpus": 2
							}
						]
					}
				`)
		case "2":
			fmt.Fprintf(w, `{ "flavors": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	client := serviceClient()
	pages := 0
	err := List(client, ListFilterOptions{}).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		expected := []Flavor{
			Flavor{ID: "1", Name: "m1.tiny", Disk: 1, RAM: 512, VCPUs: 1},
			Flavor{ID: "2", Name: "m2.small", Disk: 10, RAM: 1024, VCPUs: 2},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %#v, but was %#v", expected, actual)
		}

		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
