package images

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestListImageDetails(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/images/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, ListOutput)
		case "e19a734c-c7e6-443a-830c-242209c4d65d":
			fmt.Fprintf(w, `{ "images": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	count := 0
	err := ListDetail(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractImages(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedImageSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	if count < 1 {
		t.Fatalf("At least one page of image results expected.")
	}
}
