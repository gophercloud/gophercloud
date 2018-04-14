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
	name                  string
	httpMethod            string
	jsonBody              string
	uriPath               string
	uriQueryParams        map[string]string
	expectedQuotaSet      quotasets.QuotaSet
	expectedQuotaUsageSet quotasets.QuotaUsageSet
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

	{
		name: "GET details request",
		jsonBody: `
{
	"quota_set" : {
		"id": "555544443333222211110000ffffeeee",
		"volumes" : {
			"in_use": 15,
			"limit": 16,
			"reserved": 17
		},
		"snapshots" : {
			"in_use": 18,
			"limit": 19,
			"reserved": 20
		},
		"gigabytes" : {
			"in_use": 21,
			"limit": 22,
			"reserved": 23
		},
		"per_volume_gigabytes" : {
			"in_use": 24,
			"limit": 25,
			"reserved": 26
		},
		"backups" : {
			"in_use": 27,
			"limit": 28,
			"reserved": 29
		},
		"backup_gigabytes" : {
			"in_use": 30,
			"limit": 31,
			"reserved": 32
		}
		}
	}
}`,
		expectedQuotaUsageSet: quotasets.QuotaUsageSet{
			ID:                 FirstTenantID,
			Volumes:            quotasets.QuotaUsage{InUse: 15, Limit: 16, Reserved: 17},
			Snapshots:          quotasets.QuotaUsage{InUse: 18, Limit: 19, Reserved: 20},
			Gigabytes:          quotasets.QuotaUsage{InUse: 21, Limit: 22, Reserved: 23},
			PerVolumeGigabytes: quotasets.QuotaUsage{InUse: 24, Limit: 25, Reserved: 26},
			Backups:            quotasets.QuotaUsage{InUse: 27, Limit: 28, Reserved: 29},
			BackupGigabytes:    quotasets.QuotaUsage{InUse: 30, Limit: 31, Reserved: 32},
		},
		uriPath:        "/os-quota-sets/" + FirstTenantID,
		uriQueryParams: map[string]string{"usage": "true"},
		httpMethod:     "GET",
	},
}

// HandleSuccessfulRequest configures the test server to respond to an HTTP request.
func HandleSuccessfulRequest(t *testing.T, httpMethod, uriPath, jsonOutput string, uriQueryParams map[string]string) {

	th.Mux.HandleFunc(uriPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, httpMethod)
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")

		if uriQueryParams != nil {
			th.TestFormValues(t, r, uriQueryParams)
		}

		fmt.Fprintf(w, jsonOutput)
	})
}
