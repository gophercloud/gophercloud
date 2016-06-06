package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

var metadata = map[string]string{"gophercloud-test": "containers"}

func TestListContainerInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerInfoSuccessfully(t)

	count := 0
	err := containers.List(fake.ServiceClient(), &containers.ListOpts{Full: true}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := containers.ExtractInfo(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedListInfo, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAllContainerInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerInfoSuccessfully(t)

	allPages, err := containers.List(fake.ServiceClient(), &containers.ListOpts{Full: true}).AllPages()
	th.AssertNoErr(t, err)
	actual, err := containers.ExtractInfo(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedListInfo, actual)
}

func TestListContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerNamesSuccessfully(t)

	count := 0
	err := containers.List(fake.ServiceClient(), &containers.ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := containers.ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract container names: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, ExpectedListNames, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAllContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerNamesSuccessfully(t)

	allPages, err := containers.List(fake.ServiceClient(), &containers.ListOpts{Full: false}).AllPages()
	th.AssertNoErr(t, err)
	actual, err := containers.ExtractNames(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedListNames, actual)
}

func TestCreateContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateContainerSuccessfully(t)

	options := containers.CreateOpts{ContentType: "application/json", Metadata: map[string]string{"foo": "bar"}}
	res := containers.Create(fake.ServiceClient(), "testContainer", options)
	c, err := res.Extract()
	th.CheckNoErr(t, err)
	th.CheckEquals(t, "bar", res.Header["X-Container-Meta-Foo"][0])
	th.CheckEquals(t, "1234567", c.TransID)
}

func TestDeleteContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteContainerSuccessfully(t)

	res := containers.Delete(fake.ServiceClient(), "testContainer")
	th.CheckNoErr(t, res.Err)
}

func TestUpateContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateContainerSuccessfully(t)

	options := &containers.UpdateOpts{Metadata: map[string]string{"foo": "bar"}}
	res := containers.Update(fake.ServiceClient(), "testContainer", options)
	th.CheckNoErr(t, res.Err)
}

func TestGetContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetContainerSuccessfully(t)

	_, err := containers.Get(fake.ServiceClient(), "testContainer").ExtractMetadata()
	th.CheckNoErr(t, err)
}
