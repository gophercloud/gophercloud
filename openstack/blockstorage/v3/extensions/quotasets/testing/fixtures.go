package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const FirstTenantID = "555544443333222211110000ffffeeee"

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
	"quota_set" : {
		"volumes" : 8,
		"snapshots" : 9,
		"gigabytes" : 10,
		"per_volume_gigabytes" : 11,
		"backups" : 12,
		"backup_gigabytes" : 13,
		"groups" : 14
	}
}
`

// FirstQuotaSet is the first result in ListOutput.
var FirstQuotaSet = quotasets.QuotaSet{
	Volumes:            8,
	Snapshots:          9,
	Gigabytes:          10,
	PerVolumeGigabytes: 11,
	Backups:            12,
	BackupGigabytes:    13,
	Groups:             14,
}

// GetDetailsOutput is a sample response to a Get call with the detailed option.
const GetDetailsOutput = `
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
		},
		"groups" : {
			"in_use": 33,
			"limit": 34,
			"reserved": 35
		}
	}
}
`

// FirstQuotaDetailsSet is the first result in ListOutput.
var FirstQuotaDetailsSet = quotasets.QuotaDetailSet{
	ID:                 FirstTenantID,
	Volumes:            quotasets.QuotaDetail{InUse: 15, Limit: 16, Reserved: 17},
	Snapshots:          quotasets.QuotaDetail{InUse: 18, Limit: 19, Reserved: 20},
	Gigabytes:          quotasets.QuotaDetail{InUse: 21, Limit: 22, Reserved: 23},
	PerVolumeGigabytes: quotasets.QuotaDetail{InUse: 24, Limit: 25, Reserved: 26},
	Backups:            quotasets.QuotaDetail{InUse: 27, Limit: 28, Reserved: 29},
	BackupGigabytes:    quotasets.QuotaDetail{InUse: 30, Limit: 31, Reserved: 32},
	Groups:             quotasets.QuotaDetail{InUse: 33, Limit: 34, Reserved: 35},
}

//The expected update Body. Is also returned by PUT request
const UpdateOutput = `{"quota_set":{"volumes":8,"snapshots":9,"gigabytes":10,"per_volume_gigabytes":11,"backups":12,"backup_gigabytes":13,"groups":14}}`

//Result of Quota-update
var UpdatedQuotaSet = quotasets.UpdateOpts{
	Volumes:            gophercloud.IntToPointer(8),
	Snapshots:          gophercloud.IntToPointer(9),
	Gigabytes:          gophercloud.IntToPointer(10),
	PerVolumeGigabytes: gophercloud.IntToPointer(11),
	Backups:            gophercloud.IntToPointer(12),
	BackupGigabytes:    gophercloud.IntToPointer(13),
	Groups:             gophercloud.IntToPointer(14),
}

//The expected partialupdate Body. Is also returned by PUT request
const PartialUpdateBody = `{"quota_set":{"volumes":200, "force":true}}`

// PartialQuotaSet something something
var PartialQuotaSet = quotasets.QuotaSet{
	Volumes:            200,
	Snapshots:          0,
	Gigabytes:          0,
	PerVolumeGigabytes: 0,
	Backups:            0,
	BackupGigabytes:    0,
	Groups:             0,
}

// HandleGetSuccessfully configures the test server to respond to a Get request for sample tenant
func HandleGetSuccessfully(t *testing.T, uriPath, jsonOutput string) {
	th.Mux.HandleFunc(uriPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, jsonOutput)
	})
}

// HandlePutSuccessfully configures the test server to respond to a Put request for sample tenant
func HandlePutSuccessfully(t *testing.T, uriPath, jsonOutput string) {
	th.Mux.HandleFunc(uriPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, jsonOutput)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, jsonOutput)
	})
}

// HandleDeleteSuccessfully configures the test server to respond to a Delete request for sample tenant
func HandleDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-quota-sets/"+FirstTenantID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestBody(t, r, "")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(202)
	})
}
