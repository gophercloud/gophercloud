package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/simpletenantusage"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
    "tenant_usages": [
        {
            "start": "2012-10-08T21:10:44.587336",
            "stop": "2012-10-08T22:10:44.587336",
            "tenant_id": "6f70656e737461636b20342065766572",
            "total_hours": 1.0,
            "total_local_gb_usage": 1.0,
            "total_memory_mb_usage": 512.0,
            "total_vcpus_usage": 1.0
        }
    ],
    "tenant_usages_links": [
        {
            "href": "http://openstack.example.com/v2.1/6f70656e737461636b20342065766572/os-simple-tenant-usage?end=2016-10-12+18%3A22%3A04.868106&limit=1&marker=1f1deceb-17b5-4c04-84c7-e0d4499c8fe0&start=2016-10-12+18%3A22%3A04.868106",
            "rel": "next"
        }
    ]
}
`

// GetTenantOutput is a sample response to a Get call for a specific tenant.
const GetTenantOutput = `
{
    "tenant_usage": {
        "server_usages": [
            {
                "ended_at": null,
                "flavor": "m1.tiny",
                "hours": 1.0,
                "instance_id": "1f1deceb-17b5-4c04-84c7-e0d4499c8fe0",
                "local_gb": 1,
                "memory_mb": 512,
                "name": "instance-2",
                "started_at": "2012-10-08T20:10:44.541277",
                "state": "active",
                "tenant_id": "6f70656e737461636b20342065766572",
                "uptime": 3600,
                "vcpus": 1
            }
        ],
        "start": "2012-10-08T20:10:44.587336",
        "stop": "2012-10-08T21:10:44.587336",
        "tenant_id": "6f70656e737461636b20342065766572",
        "total_hours": 1.0,
        "total_local_gb_usage": 1.0,
        "total_memory_mb_usage": 512.0,
        "total_vcpus_usage": 1.0
    },
    "tenant_usage_links": [
        {
            "href": "http://openstack.example.com/v2.1/6f70656e737461636b20342065766572/os-simple-tenant-usage/6f70656e737461636b20342065766572?end=2016-10-12+18%3A22%3A04.868106&limit=1&marker=1f1deceb-17b5-4c04-84c7-e0d4499c8fe0&start=2016-10-12+18%3A22%3A04.868106",
            "rel": "next"
        }
    ]
}
`

const FirstTenantID = "aabbccddeeff112233445566"

// SimpleTenantUsageResults is the decoded output corresponding to GetOutput above.
var SimpleTenantUsageResults = []simpletenantusage.TenantUsage{
	simpletenantusage.TenantUsage{
		ServerUsages:       []simpletenantusage.ServerUsage(nil),
		Start:              gophercloud.JSONRFC3339MilliNoZ(time.Date(2012, 10, 8, 21, 10, 44, 587336000, time.UTC)),
		Stop:               gophercloud.JSONRFC3339MilliNoZ(time.Date(2012, 10, 8, 22, 10, 44, 587336000, time.UTC)),
		TenantID:           "6f70656e737461636b20342065766572",
		TotalHours:         1.0,
		TotalLocalGBUsage:  1.0,
		TotalMemoryMBUsage: 512.0,
		TotalVCPUsUsage:    1.0,
	},
}

// SimpleTenantUsageOneTenantResults is the decoded output corresponding to GetTenantOutput above.
var SimpleTenantUsageOneTenantResults = simpletenantusage.TenantUsage{
	ServerUsages: []simpletenantusage.ServerUsage{
		simpletenantusage.ServerUsage{
			Flavor:     "m1.tiny",
			Hours:      1.0,
			InstanceID: "1f1deceb-17b5-4c04-84c7-e0d4499c8fe0",
			LocalGB:    1,
			MemoryMB:   512,
			Name:       "instance-2",
			StartedAt:  gophercloud.JSONRFC3339MilliNoZ(time.Date(2012, 10, 8, 20, 10, 44, 541277000, time.UTC)),
			State:      "active",
			TenantID:   "6f70656e737461636b20342065766572",
			Uptime:     3600,
			VCPUs:      1,
		},
	},
	Start:              gophercloud.JSONRFC3339MilliNoZ(time.Date(2012, 10, 8, 20, 10, 44, 587336000, time.UTC)),
	Stop:               gophercloud.JSONRFC3339MilliNoZ(time.Date(2012, 10, 8, 21, 10, 44, 587336000, time.UTC)),
	TenantID:           "6f70656e737461636b20342065766572",
	TotalHours:         1.0,
	TotalLocalGBUsage:  1.0,
	TotalMemoryMBUsage: 512.0,
	TotalVCPUsUsage:    1.0,
}

// HandleGetSuccessfully configures the test server to respond to a Get request
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-simple-tenant-usage/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetOutput)
	})
}

// HandleGetTenantSuccessfully configures the test server to respond to a Get request for sample tenant
func HandleGetTenantSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-simple-tenant-usage/"+FirstTenantID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetTenantOutput)
	})
}
