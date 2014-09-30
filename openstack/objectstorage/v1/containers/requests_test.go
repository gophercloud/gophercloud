package containers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/testhelper"
)

const (
	tokenId = "abcabcabcabc"
)

var metadata = map[string]string{"gophercloud-test": "containers"}

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenId},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestListContainerInfo(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
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
		case "marktwain":
			fmt.Fprintf(w, `[]`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	client := serviceClient()
	count := 0
	List(client, &ListOpts{Full: true}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractInfo(page)
		if err != nil {
			t.Errorf("Failed to extract container info: %v", err)
			return false, err
		}

		expected := []Container{
			Container{
				Count: 0,
				Bytes: 0,
				Name:  "janeausten",
			},
			Container{
				Count: 1,
				Bytes: 14,
				Name:  "marktwain",
			},
		}

		testhelper.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListContainerNames(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "text/plain")

		w.Header().Set("Content-Type", "text/plain")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, "janeausten\nmarktwain\n")
		case "marktwain":
			fmt.Fprintf(w, ``)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	client := serviceClient()
	count := 0
	List(client, &ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract container names: %v", err)
			return false, err
		}

		expected := []string{"janeausten", "marktwain"}

		testhelper.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreateContainer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "PUT")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	_, err := Create(client, "testContainer", nil)
	if err != nil {
		t.Fatalf("Unexpected error creating container: %v", err)
	}
}

func TestDeleteContainer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	err := Delete(client, "testContainer")
	if err != nil {
		t.Fatalf("Unexpected error deleting container: %v", err)
	}
}

func TestUpateContainer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	err := Update(client, "testContainer", nil)
	if err != nil {
		t.Fatalf("Unexpected error updating container metadata: %v", err)
	}
}

func TestGetContainer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "HEAD")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	_, err := Get(client, "testContainer").ExtractMetadata()
	if err != nil {
		t.Fatalf("Unexpected error getting container metadata: %v", err)
	}
}
