package instances

import (
	"testing"

	"github.com/rackspace/gophercloud"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var instanceID = "d4603f69-ec7e-4e9b-803f-600b9205576f"

var expectedInstance = &Instance{
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

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateInstanceSuccessfully(t)

	opts := CreateOpts{
		Name:      "json_rack_instance",
		FlavorRef: "1",
		Databases: db.BatchCreateOpts{
			db.CreateOpts{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			db.CreateOpts{Name: "nextround"},
		},
		Users: UsersOpts{
			UserOpts{
				Name:     "demouser",
				Password: "demopassword",
				Databases: db.BatchCreateOpts{
					db.CreateOpts{Name: "sampledb"},
				},
			},
		},
		Size: 2,
	}

	instance, err := Create(fake.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedInstance, instance)
}

func TestInstanceList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListInstanceSuccessfully(t)

	expectedInstance := Instance{
		Flavor: Flavor{
			ID: "1",
			Links: []gophercloud.Link{
				gophercloud.Link{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "self"},
				gophercloud.Link{Href: "https://openstack.example.com/flavors/1", Rel: "bookmark"},
			},
		},
		ID: "8fb081af-f237-44f5-80cc-b46be1840ca9",
		Links: []gophercloud.Link{
			gophercloud.Link{Href: "https://openstack.example.com/v1.0/1234/instances/8fb081af-f237-44f5-80cc-b46be1840ca9", Rel: "self"},
		},
		Name:   "xml_rack_instance",
		Status: "ACTIVE",
		Volume: Volume{Size: 2},
	}

	pages := 0
	err := List(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractInstances(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 1 {
			t.Fatalf("Expected 1 DB instance, got %d", len(actual))
		}
		th.CheckDeepEquals(t, expectedInstance, actual[0])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestGetInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetInstanceSuccessfully(t, instanceID)

	instance, err := Get(fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, instance, expectedInstance)
}

func TestDeleteInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleDeleteInstanceSuccessfully(t, instanceID)

	res := Delete(fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestEnableRootUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleEnableRootUserSuccessfully(t, instanceID)

	expected := &User{Name: "root", Password: "secretsecret"}

	user, err := EnableRootUser(fake.ServiceClient(), instanceID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, user)
}

func TestIsRootEnabled(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleIsRootEnabledSuccessfully(t, instanceID)

	isEnabled, err := IsRootEnabled(fake.ServiceClient(), instanceID)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, isEnabled)
}

func TestRestartService(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleRestartSuccessfully(t, instanceID)

	res := RestartService(fake.ServiceClient(), instanceID)

	th.AssertNoErr(t, res.Err)
}

func TestResizeInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleResizeInstanceSuccessfully(t, instanceID)

	res := ResizeInstance(fake.ServiceClient(), instanceID, "2")

	th.AssertNoErr(t, res.Err)
}

func TestResizeVolume(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleResizeVolSuccessfully(t, instanceID)

	res := ResizeVolume(fake.ServiceClient(), instanceID, 4)

	th.AssertNoErr(t, res.Err)
}
