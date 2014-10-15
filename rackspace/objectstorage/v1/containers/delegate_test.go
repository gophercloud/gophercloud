package containers

import (
	"testing"

	os "github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestListContainerInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleListContainerInfoSuccessfully(t)

	count := 0
	err := List(fake.ServiceClient(), &os.ListOpts{Full: true}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractInfo(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, os.ExpectedListInfo, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleListContainerNamesSuccessfully(t)

	count := 0
	err := List(fake.ServiceClient(), &os.ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract container names: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, os.ExpectedListNames, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestCreateContainers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleCreateContainerSuccessfully(t)

	options := os.CreateOpts{ContentType: "application/json", Metadata: map[string]string{"foo": "bar"}}
	headers, err := Create(fake.ServiceClient(), "testContainer", options).ExtractHeaders()
	th.CheckNoErr(t, err)
	th.CheckEquals(t, "bar", headers["X-Container-Meta-Foo"][0])

}

func TestDeleteContainers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleDeleteContainerSuccessfully(t)

	_, err := Delete(fake.ServiceClient(), "testContainer").ExtractHeaders()
	th.CheckNoErr(t, err)
}

func TestUpdateContainers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleUpdateContainerSuccessfully(t)

	options := &os.UpdateOpts{Metadata: map[string]string{"foo": "bar"}}
	_, err := Update(fake.ServiceClient(), "testContainer", options).ExtractHeaders()
	th.CheckNoErr(t, err)
}

func TestGetContainers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleGetContainerSuccessfully(t)

	_, err := Get(fake.ServiceClient(), "testContainer").ExtractMetadata()
	th.CheckNoErr(t, err)
}
