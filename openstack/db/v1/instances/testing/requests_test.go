package testing

import (
	"context"
	"testing"

	db "github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreate(t, fakeServer)

	opts := instances.CreateOpts{
		AvailabilityZone: "us-east1",
		Configuration:    "4a78b397-c355-4127-be45-56230b2ab74e",
		Name:             "json_rack_instance",
		FlavorRef:        "1",
		Databases: db.BatchCreateOpts{
			{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			{Name: "nextround"},
		},
		Users: users.BatchCreateOpts{
			{
				Name:     "demouser",
				Password: "demopassword",
				Databases: db.BatchCreateOpts{
					{Name: "sampledb"},
				},
			},
		},
		Size:       2,
		VolumeType: "ssd",
	}

	instance, err := instances.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedInstance, instance)
}

func TestCreateWithFault(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateWithFault(t, fakeServer)

	opts := instances.CreateOpts{
		AvailabilityZone: "us-east1",
		Configuration:    "4a78b397-c355-4127-be45-56230b2ab74e",
		Name:             "json_rack_instance",
		FlavorRef:        "1",
		Databases: db.BatchCreateOpts{
			{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			{Name: "nextround"},
		},
		Users: users.BatchCreateOpts{
			{
				Name:     "demouser",
				Password: "demopassword",
				Databases: db.BatchCreateOpts{
					{Name: "sampledb"},
				},
			},
		},
		Size:       2,
		VolumeType: "ssd",
	}

	instance, err := instances.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedInstanceWithFault, instance)
}

func TestInstanceList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleList(t, fakeServer)

	pages := 0
	err := instances.List(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := instances.ExtractInstances(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, []instances.Instance{expectedInstance}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetInstance(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGet(t, fakeServer)

	instance, err := instances.Get(context.TODO(), client.ServiceClient(fakeServer), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedGetInstance, instance)
}

func TestDeleteInstance(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDelete(t, fakeServer)

	res := instances.Delete(context.TODO(), client.ServiceClient(fakeServer), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestEnableRootUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleEnableRoot(t, fakeServer)

	expected := &users.User{Name: "root", Password: "secretsecret"}
	user, err := instances.EnableRootUser(context.TODO(), client.ServiceClient(fakeServer), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, user)
}

func TestIsRootEnabled(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleIsRootEnabled(t, fakeServer)

	isEnabled, err := instances.IsRootEnabled(context.TODO(), client.ServiceClient(fakeServer), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, isEnabled)
}

func TestRestart(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRestart(t, fakeServer)

	res := instances.Restart(context.TODO(), client.ServiceClient(fakeServer), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestResize(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleResize(t, fakeServer)

	res := instances.Resize(context.TODO(), client.ServiceClient(fakeServer), instanceID, "2")
	th.AssertNoErr(t, res.Err)
}

func TestResizeVolume(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleResizeVol(t, fakeServer)

	res := instances.ResizeVolume(context.TODO(), client.ServiceClient(fakeServer), instanceID, 4)
	th.AssertNoErr(t, res.Err)
}

func TestAttachConfigurationGroup(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAttachConfigurationGroup(t, fakeServer)

	res := instances.AttachConfigurationGroup(context.TODO(), client.ServiceClient(fakeServer), instanceID, configGroupID)
	th.AssertNoErr(t, res.Err)
}

func TestDetachConfigurationGroup(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDetachConfigurationGroup(t, fakeServer)

	res := instances.DetachConfigurationGroup(context.TODO(), client.ServiceClient(fakeServer), instanceID)
	th.AssertNoErr(t, res.Err)
}
