package testing

import (
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

var (
	timestamp  = "2015-11-12T14:22:42"
	timeVal, _ = time.Parse(gophercloud.RFC3339NoZ, timestamp)
)

var instance = `
{
  "created": "` + timestamp + `",
  "datastore": {
    "type": "mysql",
    "version": "5.6"
  },
  "flavor": {
    "id": "1",
    "links": [
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "self"
      },
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "bookmark"
      }
    ]
  },
  "links": [
    {
      "href": "https://openstack.example.com/v1.0/1234/instances/1",
      "rel": "self"
    }
  ],
  "hostname": "e09ad9a3f73309469cf1f43d11e79549caf9acf2.openstack.example.com",
  "id": "{instanceID}",
  "name": "json_rack_instance",
  "status": "BUILD",
  "updated": "` + timestamp + `",
  "volume": {
    "size": 2
  }
}
`

var instanceGet = `
{
  "created": "` + timestamp + `",
  "datastore": {
    "type": "mysql",
    "version": "5.6"
  },
  "flavor": {
    "id": "1",
    "links": [
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "self"
      },
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "bookmark"
      }
    ]
  },
  "links": [
    {
      "href": "https://openstack.example.com/v1.0/1234/instances/1",
      "rel": "self"
    }
  ],
  "id": "{instanceID}",
  "name": "test",
  "status": "ACTIVE",
  "operating_status": "HEALTHY",
  "updated": "` + timestamp + `",
  "volume": {
    "size": 1,
    "used": 0.12
  },
  "addresses": [
    {
      "address": "10.1.0.62",
      "type": "private"
    },
    {
      "address": "172.24.5.114",
      "type": "public"
    }
  ]
}
`

var createReq = `
{
	"instance": {
		"availability_zone": "us-east1",
		"configuration": "4a78b397-c355-4127-be45-56230b2ab74e",
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
			"size": 2,
			"type": "ssd"
		}
	}
}
`

var instanceWithFault = `
{
  "created": "` + timestamp + `",
  "datastore": {
    "type": "mysql",
    "version": "5.6"
  },
  "flavor": {
    "id": "1",
    "links": [
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "self"
      },
      {
        "href": "https://openstack.example.com/v1.0/1234/flavors/1",
        "rel": "bookmark"
      }
    ]
  },
  "links": [
    {
      "href": "https://openstack.example.com/v1.0/1234/instances/1",
      "rel": "self"
    }
  ],
  "hostname": "e09ad9a3f73309469cf1f43d11e79549caf9acf2.openstack.example.com",
  "id": "{instanceID}",
  "name": "json_rack_instance",
  "status": "BUILD",
  "updated": "` + timestamp + `",
  "volume": {
    "size": 2
  },
  "fault": {
    "message": "some error message",
    "created": "` + timestamp + `",
    "details": "some details about the error"
  }
}
`

var (
	instanceID    = "{instanceID}"
	configGroupID = "00000000-0000-0000-0000-000000000000"
	rootURL       = "/instances"
	resURL        = rootURL + "/" + instanceID
	uRootURL      = resURL + "/root"
	aURL          = resURL + "/action"
)

var (
	restartReq                  = `{"restart": {}}`
	resizeReq                   = `{"resize": {"flavorRef": "2"}}`
	resizeVolReq                = `{"resize": {"volume": {"size": 4}}}`
	attachConfigurationGroupReq = `{"instance": {"configuration": "00000000-0000-0000-0000-000000000000"}}`
	detachConfigurationGroupReq = `{"instance": {}}`
)

var (
	createResp          = fmt.Sprintf(`{"instance": %s}`, instance)
	createWithFaultResp = fmt.Sprintf(`{"instance": %s}`, instanceWithFault)
	listInstancesResp   = fmt.Sprintf(`{"instances":[%s]}`, instance)
	getInstanceResp     = fmt.Sprintf(`{"instance": %s}`, instanceGet)
	enableUserResp      = `{"user":{"name":"root","password":"secretsecret"}}`
	isUserEnabledResp   = `{"rootEnabled":true}`
)

var expectedInstance = instances.Instance{
	Created: timeVal,
	Updated: timeVal,
	Flavor: instances.Flavor{
		ID: "1",
		Links: []gophercloud.Link{
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "self"},
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "bookmark"},
		},
	},
	Hostname: "e09ad9a3f73309469cf1f43d11e79549caf9acf2.openstack.example.com",
	ID:       instanceID,
	Links: []gophercloud.Link{
		{Href: "https://openstack.example.com/v1.0/1234/instances/1", Rel: "self"},
	},
	Name:   "json_rack_instance",
	Status: "BUILD",
	Volume: instances.Volume{Size: 2},
	Datastore: datastores.DatastorePartial{
		Type:    "mysql",
		Version: "5.6",
	},
}

var expectedGetInstance = instances.Instance{
	Created: timeVal,
	Updated: timeVal,
	Flavor: instances.Flavor{
		ID: "1",
		Links: []gophercloud.Link{
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "self"},
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "bookmark"},
		},
	},
	ID: instanceID,
	Links: []gophercloud.Link{
		{Href: "https://openstack.example.com/v1.0/1234/instances/1", Rel: "self"},
	},
	Name:   "test",
	Status: "ACTIVE",
	Volume: instances.Volume{Size: 1, Used: 0.12},
	Datastore: datastores.DatastorePartial{
		Type:    "mysql",
		Version: "5.6",
	},
	Addresses: []instances.Address{
		{Type: "private", Address: "10.1.0.62"},
		{Type: "public", Address: "172.24.5.114"},
	},
}

var expectedInstanceWithFault = instances.Instance{
	Created: timeVal,
	Updated: timeVal,
	Flavor: instances.Flavor{
		ID: "1",
		Links: []gophercloud.Link{
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "self"},
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "bookmark"},
		},
	},
	Hostname: "e09ad9a3f73309469cf1f43d11e79549caf9acf2.openstack.example.com",
	ID:       instanceID,
	Links: []gophercloud.Link{
		{Href: "https://openstack.example.com/v1.0/1234/instances/1", Rel: "self"},
	},
	Name:   "json_rack_instance",
	Status: "BUILD",
	Volume: instances.Volume{Size: 2},
	Datastore: datastores.DatastorePartial{
		Type:    "mysql",
		Version: "5.6",
	},
	Fault: &instances.Fault{
		Created: timeVal,
		Message: "some error message",
		Details: "some details about the error",
	},
}

func HandleCreate(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, rootURL, "POST", createReq, createResp, 200)
}

func HandleCreateWithFault(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, rootURL, "POST", createReq, createWithFaultResp, 200)
}

func HandleList(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, rootURL, "GET", "", listInstancesResp, 200)
}

func HandleGet(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "GET", "", getInstanceResp, 200)
}

func HandleDelete(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "DELETE", "", "", 202)
}

func HandleEnableRoot(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, uRootURL, "POST", "", enableUserResp, 200)
}

func HandleIsRootEnabled(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, uRootURL, "GET", "", isUserEnabledResp, 200)
}

func HandleRestart(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, aURL, "POST", restartReq, "", 202)
}

func HandleResize(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, aURL, "POST", resizeReq, "", 202)
}

func HandleResizeVol(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, aURL, "POST", resizeVolReq, "", 202)
}

func HandleAttachConfigurationGroup(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "PUT", attachConfigurationGroupReq, "", 202)
}

func HandleDetachConfigurationGroup(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "PUT", detachConfigurationGroupReq, "", 202)
}
