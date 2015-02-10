package instances

import (
	"testing"

	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/db/v1/instances"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateInstanceSuccessfully(t)

	opts := CreateOpts{
		Name:      "json_rack_instance",
		FlavorRef: "1",
		Databases: os.DatabasesOpts{
			os.DatabaseOpts{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			os.DatabaseOpts{Name: "nextround"},
		},
		Users: os.UsersOpts{
			os.UserOpts{
				Name:     "demouser",
				Password: "demopassword",
				Databases: os.DatabasesOpts{
					os.DatabaseOpts{Name: "sampledb"},
				},
			},
		},
		Size:         2,
		RestorePoint: "1234567890",
	}

	instance, err := Create(fake.ServiceClient(), opts).Extract()

	expected := &Instance{
		Created:   "2014-02-13T21:47:13",
		Updated:   "2014-02-13T21:47:13",
		Datastore: Datastore{Type: "mysql", Version: "5.6"},
		Flavor: os.Flavor{
			ID: "1",
			Links: []gophercloud.Link{
				gophercloud.Link{Href: "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1", Rel: "self"},
				gophercloud.Link{Href: "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1", Rel: "bookmark"},
			},
		},
		Hostname: "e09ad9a3f73309469cf1f43d11e79549caf9acf2.rackspaceclouddb.com",
		ID:       "d4603f69-ec7e-4e9b-803f-600b9205576f",
		Links: []gophercloud.Link{
			gophercloud.Link{Href: "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1", Rel: "self"},
		},
		Name:   "json_rack_instance",
		Status: "BUILD",
		Volume: os.Volume{Size: 2},
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, instance)
}
