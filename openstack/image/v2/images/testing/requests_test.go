package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageListSuccessfully(t, fakeServer)

	t.Logf("Id\tName\tOwner\tChecksum\tSizeBytes")

	pager := images.List(client.ServiceClient(fakeServer), images.ListOpts{Limit: 1})
	t.Logf("Pager state %v", pager)
	count, pages := 0, 0
	err := pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		t.Logf("Page %v", page)
		images, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}

		for _, i := range images {
			t.Logf("%s\t%s\t%s\t%s\t%v\t\n", i.ID, i.Name, i.Owner, i.Checksum, i.SizeBytes)
			count++
		}

		return true, nil
	})
	th.AssertNoErr(t, err)

	t.Logf("--------\n%d images listed on %d pages.\n", count, pages)
	th.AssertEquals(t, 3, pages)
	th.AssertEquals(t, 3, count)
}

func TestAllPagesImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageListSuccessfully(t, fakeServer)

	pages, err := images.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	images, err := images.ExtractImages(pages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 3, len(images))
}

func TestCreateImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageCreationSuccessfully(t, fakeServer)

	id := "e7db3b45-8db7-47ad-8109-3fb55c2c24fd"
	name := "Ubuntu 12.10"

	actualImage, err := images.Create(context.TODO(), client.ServiceClient(fakeServer), images.CreateOpts{
		ID:   id,
		Name: name,
		Properties: map[string]string{
			"architecture": "x86_64",
		},
		Tags: []string{"ubuntu", "quantal"},
	}).Extract()

	th.AssertNoErr(t, err)

	containerFormat := "bare"
	diskFormat := "qcow2"
	owner := "b4eedccc6fb74fa8a7ad6b08382b852b"
	minDiskGigabytes := 0
	minRAMMegabytes := 0
	file := actualImage.File
	createdDate := actualImage.CreatedAt
	lastUpdate := actualImage.UpdatedAt
	schema := "/v2/schemas/image"

	expectedImage := images.Image{
		ID:   "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Name: "Ubuntu 12.10",
		Tags: []string{"ubuntu", "quantal"},

		Status: images.ImageStatusQueued,

		ContainerFormat: containerFormat,
		DiskFormat:      diskFormat,

		MinDiskGigabytes: minDiskGigabytes,
		MinRAMMegabytes:  minRAMMegabytes,

		Owner: owner,

		Visibility:  images.ImageVisibilityPrivate,
		File:        file,
		CreatedAt:   createdDate,
		UpdatedAt:   lastUpdate,
		Schema:      schema,
		VirtualSize: 0,
		Properties: map[string]any{
			"hw_disk_bus":       "scsi",
			"hw_disk_bus_model": "virtio-scsi",
			"hw_scsi_model":     "virtio-scsi",
		},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestCreateImageNulls(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageCreationSuccessfullyNulls(t, fakeServer)

	id := "e7db3b45-8db7-47ad-8109-3fb55c2c24fd"
	name := "Ubuntu 12.10"

	actualImage, err := images.Create(context.TODO(), client.ServiceClient(fakeServer), images.CreateOpts{
		ID:   id,
		Name: name,
		Tags: []string{"ubuntu", "quantal"},
		Properties: map[string]string{
			"architecture": "x86_64",
		},
	}).Extract()

	th.AssertNoErr(t, err)

	containerFormat := "bare"
	diskFormat := "qcow2"
	owner := "b4eedccc6fb74fa8a7ad6b08382b852b"
	minDiskGigabytes := 0
	minRAMMegabytes := 0
	file := actualImage.File
	createdDate := actualImage.CreatedAt
	lastUpdate := actualImage.UpdatedAt
	schema := "/v2/schemas/image"
	properties := map[string]any{
		"architecture": "x86_64",
	}
	sizeBytes := int64(0)

	expectedImage := images.Image{
		ID:   "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Name: "Ubuntu 12.10",
		Tags: []string{"ubuntu", "quantal"},

		Status: images.ImageStatusQueued,

		ContainerFormat: containerFormat,
		DiskFormat:      diskFormat,

		MinDiskGigabytes: minDiskGigabytes,
		MinRAMMegabytes:  minRAMMegabytes,

		Owner: owner,

		Visibility: images.ImageVisibilityPrivate,
		File:       file,
		CreatedAt:  createdDate,
		UpdatedAt:  lastUpdate,
		Schema:     schema,
		Properties: properties,
		SizeBytes:  sizeBytes,
		OpenStackImageImportMethods: []string{
			"glance-direct",
			"web-download",
		},
		OpenStackImageStoreIDs: []string{
			"123",
			"456",
		},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestGetImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageGetSuccessfully(t, fakeServer)

	actualImage, err := images.Get(context.TODO(), client.ServiceClient(fakeServer), "1bea47ed-f6a9-463b-b423-14b9cca9ad27").Extract()

	th.AssertNoErr(t, err)

	checksum := "64d7c1cd2b6f60c92c14662941cb7913"
	sizeBytes := int64(13167616)
	containerFormat := "bare"
	diskFormat := "qcow2"
	minDiskGigabytes := 0
	minRAMMegabytes := 0
	owner := "5ef70662f8b34079a6eddb8da9d75fe8"
	file := actualImage.File
	createdDate := actualImage.CreatedAt
	lastUpdate := actualImage.UpdatedAt
	schema := "/v2/schemas/image"

	expectedImage := images.Image{
		ID:   "1bea47ed-f6a9-463b-b423-14b9cca9ad27",
		Name: "cirros-0.3.2-x86_64-disk",
		Tags: []string{},

		Status: images.ImageStatusActive,

		ContainerFormat: containerFormat,
		DiskFormat:      diskFormat,

		MinDiskGigabytes: minDiskGigabytes,
		MinRAMMegabytes:  minRAMMegabytes,

		Owner: owner,

		Protected:  false,
		Visibility: images.ImageVisibilityPublic,
		Hidden:     false,

		Checksum:    checksum,
		SizeBytes:   sizeBytes,
		File:        file,
		CreatedAt:   createdDate,
		UpdatedAt:   lastUpdate,
		Schema:      schema,
		VirtualSize: 0,
		Properties: map[string]any{
			"hw_disk_bus":       "scsi",
			"hw_disk_bus_model": "virtio-scsi",
			"hw_scsi_model":     "virtio-scsi",
		},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestDeleteImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageDeleteSuccessfully(t, fakeServer)

	result := images.Delete(context.TODO(), client.ServiceClient(fakeServer), "1bea47ed-f6a9-463b-b423-14b9cca9ad27")
	th.AssertNoErr(t, result.Err)
}

func TestUpdateImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageUpdateSuccessfully(t, fakeServer)

	actualImage, err := images.Update(context.TODO(), client.ServiceClient(fakeServer), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea", images.UpdateOpts{
		images.ReplaceImageName{NewName: "Fedora 17"},
		images.ReplaceImageTags{NewTags: []string{"fedora", "beefy"}},
		images.ReplaceImageMinDisk{NewMinDisk: 21},
		images.ReplaceImageMinRam{NewMinRam: 1024},
		images.ReplaceImageHidden{NewHidden: false},
		images.ReplaceImageProtected{NewProtected: true},
		images.UpdateImageProperty{
			Op:    images.AddOp,
			Name:  "empty_value",
			Value: "",
		},
	}).Extract()

	th.AssertNoErr(t, err)

	sizebytes := int64(2254249)
	checksum := "2cec138d7dae2aa59038ef8c9aec2390"
	file := actualImage.File
	createdDate := actualImage.CreatedAt
	lastUpdate := actualImage.UpdatedAt
	schema := "/v2/schemas/image"

	expectedImage := images.Image{
		ID:         "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		Name:       "Fedora 17",
		Status:     images.ImageStatusActive,
		Visibility: images.ImageVisibilityPublic,
		Hidden:     false,
		Protected:  true,

		SizeBytes: sizebytes,
		Checksum:  checksum,

		Tags: []string{
			"fedora",
			"beefy",
		},

		Owner:            "",
		MinRAMMegabytes:  1024,
		MinDiskGigabytes: 21,

		DiskFormat:      "",
		ContainerFormat: "",
		File:            file,
		CreatedAt:       createdDate,
		UpdatedAt:       lastUpdate,
		Schema:          schema,
		VirtualSize:     0,
		Properties: map[string]any{
			"hw_disk_bus":       "scsi",
			"hw_disk_bus_model": "virtio-scsi",
			"hw_scsi_model":     "virtio-scsi",
			"empty_value":       "",
		},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestImageDateQuery(t *testing.T) {
	date := time.Date(2014, 1, 1, 1, 1, 1, 0, time.UTC)

	listOpts := images.ListOpts{
		CreatedAtQuery: &images.ImageDateQuery{
			Date:   date,
			Filter: images.FilterGTE,
		},
		UpdatedAtQuery: &images.ImageDateQuery{
			Date: date,
		},
	}

	expectedQueryString := "?created_at=gte%3A2014-01-01T01%3A01%3A01Z&updated_at=2014-01-01T01%3A01%3A01Z"
	actualQueryString, err := listOpts.ToImageListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expectedQueryString, actualQueryString)
}

func TestImageListByTags(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageListByTagsSuccessfully(t, fakeServer)

	listOpts := images.ListOpts{
		Tags: []string{"foo", "bar"},
	}

	expectedQueryString := "?tag=foo&tag=bar"
	actualQueryString, err := listOpts.ToImageListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expectedQueryString, actualQueryString)

	pages, err := images.List(client.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allImages, err := images.ExtractImages(pages)
	th.AssertNoErr(t, err)

	checksum := "64d7c1cd2b6f60c92c14662941cb7913"
	sizeBytes := int64(13167616)
	containerFormat := "bare"
	diskFormat := "qcow2"
	minDiskGigabytes := 0
	minRAMMegabytes := 0
	owner := "5ef70662f8b34079a6eddb8da9d75fe8"
	file := allImages[0].File
	createdDate := allImages[0].CreatedAt
	lastUpdate := allImages[0].UpdatedAt
	schema := "/v2/schemas/image"
	tags := []string{"foo", "bar"}

	expectedImage := images.Image{
		ID:   "1bea47ed-f6a9-463b-b423-14b9cca9ad27",
		Name: "cirros-0.3.2-x86_64-disk",
		Tags: tags,

		Status: images.ImageStatusActive,

		ContainerFormat: containerFormat,
		DiskFormat:      diskFormat,

		MinDiskGigabytes: minDiskGigabytes,
		MinRAMMegabytes:  minRAMMegabytes,

		Owner: owner,

		Protected:  false,
		Visibility: images.ImageVisibilityPublic,

		Checksum:    checksum,
		SizeBytes:   sizeBytes,
		File:        file,
		CreatedAt:   createdDate,
		UpdatedAt:   lastUpdate,
		Schema:      schema,
		VirtualSize: 0,
		Properties: map[string]any{
			"hw_disk_bus":       "scsi",
			"hw_disk_bus_model": "virtio-scsi",
			"hw_scsi_model":     "virtio-scsi",
		},
	}

	th.AssertDeepEquals(t, expectedImage, allImages[0])
}

func TestUpdateImageProperties(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleImageUpdatePropertiesSuccessfully(t, fakeServer)

	actualImage, err := images.Update(context.TODO(), client.ServiceClient(fakeServer), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea", images.UpdateOpts{
		images.UpdateImageProperty{
			Op:    images.AddOp,
			Name:  "hw_disk_bus",
			Value: "scsi",
		},
		images.UpdateImageProperty{
			Op:    images.AddOp,
			Name:  "hw_disk_bus_model",
			Value: "virtio-scsi",
		},
		images.UpdateImageProperty{
			Op:    images.AddOp,
			Name:  "hw_scsi_model",
			Value: "virtio-scsi",
		},
	}).Extract()

	th.AssertNoErr(t, err)

	sizebytes := int64(2254249)
	checksum := "2cec138d7dae2aa59038ef8c9aec2390"
	file := actualImage.File
	createdDate := actualImage.CreatedAt
	lastUpdate := actualImage.UpdatedAt
	schema := "/v2/schemas/image"

	expectedImage := images.Image{
		ID:         "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		Name:       "Fedora 17",
		Status:     images.ImageStatusActive,
		Visibility: images.ImageVisibilityPublic,

		SizeBytes: sizebytes,
		Checksum:  checksum,

		Tags: []string{
			"fedora",
			"beefy",
		},

		Owner:            "",
		MinRAMMegabytes:  0,
		MinDiskGigabytes: 0,

		DiskFormat:      "",
		ContainerFormat: "",
		File:            file,
		CreatedAt:       createdDate,
		UpdatedAt:       lastUpdate,
		Schema:          schema,
		VirtualSize:     0,
		Properties: map[string]any{
			"hw_disk_bus":       "scsi",
			"hw_disk_bus_model": "virtio-scsi",
			"hw_scsi_model":     "virtio-scsi",
		},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}
