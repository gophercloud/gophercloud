package servers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "bzbzbzbzbz"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestListServers(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, serverListBody)
		case "9e5476bd-a4ec-4653-93d6-72c93aa682ba":
			fmt.Fprintf(w, `{ "servers": [] }`)
		default:
			t.Fatalf("/servers/detail invoked with unexpected marker=[%s]", marker)
		}
	})

	client := serviceClient()
	pages := 0
	err := List(client).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractServers(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 servers, got %d", len(actual))
		}
		equalServers(t, serverHerp, actual[0])
		equalServers(t, serverDerp, actual[1])

		return true, nil
	})

	if err != nil {
		t.Fatalf("Unexpected error from EachPage: %v", err)
	}
	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestCreateServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestDeleteServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestGetServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestUpdateServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestChangeServerAdminPassword(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestRebootServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestRebuildServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestResizeServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestConfirmResize(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}

func TestRevertResize(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	t.Error("Pending")
}
