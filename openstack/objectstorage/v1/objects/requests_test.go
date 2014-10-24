package objects

import (
	"bytes"
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestDownloadObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDownloadObjectSuccessfully(t)

	content, err := Download(fake.ServiceClient(), "testContainer", "testObject", nil).ExtractContent()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, string(content), "Successful download with Gophercloud")
}

func TestListObjectInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListObjectsInfoSuccessfully(t)

	count := 0
	options := &ListOpts{Full: true}
	err := List(fake.ServiceClient(), "testContainer", options).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractInfo(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedListInfo, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListObjectNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListObjectNamesSuccessfully(t)

	count := 0
	options := &ListOpts{Full: false}
	err := List(fake.ServiceClient(), "testContainer", options).EachPage(func(page pagination.Page) (bool, error) {
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

func TestCreateObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateObjectSuccessfully(t)

	content := bytes.NewBufferString("Did gyre and gimble in the wabe")
	options := &CreateOpts{ContentType: "application/json"}
	res := Create(fake.ServiceClient(), "testContainer", "testObject", content, options)
	th.AssertNoErr(t, res.Err)
}

func TestCopyObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCopyObjectSuccessfully(t)

	options := &CopyOpts{Destination: "/newTestContainer/newTestObject"}
	res := Copy(fake.ServiceClient(), "testContainer", "testObject", options)
	th.AssertNoErr(t, res.Err)
}

func TestDeleteObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteObjectSuccessfully(t)

	res := Delete(fake.ServiceClient(), "testContainer", "testObject", nil)
	th.AssertNoErr(t, res.Err)
}

func TestUpateObjectMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateObjectSuccessfully(t)

	options := &UpdateOpts{Metadata: map[string]string{"Gophercloud-Test": "objects"}}
	res := Update(fake.ServiceClient(), "testContainer", "testObject", options)
	th.AssertNoErr(t, res.Err)
}

func TestGetObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetObjectSuccessfully(t)

	expected := map[string]string{"Gophercloud-Test": "objects"}
	actual, err := Get(fake.ServiceClient(), "testContainer", "testObject", nil).ExtractMetadata()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}
