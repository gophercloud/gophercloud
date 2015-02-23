package instances

import (
	"testing"

	"github.com/rackspace/gophercloud"
	osDBs "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/openstack/db/v1/flavors"
	os "github.com/rackspace/gophercloud/openstack/db/v1/instances"
	osUsers "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/rackspace/db/v1/datastores"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/gophercloud/testhelper/fixture"
)

var (
	instanceID = "{instanceID}"
	_rootURL   = "/instances"
	resURL     = "/instances/" + instanceID
)

var expectedInstance = &Instance{
	Created:   "2014-02-13T21:47:13",
	Updated:   "2014-02-13T21:47:13",
	Datastore: datastores.DatastorePartial{Type: "mysql", Version: "5.6"},
	Flavor: flavors.Flavor{
		ID: "1",
		Links: []gophercloud.Link{
			gophercloud.Link{Href: "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1", Rel: "self"},
			gophercloud.Link{Href: "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1", Rel: "bookmark"},
		},
	},
	Hostname: "e09ad9a3f73309469cf1f43d11e79549caf9acf2.rackspaceclouddb.com",
	ID:       instanceID,
	Links: []gophercloud.Link{
		gophercloud.Link{Href: "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1", Rel: "self"},
	},
	Name:   "json_rack_instance",
	Status: "BUILD",
	Volume: os.Volume{Size: 2},
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, _rootURL, "POST", createReq, createResp, 200)

	opts := CreateOpts{
		Name:      "json_rack_instance",
		FlavorRef: "1",
		Databases: osDBs.BatchCreateOpts{
			osDBs.CreateOpts{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			osDBs.CreateOpts{Name: "nextround"},
		},
		Users: osUsers.BatchCreateOpts{
			osUsers.CreateOpts{
				Name:     "demouser",
				Password: "demopassword",
				Databases: osDBs.BatchCreateOpts{
					osDBs.CreateOpts{Name: "sampledb"},
				},
			},
		},
		Size:         2,
		RestorePoint: "1234567890",
	}

	instance, err := Create(fake.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedInstance, instance)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, resURL, "GET", "", getResp, 200)

	instance, err := Get(fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedInstance, instance)
}

func TestDeleteInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleDelete(t)

	res := Delete(fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestEnableRootUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleEnableRoot(t)

	expected := &osUsers.User{Name: "root", Password: "secretsecret"}

	user, err := EnableRootUser(fake.ServiceClient(), instanceID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, user)
}

func TestRestartService(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleRestart(t)

	res := RestartService(fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestResizeInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleResize(t)

	res := ResizeInstance(fake.ServiceClient(), instanceID, "2")
	th.AssertNoErr(t, res.Err)
}

func TestResizeVolume(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleResizeVol(t)

	res := ResizeVolume(fake.ServiceClient(), instanceID, 4)
	th.AssertNoErr(t, res.Err)
}
