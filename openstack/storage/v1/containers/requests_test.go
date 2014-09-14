package containers

import (
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
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
	})

	client := serviceClient()
	_, err := List(client, ListOpts{Full: true})
	if err != nil {
		t.Fatalf("Unexpected error listing containers info: %v", err)
	}
}

func TestListContainerNames(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "text/plain")
	})

	client := serviceClient()
	_, err := List(client, ListOpts{})
	if err != nil {
		t.Fatalf("Unexpected error listing containers info: %v", err)
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
	_, err := Create(client, CreateOpts{
		Name: "testContainer",
	})
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
	err := Delete(client, DeleteOpts{
		Name: "testContainer",
	})
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
	err := Update(client, UpdateOpts{
		Name: "testContainer",
	})
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
	_, err := Get(client, GetOpts{
			Name: "testContainer",
	})
	if err != nil {
		t.Fatalf("Unexpected error getting container metadata: %v", err)
	}
}
