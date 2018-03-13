package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/buildinfo"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListBuildInfos(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
    "build_info": [
        {
					  "api": {
							"revision": "1.0"
						},
						"engine": {
							"revision": "2.0"
						}
        }
    ]
		}`)
	})

	rBuildInfo := buildinfo.Get(fake.ServiceClient())
	actual, err := rBuildInfo.ExtractBuildInfo()
	if err != nil {
		t.Errorf("Failed to extract build-info: %V", err)
		return
	}

	expected := []buildinfo.BuildInfo{
		{
			API:    map[string]interface{}{"revision": "1.0"},
			Engine: map[string]interface{}{"revision": "2.0"},
		},
	}

	th.AssertDeepEquals(t, expected, actual)
}

func TestNonJSONCannotBeExtractedIntoBuildInfos(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rBuildInfo := buildinfo.Get(fake.ServiceClient())
	_, err := rBuildInfo.ExtractBuildInfo()
	if err == nil {
		t.Fatalf("Expected error, got nil")
		return
	}
}
