package testing

import (
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestDeleteNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/6dc6d336e3fc4c0a951b5698cd1236ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	deleteResult := nodes.Delete(fake.ServiceClient(), "6dc6d336e3fc4c0a951b5698cd1236ee")
	th.AssertNoErr(t, deleteResult.ExtractErr())
}
