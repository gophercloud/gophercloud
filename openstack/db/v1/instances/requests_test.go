package instances

import (
	"testing"

	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/gophercloud/testhelper/fixture"
)

var (
	instanceID = "{instanceID}"
	rootURL    = "/instances"
	resURL     = rootURL + "/" + instanceID
	uRootURL   = resURL + "/root"
	aURL       = resURL + "/action"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, rootURL, "POST", createReq, createResp, 200)

	opts := CreateOpts{
		Name:      "json_rack_instance",
		FlavorRef: "1",
		Databases: db.BatchCreateOpts{
			db.CreateOpts{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			db.CreateOpts{Name: "nextround"},
		},
		Users: users.BatchCreateOpts{
			users.CreateOpts{
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
	th.AssertDeepEquals(t, &expectedInstance, instance)
}

func TestInstanceList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, rootURL, "GET", "", listInstancesResp, 200)

	pages := 0
	err := List(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractInstances(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, []Instance{expectedInstance}, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, resURL, "GET", "", getInstanceResp, 200)

	instance, err := Get(fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedInstance, instance)
}

func TestDeleteInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, resURL, "DELETE", "", "", 202)

	res := Delete(fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestEnableRootUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, uRootURL, "POST", "", enableUserResp, 200)

	expected := &users.User{Name: "root", Password: "secretsecret"}
	user, err := EnableRootUser(fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, user)
}

func TestIsRootEnabled(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, uRootURL, "GET", "", isUserEnabledResp, 200)

	isEnabled, err := IsRootEnabled(fake.ServiceClient(), instanceID)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, isEnabled)
}

func TestRestartService(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, aURL, "POST", restartReq, "", 202)

	res := RestartService(fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestResizeInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, aURL, "POST", resizeReq, "", 202)

	res := ResizeInstance(fake.ServiceClient(), instanceID, "2")
	th.AssertNoErr(t, res.Err)
}

func TestResizeVolume(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, aURL, "POST", resizeVolReq, "", 202)

	res := ResizeVolume(fake.ServiceClient(), instanceID, 4)
	th.AssertNoErr(t, res.Err)
}
