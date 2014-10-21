package servers

import (
	"fmt"
	"net/http"
	"testing"

	os "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestListServers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, ListOutput)
	})

	count := 0
	err := List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractServers(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedServerSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreateServer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"server": {
				"name": "derp",
				"imageRef": "f90f6034-2570-4974-8351-6b49732ef2eb",
				"flavorRef": "1"
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, CreateOutput)
	})

	client := client.ServiceClient()
	actual, err := Create(client, os.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &CreatedServer, actual)
}

func TestGetServer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/8c65cb68-0681-4c30-bc88-6b83a8a26aee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, GetOutput)
	})

	actual, err := Get(client.ServiceClient(), "8c65cb68-0681-4c30-bc88-6b83a8a26aee").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &GophercloudServer, actual)
}

func TestChangeAdminPassword(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/1234asdf/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{ "changePassword": { "adminPass": "new-password" } }`)

		w.WriteHeader(http.StatusAccepted)
	})

	res := ChangeAdminPassword(client.ServiceClient(), "1234asdf", "new-password")
	th.AssertNoErr(t, res.Err)
}
