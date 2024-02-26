package testing

import (
	"context"
	"testing"

	db "github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreate(t)

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

	instance, err := instances.Create(context.TODO(), fake.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedInstance, instance)
}

func TestCreateWithFault(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateWithFault(t)

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

	instance, err := instances.Create(context.TODO(), fake.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedInstanceWithFault, instance)
}

func TestInstanceList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleList(t)

	pages := 0
	err := instances.List(fake.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGet(t)

	instance, err := instances.Get(context.TODO(), fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedGetInstance, instance)
}

func TestDeleteInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDelete(t)

	res := instances.Delete(context.TODO(), fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestEnableRootUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleEnableRoot(t)

	expected := &users.User{Name: "root", Password: "secretsecret"}
	user, err := instances.EnableRootUser(context.TODO(), fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, user)
}

func TestIsRootEnabled(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleIsRootEnabled(t)

	isEnabled, err := instances.IsRootEnabled(context.TODO(), fake.ServiceClient(), instanceID).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, isEnabled)
}

func TestRestart(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleRestart(t)

	res := instances.Restart(context.TODO(), fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}

func TestResize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleResize(t)

	res := instances.Resize(context.TODO(), fake.ServiceClient(), instanceID, "2")
	th.AssertNoErr(t, res.Err)
}

func TestResizeVolume(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleResizeVol(t)

	res := instances.ResizeVolume(context.TODO(), fake.ServiceClient(), instanceID, 4)
	th.AssertNoErr(t, res.Err)
}

func TestAttachConfigurationGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAttachConfigurationGroup(t)

	res := instances.AttachConfigurationGroup(context.TODO(), fake.ServiceClient(), instanceID, configGroupID)
	th.AssertNoErr(t, res.Err)
}

func TestDetachConfigurationGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDetachConfigurationGroup(t)

	res := instances.DetachConfigurationGroup(context.TODO(), fake.ServiceClient(), instanceID)
	th.AssertNoErr(t, res.Err)
}
