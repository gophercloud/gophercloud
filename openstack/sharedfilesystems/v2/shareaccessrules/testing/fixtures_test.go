package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const (
	shareAccessRulesEndpoint = "/share-access-rules"
	shareAccessRuleID        = "507bf114-36f2-4f56-8cf4-857985ca87c1"
	shareID                  = "fb213952-2352-41b4-ad7b-2c4c69d13eef"
)

var getResponse = `{
    "access": {
        "access_level": "rw",
        "state": "error",
        "id": "507bf114-36f2-4f56-8cf4-857985ca87c1",
        "share_id": "fb213952-2352-41b4-ad7b-2c4c69d13eef",
        "access_type": "cert",
        "access_to": "example.com",
        "access_key": null,
        "created_at": "2018-07-17T02:01:04.000000",
        "updated_at": "2018-07-17T02:01:04.000000",
        "metadata": {
            "key1": "value1",
            "key2": "value2"
        }
    }
}`

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc(shareAccessRulesEndpoint+"/"+shareAccessRuleID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getResponse)
	})
}

var listResponse = `{
    "access_list": [
        {
            "access_level": "rw",
            "state": "error",
            "id": "507bf114-36f2-4f56-8cf4-857985ca87c1",
            "access_type": "cert",
            "access_to": "example.com",
            "access_key": null,
            "created_at": "2018-07-17T02:01:04.000000",
            "updated_at": "2018-07-17T02:01:04.000000",
            "metadata": {
                "key1": "value1",
                "key2": "value2"
            }
        },
        {
            "access_level": "rw",
            "state": "active",
            "id": "a25b2df3-90bd-4add-afa6-5f0dbbd50452",
            "access_type": "ip",
            "access_to": "0.0.0.0/0",
            "access_key": null,
            "created_at": "2018-07-16T01:03:21.000000",
            "updated_at": "2018-07-16T01:03:21.000000",
            "metadata": {
                "key3": "value3",
                "key4": "value4"
            }
        }
    ]
}`
