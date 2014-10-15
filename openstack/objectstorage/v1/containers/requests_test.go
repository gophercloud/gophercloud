package containers

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var metadata = map[string]string{"gophercloud-test": "containers"}

func TestListContainerInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerInfoSuccessfully(t)

	count := 0
	err := List(fake.ServiceClient(), &ListOpts{Full: true}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractInfo(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedListInfo, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerNamesSuccessfully(t)

	count := 0
	err := List(fake.ServiceClient(), &ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractNames(page)
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

func TestCreateContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateContainerSuccessfully(t)

	options := CreateOpts{ContentType: "application/json", Metadata: map[string]string{"foo": "bar"}}
	headers, err := Create(fake.ServiceClient(), "testContainer", options).ExtractHeaders()
	th.CheckNoErr(t, err)
	th.CheckEquals(t, "bar", headers["X-Container-Meta-Foo"][0])
}

func TestDeleteContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteContainerSuccessfully(t)

	_, err := Delete(fake.ServiceClient(), "testContainer").ExtractHeaders()
	th.CheckNoErr(t, err)
}

func TestUpateContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateContainerSuccessfully(t)

	options := &UpdateOpts{Metadata: map[string]string{"foo": "bar"}}
	_, err := Update(fake.ServiceClient(), "testContainer", options).ExtractHeaders()
	th.CheckNoErr(t, err)
}

func TestGetContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetContainerSuccessfully(t)

	_, err := Get(fake.ServiceClient(), "testContainer").ExtractMetadata()
	th.CheckNoErr(t, err)
}
