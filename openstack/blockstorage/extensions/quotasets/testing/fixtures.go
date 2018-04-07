package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const FirstTenantID = "555544443333222211110000ffffeeee"

var successTestCases = []struct {
	name             string
	httpMethod       string
	jsonBody         string
	uriPath          string
	uriQueryParams   map[string]string
	expectedQuotaSet quotasets.QuotaSet
}{
	{
		name: "simple GET request",
		jsonBody: `
{
	"quota_set" : {
		"volumes" : 8,
		"snapshots" : 9,
		"gigabytes" : 10,
		"per_volume_gigabytes" : 11,
		"backups" : 12,
		"backup_gigabytes" : 13
	}
}`,
		expectedQuotaSet: quotasets.QuotaSet{
			Volumes:            8,
			Snapshots:          9,
			Gigabytes:          10,
			PerVolumeGigabytes: 11,
			Backups:            12,
			BackupGigabytes:    13,
		},
		uriPath:    "/os-quota-sets/" + FirstTenantID,
		httpMethod: "GET",
	},
}

// HandleSuccessfulRequest configures the test server to respond to an HTTP request.
func HandleSuccessfulRequest(t *testing.T, httpMethod, uriPath, jsonOutput string) {

	th.Mux.HandleFunc(uriPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, httpMethod)
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, jsonOutput)
	})
}
