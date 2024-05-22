package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shareaccessrules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	resp := shareaccessrules.Get(context.TODO(), client.ServiceClient(), "507bf114-36f2-4f56-8cf4-857985ca87c1")
	th.AssertNoErr(t, resp.Err)

	accessRule, err := resp.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &shareaccessrules.ShareAccess{
		ShareID:     "fb213952-2352-41b4-ad7b-2c4c69d13eef",
		CreatedAt:   time.Date(2018, 7, 17, 2, 1, 4, 0, time.UTC),
		UpdatedAt:   time.Date(2018, 7, 17, 2, 1, 4, 0, time.UTC),
		AccessType:  "cert",
		AccessTo:    "example.com",
		AccessKey:   "",
		State:       "error",
		AccessLevel: "rw",
		ID:          "507bf114-36f2-4f56-8cf4-857985ca87c1",
		Metadata: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}, accessRule)
}

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc(shareAccessRulesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, listResponse)
	})
}
