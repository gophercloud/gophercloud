package objects

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

const ( 
	tokenId = "abcabcabcabc"
)

var metadata = map[string]string{"gophercloud-test": "objects"}

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenId},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestDownloadObject(t * testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	
	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
	})

	client := serviceClient()
	_, err := Download(client, DownloadOpts{
		Container: "testContainer",
		Name: "testObject",
	})
	if err != nil {
		t.Fatalf("Unexpected error downloading object: %v", err)
	}
}

func TestListObjectInfo(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
	})

	client := serviceClient()
	_, err := List(client, ListOpts{
		Full: true,
		Container: "testContainer",
	})
	if err != nil {
		t.Fatalf("Unexpected error listing objects info: %v", err)
	}
}

func TestListObjectNames(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "text/plain")
	})

	client := serviceClient()
	_, err := List(client, ListOpts{
		Container: "testContainer",
	})
	if err != nil {
		t.Fatalf("Unexpected error listing object names: %v", err)
	}
}

func TestCreateObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "PUT")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	client := serviceClient()
	err := Create(client, CreateOpts{
		Content: bytes.NewBufferString("Did gyre and gimble in the wabe:"),
		Container: "testContainer",
		Name: "testObject",
	})
	if err != nil {
		t.Fatalf("Unexpected error creating object: %v", err)
	}
}

func TestCopyObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "COPY")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "Destination", "/newTestContainer/newTestObject")
		w.WriteHeader(http.StatusCreated)
	})

	client := serviceClient()
	err := Copy(client, CopyOpts{
		NewContainer: "newTestContainer",
		NewName: "newTestObject",
		Container: "testContainer",
		Name: "testObject",
	})
	if err != nil {
		t.Fatalf("Unexpected error copying object: %v", err)
	}
}

func TestDeleteObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	err := Delete(client, DeleteOpts{
		Container: "testContainer",
		Name: "testObject",
	})
	if err != nil {
		t.Fatalf("Unexpected error deleting object: %v", err)
	}
}

func TestUpateObjectMetadata(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Object-Meta-Gophercloud-Test", "objects")
		w.WriteHeader(http.StatusAccepted)
	})

	client := serviceClient()
	err := Update(client, UpdateOpts{
		Container: "testContainer",
		Name: "testObject",
		Metadata: metadata,
	})
	if err != nil {
		t.Fatalf("Unexpected error updating object metadata: %v", err)
	}
}

func TestGetContainer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "HEAD")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	_, err := Get(client, GetOpts{
			Container: "testContainer",
			Name: "testObject",
	})
	if err != nil {
		t.Fatalf("Unexpected error getting object metadata: %v", err)
	}
}
