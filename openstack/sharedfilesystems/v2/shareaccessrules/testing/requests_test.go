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
)

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponse(t, fakeServer)

	resp := shareaccessrules.Get(context.TODO(), client.ServiceClient(fakeServer), "507bf114-36f2-4f56-8cf4-857985ca87c1")
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
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}, accessRule)
}

func MockListResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc(shareAccessRulesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, listResponse)
	})
}

func TestUpdateMetadataSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUpdateMetadataResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for
	// access rule metadata is 2.45
	c.Microversion = "2.45"

	opts := shareaccessrules.UpdateMetadataOpts{
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	actual, err := shareaccessrules.UpdateMetadata(context.TODO(), c, shareAccessRuleID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, actual)
}

func TestDeleteMetadatumSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockDeleteMetadatumResponse(t, fakeServer, "key1")

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for
	// access rule metadata is 2.45
	c.Microversion = "2.45"

	err := shareaccessrules.DeleteMetadatum(context.TODO(), c, shareAccessRuleID, "key1").ExtractErr()
	th.AssertNoErr(t, err)
}
