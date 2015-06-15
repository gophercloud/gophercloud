package v2

// TODO
// compare with openstack/compute/v2/servers/requests_test.go

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fakeclient "github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageCreationSuccessfully(t)

	id := "e7db3b45-8db7-47ad-8109-3fb55c2c24fd"
	name := "Ubuntu 12.10"

	actualImage, err := Create(fakeclient.ServiceClient(), CreateOpts{
		Id: &id,
		Name: &name,
		Tags: []string{"ubuntu", "quantal"},
	}).Extract()

	th.AssertNoErr(t, err)

	container_format := "bare"
	disk_format := "qcow2"
	owner := "b4eedccc6fb74fa8a7ad6b08382b852b"
	min_disk_gigabytes := 0
	min_ram_megabytes := 0

	expectedImage := Image{
		Id: "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Name: "Ubuntu 12.10",
		Tags: []string{"ubuntu", "quantal"},
		
		Status: ImageStatusQueued,
		
		ContainerFormat: &container_format,
		DiskFormat: &disk_format,

		MinDiskGigabytes: &min_disk_gigabytes,
		MinRamMegabytes: &min_ram_megabytes,
		
		Owner: &owner,

		Visibility: ImageVisibilityPrivate,

		Metadata: make(map[string]string),
		Properties: make(map[string]string),
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestGetImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageGetSuccessfully(t)

	actualImage, err := Get(fakeclient.ServiceClient(), "1bea47ed-f6a9-463b-b423-14b9cca9ad27").Extract()

	th.AssertNoErr(t, err)

	checksum := "64d7c1cd2b6f60c92c14662941cb7913"
	size_bytes := 13167616
	container_format := "bare"
	disk_format := "qcow2"
	min_disk_gigabytes := 0
	min_ram_megabytes := 0
	owner := "5ef70662f8b34079a6eddb8da9d75fe8"

	expectedImage := Image{
		Id: "1bea47ed-f6a9-463b-b423-14b9cca9ad27",
		Name: "cirros-0.3.2-x86_64-disk",
		Tags: []string{},

		Status: ImageStatusActive,

		ContainerFormat: &container_format,
		DiskFormat: &disk_format,

		MinDiskGigabytes: &min_disk_gigabytes,
		MinRamMegabytes: &min_ram_megabytes,

		Owner: &owner,

		Protected: false,
		Visibility: ImageVisibilityPublic,

		Checksum: &checksum,
		SizeBytes: &size_bytes,

		Metadata: make(map[string]string),
		Properties: make(map[string]string),
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestDeleteImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageDeleteSuccessfully(t)

	Delete(fakeclient.ServiceClient(), "1bea47ed-f6a9-463b-b423-14b9cca9ad27")
	// TODO
}

func TestUpdateImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageUpdateSuccessfully(t)

	actualImage, err := Update(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea", UpdateOpts{
		ReplaceImageName{NewName: "Fedora 17"},
		ReplaceImageTags{NewTags: []string{"fedora", "beefy"}, },
	}).Extract()

	th.AssertNoErr(t, err)

	sizebytes := 2254249
	checksum := "2cec138d7dae2aa59038ef8c9aec2390"

	expectedImage := Image{
		Id: "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		Name: "Fedora 17",
		Status: ImageStatusActive,
		Visibility: ImageVisibilityPublic,

		SizeBytes: &sizebytes,
		Checksum: &checksum,

		Tags: []string{
			"fedora",
			"beefy",
		},

		Owner: nil,
		MinRamMegabytes: nil,
		MinDiskGigabytes: nil,

		DiskFormat: nil,
		ContainerFormat: nil,

		Metadata: make(map[string]string),
		Properties: make(map[string]string),
	}
	
	th.AssertDeepEquals(t, &expectedImage, actualImage)
}
