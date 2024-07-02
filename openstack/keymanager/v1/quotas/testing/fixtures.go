package testing

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/quotas"
)

const GetResponseRaw_1 = `
{
    "quotas": {
      "secrets": 10,
      "orders": 20,
      "containers": -1,
      "consumers": 10,
      "cas": 5
    }
}
`

var GetResponse = quotas.Quota{
	Secrets:    gophercloud.IntToPointer(10),
	Orders:     gophercloud.IntToPointer(20),
	Containers: gophercloud.IntToPointer(-1),
	Consumers:  gophercloud.IntToPointer(10),
	CAS:        gophercloud.IntToPointer(5),
}

const GetProjectResponseRaw_1 = `
{
    "project_quotas": {
      "secrets": 10,
      "orders": 20,
      "containers": -1,
      "consumers": 10,
      "cas": 5
    }
}
`

const ListResponseRaw_1 = `
{
    "project_quotas": [
      {
        "project_id": "1234",
        "project_quotas": {
             "secrets": 2000,
             "orders": 0,
             "containers": -1,
             "consumers": null,
             "cas": null
         }
      },
      {
        "project_id": "5678",
        "project_quotas": {
             "secrets": 200,
             "orders": 100,
             "containers": -1,
             "consumers": null,
             "cas": null
         }
      }
    ],
    "total" : 30
  }
`

var ExpectedQuotasSlice = []quotas.ProjectQuota{
	{
		ProjectID: "1234",
		Quota: quotas.Quota{
			Secrets:    gophercloud.IntToPointer(2000),
			Orders:     gophercloud.IntToPointer(0),
			Containers: gophercloud.IntToPointer(-1),
			Consumers:  nil,
			CAS:        nil,
		},
	},
	{
		ProjectID: "5678",
		Quota: quotas.Quota{
			Secrets:    gophercloud.IntToPointer(200),
			Orders:     gophercloud.IntToPointer(100),
			Containers: gophercloud.IntToPointer(-1),
			Consumers:  nil,
			CAS:        nil,
		},
	},
}
