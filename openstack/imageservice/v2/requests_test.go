package v2

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fakeclient "github.com/rackspace/gophercloud/testhelper/client"
)

func TestListImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageListSuccessfully(t)

	t.Logf("Test setup %+v\n", th.Server)

	t.Logf("Id\tName\tOwner\tChecksum\tSizeBytes")

	pager := List(fakeclient.ServiceClient(), ListOpts{Limit: 1})
	t.Logf("Pager state %v", pager)
	count, pages := 0, 0
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		pages++
		t.Logf("Page %v", page)
		images, err := ExtractImages(page)
		if err != nil {
			return false, err
		}

		for _, i := range images {
			t.Logf("%s\t%s\t%s\t%s\t%v\t\n", i.ID, i.Name, *(i.Owner), *(i.Checksum), i.SizeBytes)
			count++
		}

		return true, nil
	})
	th.AssertNoErr(t, err)

	t.Logf("--------\n%d images listed on %d pages.\n", count, pages)
	th.AssertEquals(t, 3, pages)
	th.AssertEquals(t, 3, count)
}

func TestCreateImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageCreationSuccessfully(t)

	id := "e7db3b45-8db7-47ad-8109-3fb55c2c24fd"
	name := "Ubuntu 12.10"

	actualImage, err := Create(fakeclient.ServiceClient(), CreateOpts{
		ID:   &id,
		Name: &name,
		Tags: []string{"ubuntu", "quantal"},
	}).Extract()

	th.AssertNoErr(t, err)

	containerFormat := "bare"
	diskFormat := "qcow2"
	owner := "b4eedccc6fb74fa8a7ad6b08382b852b"
	minDiskGigabytes := 0
	minRAMMegabytes := 0

	expectedImage := Image{
		ID:   "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Name: "Ubuntu 12.10",
		Tags: []string{"ubuntu", "quantal"},

		Status: ImageStatusQueued,

		ContainerFormat: &containerFormat,
		DiskFormat:      &diskFormat,

		MinDiskGigabytes: &minDiskGigabytes,
		MinRAMMegabytes:  &minRAMMegabytes,

		Owner: &owner,

		Visibility: ImageVisibilityPrivate,

		Metadata:   make(map[string]string),
		Properties: make(map[string]string),
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestCreateImageNulls(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageCreationSuccessfullyNulls(t)

	id := "e7db3b45-8db7-47ad-8109-3fb55c2c24fd"
	name := "Ubuntu 12.10"

	actualImage, err := Create(fakeclient.ServiceClient(), CreateOpts{
		ID:   &id,
		Name: &name,
		Tags: []string{"ubuntu", "quantal"},
	}).Extract()

	th.AssertNoErr(t, err)

	containerFormat := "bare"
	diskFormat := "qcow2"
	owner := "b4eedccc6fb74fa8a7ad6b08382b852b"
	minDiskGigabytes := 0
	minRAMMegabytes := 0

	expectedImage := Image{
		ID:   "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Name: "Ubuntu 12.10",
		Tags: []string{"ubuntu", "quantal"},

		Status: ImageStatusQueued,

		ContainerFormat: &containerFormat,
		DiskFormat:      &diskFormat,

		MinDiskGigabytes: &minDiskGigabytes,
		MinRAMMegabytes:  &minRAMMegabytes,

		Owner: &owner,

		Visibility: ImageVisibilityPrivate,

		Metadata:   make(map[string]string),
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
	sizeBytes := 13167616
	containerFormat := "bare"
	diskFormat := "qcow2"
	minDiskGigabytes := 0
	minRAMMegabytes := 0
	owner := "5ef70662f8b34079a6eddb8da9d75fe8"

	expectedImage := Image{
		ID:   "1bea47ed-f6a9-463b-b423-14b9cca9ad27",
		Name: "cirros-0.3.2-x86_64-disk",
		Tags: []string{},

		Status: ImageStatusActive,

		ContainerFormat: &containerFormat,
		DiskFormat:      &diskFormat,

		MinDiskGigabytes: &minDiskGigabytes,
		MinRAMMegabytes:  &minRAMMegabytes,

		Owner: &owner,

		Protected:  false,
		Visibility: ImageVisibilityPublic,

		Checksum:  &checksum,
		SizeBytes: &sizeBytes,

		Metadata:   make(map[string]string),
		Properties: make(map[string]string),
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestDeleteImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageDeleteSuccessfully(t)

	result := Delete(fakeclient.ServiceClient(), "1bea47ed-f6a9-463b-b423-14b9cca9ad27")
	th.AssertNoErr(t, result.Err)
}

func TestUpdateImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageUpdateSuccessfully(t)

	actualImage, err := Update(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea", UpdateOpts{
		ReplaceImageName{NewName: "Fedora 17"},
		ReplaceImageTags{NewTags: []string{"fedora", "beefy"}},
	}).Extract()

	th.AssertNoErr(t, err)

	sizebytes := 2254249
	checksum := "2cec138d7dae2aa59038ef8c9aec2390"

	expectedImage := Image{
		ID:         "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		Name:       "Fedora 17",
		Status:     ImageStatusActive,
		Visibility: ImageVisibilityPublic,

		SizeBytes: &sizebytes,
		Checksum:  &checksum,

		Tags: []string{
			"fedora",
			"beefy",
		},

		Owner:            nil,
		MinRAMMegabytes:  nil,
		MinDiskGigabytes: nil,

		DiskFormat:      nil,
		ContainerFormat: nil,

		Metadata:   make(map[string]string),
		Properties: make(map[string]string),
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}

func TestPutImageData(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandlePutImageDataSuccessfully(t)

	PutImageData(
		fakeclient.ServiceClient(),
		"da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		readSeekerOfBytes([]byte{5, 3, 7, 24}))

	// TODO
}

func readSeekerOfBytes(bs []byte) io.ReadSeeker {
	return &RS{bs: bs}
}

// implements io.ReadSeeker
type RS struct {
	bs     []byte
	offset int
}

func (rs *RS) Read(p []byte) (int, error) {
	leftToRead := len(rs.bs) - rs.offset

	if 0 < leftToRead {
		bytesToWrite := min(leftToRead, len(p))
		for i := 0; i < bytesToWrite; i++ {
			p[i] = rs.bs[rs.offset]
			rs.offset++
		}
		return bytesToWrite, nil
	}
	return 0, io.EOF
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (rs *RS) Seek(offset int64, whence int) (int64, error) {
	var offsetInt = int(offset)
	if whence == 0 {
		rs.offset = offsetInt
	} else if whence == 1 {
		rs.offset = rs.offset + offsetInt
	} else if whence == 2 {
		rs.offset = len(rs.bs) - offsetInt
	} else {
		return 0, fmt.Errorf("For parameter `whence`, expected value in {0,1,2} but got: %#v", whence)
	}

	return int64(rs.offset), nil
}

func TestGetImageData(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetImageDataSuccessfully(t)

	rdr, err := GetImageData(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea").Extract()

	th.AssertNoErr(t, err)

	bs, err := ioutil.ReadAll(rdr)

	th.AssertNoErr(t, err)

	th.AssertByteArrayEquals(t, []byte{34, 87, 0, 23, 23, 23, 56, 255, 254, 0}, bs)
}

func TestCreateMemberSuccessfully(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateImageMemberSuccessfully(t)
	im, err := CreateMember(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ImageMember{
		CreatedAt: "2013-09-20T19:22:19Z",
		ImageID:   "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		MemberID:  "8989447062e04a818baf9e073fd04fa7",
		Schema:    "/v2/schemas/member",
		Status:    "pending",
		UpdatedAt: "2013-09-20T19:25:31Z",
	}, *im)

}

func TestCreateMemberMemberConflict(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateImageMemberConflict(t)

	result := CreateMember(fakeclient.ServiceClient(), "da3b75d9-memberConflict",
		"8989447062e04a818baf9e073fd04fa7")

	if result.Err == nil {
		t.Fatalf("Expected error in result defined (Err: %v)", result.Err)
	}

	message := result.Err.Error()
	if !strings.Contains(message, "is already member for image") {
		t.Fatalf("Wrong error message: %s", message)
	}

}
func TestCreateMemberInvalidVisibility(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateImageMemberInvalidVisibility(t)

	result := CreateMember(fakeclient.ServiceClient(), "da3b75d9-invalid-visibility",
		"8989447062e04a818baf9e073fd04fa7")

	if result.Err == nil {
		t.Fatalf("Expected error in result defined (Err: %v)", result.Err)
	}

	message := result.Err.Error()
	if !strings.Contains(message, "which 'visibility' attribute is private") {
		t.Fatalf("Wrong error message: %s", message)
	}
}

func TestMemberListSuccessfully(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageMemberList(t)

	images, err := ListMembers(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea").Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, images)
	th.AssertEquals(t, 2, len(*images))
}

func TestMemberListEmpty(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageMemberEmptyList(t)

	images, err := ListMembers(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea").Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, images)
	th.AssertEquals(t, 0, len(*images))
}

func TestShowMemberDetails(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageMemberDetails(t)
	md, err := ShowMemberDetails(fakeclient.ServiceClient(),
		"da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7").Extract()

	th.AssertNoErr(t, err)
	th.AssertNotNil(t, md)

	th.AssertDeepEquals(t, ImageMember{
		CreatedAt: "2013-11-26T07:21:21Z",
		ImageID:   "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		MemberID:  "8989447062e04a818baf9e073fd04fa7",
		Schema:    "/v2/schemas/member",
		Status:    "pending",
		UpdatedAt: "2013-11-26T07:21:21Z",
	}, *md)
}

func TestDeleteMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	counter := HandleImageMemberDeleteSuccessfully(t)

	result := DeleteMember(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7")
	th.AssertEquals(t, 1, counter.Counter)
	th.AssertNoErr(t, result.Err)
}

func TestDeleteMemberByNonOwner(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	counter := HandleImageMemberDeleteByNonOwner(t)

	result := DeleteMember(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7")
	th.AssertEquals(t, 1, counter.Counter)

	if result.Err == nil {
		t.Fatalf("Expected error in result defined (Err: %v)", result.Err)
	}

	message := result.Err.Error()
	if !strings.Contains(message, "You must be the owner of the specified image") {
		t.Fatalf("Wrong error message: %s", message)
	}
}

func TestMemberUpdateSuccessfully(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	counter := HandleImageMemberUpdate(t)
	im, err := UpdateMember(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7", "accepted").Extract()
	th.AssertEquals(t, 1, counter.Counter)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ImageMember{
		CreatedAt: "2013-11-26T07:21:21Z",
		ImageID:   "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		MemberID:  "8989447062e04a818baf9e073fd04fa7",
		Schema:    "/v2/schemas/member",
		Status:    "accepted",
		UpdatedAt: "2013-11-26T07:21:21Z",
	}, *im)

}
