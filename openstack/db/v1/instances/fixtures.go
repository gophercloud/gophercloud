package instances

import (
	"fmt"

	"github.com/rackspace/gophercloud"
)

const instance = `
{
  "created": "2014-02-13T21:47:13",
  "datastore": {
    "type": "mysql",
    "version": "5.6"
  },
  "flavor": {
    "id": "1",
    "links": [
      {
        "href": "https://my-openstack.com/v1.0/1234/flavors/1",
        "rel": "self"
      },
      {
        "href": "https://my-openstack.com/v1.0/1234/flavors/1",
        "rel": "bookmark"
      }
    ]
  },
  "links": [
    {
      "href": "https://my-openstack.com/v1.0/1234/instances/1",
      "rel": "self"
    }
  ],
  "hostname": "e09ad9a3f73309469cf1f43d11e79549caf9acf2.my-openstack.com",
  "id": "{instanceID}",
  "name": "json_rack_instance",
  "status": "BUILD",
  "updated": "2014-02-13T21:47:13",
  "volume": {
    "size": 2
  }
}
`

var createReq = `
{
	"instance": {
		"databases": [
			{
				"character_set": "utf8",
				"collate": "utf8_general_ci",
				"name": "sampledb"
			},
			{
				"name": "nextround"
			}
		],
		"flavorRef": "1",
		"name": "json_rack_instance",
		"users": [
			{
				"databases": [
					{
						"name": "sampledb"
					}
				],
				"name": "demouser",
				"password": "demopassword"
			}
		],
		"volume": {
			"size": 2
		}
	}
}
`

var (
	restartReq   = `{"restart": true}`
	resizeReq    = `{"resize": {"flavorRef": "2"}}`
	resizeVolReq = `{"resize": {"volume": {"size": 4}}}`
)

var (
	createResp        = fmt.Sprintf(`{"instance": %s}`, instance)
	listInstancesResp = fmt.Sprintf(`{"instances":[%s]}`, instance)
	getInstanceResp   = createResp
	enableUserResp    = `{"user":{"name":"root","password":"secretsecret"}}`
	isUserEnabledResp = `{"rootEnabled":true}`
)

var expectedInstance = Instance{
	Created: "2014-02-13T21:47:13",
	Updated: "2014-02-13T21:47:13",
	Flavor: Flavor{
		ID: "1",
		Links: []gophercloud.Link{
			gophercloud.Link{Href: "https://my-openstack.com/v1.0/1234/flavors/1", Rel: "self"},
			gophercloud.Link{Href: "https://my-openstack.com/v1.0/1234/flavors/1", Rel: "bookmark"},
		},
	},
	Hostname: "e09ad9a3f73309469cf1f43d11e79549caf9acf2.my-openstack.com",
	ID:       instanceID,
	Links: []gophercloud.Link{
		gophercloud.Link{Href: "https://my-openstack.com/v1.0/1234/instances/1", Rel: "self"},
	},
	Name:   "json_rack_instance",
	Status: "BUILD",
	Volume: Volume{Size: 2},
}
