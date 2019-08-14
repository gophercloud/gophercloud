package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tags"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/uuid1/tags", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintf(w, TagsListResponse)
		th.AssertNoErr(t, err)
	})

	expected := []string{"foo", "bar", "baz"}
	actual, err := tags.List(client.ServiceClient(), "uuid1").Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}
