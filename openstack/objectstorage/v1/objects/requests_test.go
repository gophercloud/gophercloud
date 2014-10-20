package objects

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var metadata = map[string]string{"Gophercloud-Test": "objects"}

func TestDownloadObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successful download with Gophercloud")
	})

	content, err := Download(fake.ServiceClient(), "testContainer", "testObject", nil).ExtractContent()
	if err != nil {
		t.Fatalf("Unexpected error downloading object: %v", err)
	}
	if string(content) != "Successful download with Gophercloud" {
		t.Errorf("Expected %s, got %s", "Successful download with Gophercloud", content)
	}
}

func TestListObjectInfo(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `[
				{
					"hash": "451e372e48e0f6b1114fa0724aa79fa1",
					"last_modified": "2009-11-10 23:00:00 +0000 UTC",
					"bytes": 14,
					"name": "goodbye",
					"content_type": "application/octet-stream"
				},
				{
					"hash": "451e372e48e0f6b1114fa0724aa79fa1",
					"last_modified": "2009-11-10 23:00:00 +0000 UTC",
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

	count := 0

	err := List(fake.ServiceClient(), "testContainer", &ListOpts{Full: true}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractInfo(page)
		if err != nil {
			t.Errorf("Failed to extract object info: %v", err)
			return false, err
		}

		expected := []Object{
			Object{
				Hash:         "451e372e48e0f6b1114fa0724aa79fa1",
				LastModified: "2009-11-10 23:00:00 +0000 UTC",
				Bytes:        14,
				Name:         "goodbye",
				ContentType:  "application/octet-stream",
			},
			Object{
				Hash:         "451e372e48e0f6b1114fa0724aa79fa1",
				LastModified: "2009-11-10 23:00:00 +0000 UTC",
				Bytes:        14,
				Name:         "hello",
				ContentType:  "application/octet-stream",
			},
		}

		testhelper.CheckDeepEquals(t, actual, expected)

		return true, nil
	})
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListObjectNames(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "text/plain")

		w.Header().Set("Content-Type", "text/plain")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, "hello\ngoodbye\n")
		case "goodbye":
			fmt.Fprintf(w, "")
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	count := 0
	List(fake.ServiceClient(), "testContainer", &ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract object names: %v", err)
			return false, err
		}

		expected := []string{"hello", "goodbye"}

		testhelper.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreateObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "PUT")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	content := bytes.NewBufferString("Did gyre and gimble in the wabe")
	options := &CreateOpts{ContentType: "application/json"}
	_, err := Create(fake.ServiceClient(), "testContainer", "testObject", content, options).Extract()
	if err != nil {
		t.Fatalf("Unexpected error creating object: %v", err)
	}
}

func TestCopyObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "COPY")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "Destination", "/newTestContainer/newTestObject")
		w.WriteHeader(http.StatusCreated)
	})

	options := &CopyOpts{Destination: "/newTestContainer/newTestObject"}
	_, err := Copy(fake.ServiceClient(), "testContainer", "testObject", options).Extract()
	if err != nil {
		t.Fatalf("Unexpected error copying object: %v", err)
	}
}

func TestDeleteObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := Delete(fake.ServiceClient(), "testContainer", "testObject", nil).Extract()
	if err != nil {
		t.Fatalf("Unexpected error deleting object: %v", err)
	}
}

func TestUpateObjectMetadata(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Object-Meta-Gophercloud-Test", "objects")
		w.WriteHeader(http.StatusAccepted)
	})

	_, err := Update(fake.ServiceClient(), "testContainer", "testObject", &UpdateOpts{Metadata: metadata}).Extract()
	if err != nil {
		t.Fatalf("Unexpected error updating object metadata: %v", err)
	}
}

func TestGetObject(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/testContainer/testObject", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "HEAD")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("X-Object-Meta-Gophercloud-Test", "objects")
		w.WriteHeader(http.StatusNoContent)
	})

	expected := metadata
	actual, err := Get(fake.ServiceClient(), "testContainer", "testObject", nil).ExtractMetadata()
	if err != nil {
		t.Fatalf("Unexpected error getting object metadata: %v", err)
	}
	testhelper.CheckDeepEquals(t, expected, actual)
}
