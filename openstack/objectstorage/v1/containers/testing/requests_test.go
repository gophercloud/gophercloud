package testing

import (
	"context"
	"testing"
	"time"

	v1 "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestContainerNames(t *testing.T) {
	for _, tc := range [...]struct {
		name          string
		containerName string
		expectedError error
	}{
		{
			"rejects_a_slash",
			"one/two",
			v1.ErrInvalidContainerName{},
		},
		{
			"rejects_an_empty_string",
			"",
			v1.ErrEmptyContainerName{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("create", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleCreateContainerSuccessfully(t)

				_, err := containers.Create(context.TODO(), fake.ServiceClient(), tc.containerName, nil).Extract()
				th.CheckErr(t, err, &tc.expectedError)
			})
			t.Run("delete", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleDeleteContainerSuccessfully(t, WithPath("/"))

				res := containers.Delete(context.TODO(), fake.ServiceClient(), tc.containerName)
				th.CheckErr(t, res.Err, &tc.expectedError)
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
				res := containers.Update(context.TODO(), fake.ServiceClient(), tc.containerName, options)
				th.CheckErr(t, res.Err, &tc.expectedError)
			})
			t.Run("get", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleGetContainerSuccessfully(t, WithPath("/"))

				res := containers.Get(context.TODO(), fake.ServiceClient(), tc.containerName, nil)
				_, err := res.ExtractMetadata()
				th.CheckErr(t, err, &tc.expectedError)

				_, err = res.Extract()
				th.CheckErr(t, err, &tc.expectedError)
			})
		})
	}
}

func TestListContainerInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerInfoSuccessfully(t)

	count := 0
	err := containers.List(fake.ServiceClient(), &containers.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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

	allPages, err := containers.List(fake.ServiceClient(), &containers.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := containers.ExtractInfo(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedListInfo, actual)
}

func TestListContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListContainerInfoSuccessfully(t)

	count := 0
	err := containers.List(fake.ServiceClient(), &containers.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	HandleListContainerInfoSuccessfully(t)

	allPages, err := containers.List(fake.ServiceClient(), &containers.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := containers.ExtractNames(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedListNames, actual)
}

func TestListZeroContainerNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListZeroContainerNames204(t)

	allPages, err := containers.List(fake.ServiceClient(), &containers.ListOpts{}).AllPages(context.TODO())
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
	res := containers.Create(context.TODO(), fake.ServiceClient(), "testContainer", options)
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

	res := containers.Delete(context.TODO(), fake.ServiceClient(), "testContainer")
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

	resp, err := containers.BulkDelete(context.TODO(), fake.ServiceClient(), []string{"testContainer1", "testContainer2"}).Extract()
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
	res := containers.Update(context.TODO(), fake.ServiceClient(), "testContainer", options)
	th.AssertNoErr(t, res.Err)
}

func TestGetContainer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetContainerSuccessfully(t)

	getOpts := containers.GetOpts{
		Newest: true,
	}
	res := containers.Get(context.TODO(), fake.ServiceClient(), "testContainer", getOpts)
	_, err := res.ExtractMetadata()
	th.AssertNoErr(t, err)

	expected := &containers.GetHeader{
		AcceptRanges:    "bytes",
		BytesUsed:       100,
		ContentType:     "application/json; charset=utf-8",
		Date:            time.Date(2016, time.August, 17, 19, 25, 43, 0, time.UTC),
		ObjectCount:     4,
		Read:            []string{"test"},
		TransID:         "tx554ed59667a64c61866f1-0057b4ba37",
		Write:           []string{"test2", "user4"},
		StoragePolicy:   "test_policy",
		Timestamp:       1471298837.95721,
		VersionsEnabled: true,
	}
	actual, err := res.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateContainerVersioningOff(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateContainerVersioningOff(t)

	contentType := "text/plain"
	options := &containers.UpdateOpts{
		Metadata:         map[string]string{"foo": "bar"},
		ContainerWrite:   new(string),
		ContainerRead:    new(string),
		ContainerSyncTo:  new(string),
		ContainerSyncKey: new(string),
		ContentType:      &contentType,
		VersionsEnabled:  new(bool),
	}
	_, err := containers.Update(context.TODO(), fake.ServiceClient(), "testVersioning", options).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdateContainerVersioningOn(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateContainerVersioningOn(t)

	iTrue := true
	contentType := "text/plain"
	options := &containers.UpdateOpts{
		Metadata:         map[string]string{"foo": "bar"},
		ContainerWrite:   new(string),
		ContainerRead:    new(string),
		ContainerSyncTo:  new(string),
		ContainerSyncKey: new(string),
		ContentType:      &contentType,
		VersionsEnabled:  &iTrue,
	}
	_, err := containers.Update(context.TODO(), fake.ServiceClient(), "testVersioning", options).Extract()
	th.AssertNoErr(t, err)
}
