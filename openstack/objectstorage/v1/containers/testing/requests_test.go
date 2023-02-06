package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

var (
	metadata = map[string]string{"gophercloud-test": "containers"}
)

func TestContainerNames(t *testing.T) {
	for _, tc := range [...]struct {
		name          string
		containerName string
	}{
		{
			"rejects_a_slash",
			"one/two",
		},
		{
			"rejects_an_escaped_slash",
			"one%2Ftwo",
		},
		{
			"rejects_an_escaped_slash_lowercase",
			"one%2ftwo",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("create", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleCreateContainerSuccessfully(t)

				_, err := containers.Create(fake.ServiceClient(), tc.containerName, nil).Extract()
				th.CheckErr(t, err, &containers.ErrInvalidContainerName{})
			})
			t.Run("delete", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleDeleteContainerSuccessfully(t, WithPath("/"))

				res := containers.Delete(fake.ServiceClient(), tc.containerName)
				th.CheckErr(t, res.Err, &containers.ErrInvalidContainerName{})
			})
			t.Run("update", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleUpdateContainerSuccessfully(t, WithPath("/"))

				contentType := "text/plain"
				options := &containers.UpdateOpts{
					Metadata:         map[string]string{"foo": "bar"},
					ContainerWrite:   new(string),
					ContainerRead:    new(string),
					ContainerSyncTo:  new(string),
					ContainerSyncKey: new(string),
					ContentType:      &contentType,
				}
				res := containers.Update(fake.ServiceClient(), tc.containerName, options)
				th.CheckErr(t, res.Err, &containers.ErrInvalidContainerName{})
			})
			t.Run("get", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleGetContainerSuccessfully(t, WithPath("/"))

				res := containers.Get(fake.ServiceClient(), tc.containerName, nil)
				_, err := res.ExtractMetadata()
				th.CheckErr(t, err, &containers.ErrInvalidContainerName{})

				_, err = res.Extract()
				th.CheckErr(t, err, &containers.ErrInvalidContainerName{})
			})
		})
	}
}

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
	th.CheckEquals(t, 1, count)
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
	th.CheckEquals(t, 1, count)
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

func TestListZeroContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListZeroContainerNames204(t)

	allPages, err := containers.List(fake.ServiceClient(), &containers.ListOpts{Full: false}).AllPages()
	th.AssertNoErr(t, err)
	actual, err := containers.ExtractNames(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, []string{}, actual)
}

func TestCreateContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateContainerSuccessfully(t)

	options := containers.CreateOpts{ContentType: "application/json", Metadata: map[string]string{"foo": "bar"}}
	res := containers.Create(fake.ServiceClient(), "testContainer", options)
	th.CheckEquals(t, "bar", res.Header["X-Container-Meta-Foo"][0])

	expected := &containers.CreateHeader{
		ContentLength: 0,
		ContentType:   "text/html; charset=UTF-8",
		Date:          time.Date(2016, time.August, 17, 19, 25, 43, 0, time.UTC),
		TransID:       "tx554ed59667a64c61866f1-0058b4ba37",
	}
	actual, err := res.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteContainerSuccessfully(t)

	res := containers.Delete(fake.ServiceClient(), "testContainer")
	th.AssertNoErr(t, res.Err)
}

func TestBulkDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleBulkDeleteSuccessfully(t)

	expected := containers.BulkDeleteResponse{
		ResponseStatus: "foo",
		ResponseBody:   "bar",
		NumberDeleted:  2,
		Errors:         [][]string{},
	}

	resp, err := containers.BulkDelete(fake.ServiceClient(), []string{"testContainer1", "testContainer2"}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, *resp)
}

func TestUpdateContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateContainerSuccessfully(t)

	contentType := "text/plain"
	options := &containers.UpdateOpts{
		Metadata:         map[string]string{"foo": "bar"},
		ContainerWrite:   new(string),
		ContainerRead:    new(string),
		ContainerSyncTo:  new(string),
		ContainerSyncKey: new(string),
		ContentType:      &contentType,
	}
	res := containers.Update(fake.ServiceClient(), "testContainer", options)
	th.AssertNoErr(t, res.Err)
}

func TestGetContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetContainerSuccessfully(t)

	getOpts := containers.GetOpts{
		Newest: true,
	}
	res := containers.Get(fake.ServiceClient(), "testContainer", getOpts)
	_, err := res.ExtractMetadata()
	th.AssertNoErr(t, err)

	expected := &containers.GetHeader{
		AcceptRanges:  "bytes",
		BytesUsed:     100,
		ContentType:   "application/json; charset=utf-8",
		Date:          time.Date(2016, time.August, 17, 19, 25, 43, 0, time.UTC),
		ObjectCount:   4,
		Read:          []string{"test"},
		TransID:       "tx554ed59667a64c61866f1-0057b4ba37",
		Write:         []string{"test2", "user4"},
		StoragePolicy: "test_policy",
		Timestamp:     1471298837.95721,
	}
	actual, err := res.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}
