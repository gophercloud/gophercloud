// +build acceptance imageservice

package v2

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/acceptance/tools"
	images "github.com/rackspace/gophercloud/openstack/imageservice/v2"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestListImages(t *testing.T) {
	client := newClient(t)

	t.Logf("Id\tName\tOwner\tChecksum\tSizeBytes")

	pager := images.List(client, nil)
	count, pages := 0, 0
	pager.EachPage(func(page pagination.Page) (bool, error) {
		pages++
		t.Logf("---")

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

	t.Logf("--------\n%d images listed on %d pages.\n", count, pages)
}

func TestListImagesFilter(t *testing.T) {
	client := newClient(t)
	t.Logf("Id\tName\tOwner\tChecksum\tSizeBytes")

	pager := images.List(client, images.ListOpts{Limit: 1})
	count, pages := 0, 0
	pager.EachPage(func(page pagination.Page) (bool, error) {
		pages++
		t.Logf("---")

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

	t.Logf("--------\n%d images listed on %d pages.\n", count, pages)

}

func TestCreateDeleteImage(t *testing.T) {
	client := newClient(t)
	imageName := tools.RandomString("ACCPT", 16)
	containerFormat := "ami"
	createResult := images.Create(client, images.CreateOpts{Name: &imageName,
		ContainerFormat: &containerFormat,
		DiskFormat:      &containerFormat})

	th.AssertNoErr(t, createResult.Err)
	image, err := createResult.Extract()
	th.AssertNoErr(t, err)

	t.Logf("Image %v", image)

	image, err = images.Get(client, image.ID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, image.Status, images.ImageStatusQueued)

	deleteResult := images.Delete(client, image.ID)
	th.AssertNoErr(t, deleteResult.Err)
}

func TestUploadDownloadImage(t *testing.T) {
	client := newClient(t)

	//creating image
	imageName := tools.RandomString("ACCPT", 16)
	containerFormat := "ami"
	createResult := images.Create(client, images.CreateOpts{Name: &imageName,
		ContainerFormat: &containerFormat,
		DiskFormat:      &containerFormat})
	th.AssertNoErr(t, createResult.Err)
	image, err := createResult.Extract()
	th.AssertNoErr(t, err)
	t.Logf("Image %v", image)

	//checking status
	image, err = images.Get(client, image.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, image.Status, images.ImageStatusQueued)

	//uploading image data
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	putImageResult := images.PutImageData(client, image.ID, bytes.NewReader(data))
	th.AssertNoErr(t, putImageResult.Err)

	//checking status
	image, err = images.Get(client, image.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, image.Status, images.ImageStatusActive)
	th.AssertEquals(t, *image.SizeBytes, 9)

	//downloading image data
	reader, err := images.GetImageData(client, image.ID).Extract()
	th.AssertNoErr(t, err)
	receivedData, err := ioutil.ReadAll(reader)
	t.Logf("Received data %v", receivedData)
	th.AssertNoErr(t, err)
	th.AssertByteArrayEquals(t, data, receivedData)

	//deteting image
	deleteResult := images.Delete(client, image.ID)
	th.AssertNoErr(t, deleteResult.Err)

}

func TestUpdateImage(t *testing.T) {
	client := newClient(t)

	//creating image
	image := createTestImage(t, client)

	t.Logf("Image tags %v", image.Tags)

	tags := []string{"acceptance-testing"}
	updatedImage, err := images.Update(client, image.ID, images.UpdateOpts{
		images.ReplaceImageTags{
			NewTags: tags}}).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Received tags '%v'", tags)
	th.AssertDeepEquals(t, updatedImage.Tags, tags)
}

func TestImageMemberCreateListDelete(t *testing.T) {
	client := newClient(t)

	//creating image
	image := createTestImage(t, client)
	defer deleteImage(t, client, image)

	//creating member
	member, err := images.CreateMember(client, image.ID, "tenant").Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, member)

	//listing member
	var members *[]images.ImageMember
	members, err = images.ListMembers(client, image.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, members)
	th.AssertEquals(t, 1, len(*members))

	t.Logf("Members after adding one %v", members)

	//checking just created member
	m := (*members)[0]
	th.AssertEquals(t, "pending", m.Status)
	th.AssertEquals(t, "tenant", m.MemberID)

	//deleting member
	deleteResult := images.DeleteMember(client, image.ID, "tenant")
	th.AssertNoErr(t, deleteResult.Err)

	//listing member
	members, err = images.ListMembers(client, image.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, members)
	th.AssertEquals(t, 0, len(*members))

	t.Logf("Members after deleting one %v", members)
}

func TestImageMemberDetailsAndUpdate(t *testing.T) {
	// getting current tenant id
	memberTenantID := os.Getenv("OS_TENANT_ID")
	if memberTenantID == "" {
		t.Fatalf("Please define OS_TENANT_ID for image member updating test was '%s'", memberTenantID)
	}

	client := newClient(t)

	//creating image
	image := createTestImage(t, client)
	defer deleteImage(t, client, image)

	//creating member
	member, err := images.CreateMember(client, image.ID, memberTenantID).Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, member)

	//checking image member details
	member, err = images.ShowMemberDetails(client, image.ID, memberTenantID).Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, member)

	th.AssertEquals(t, memberTenantID, member.MemberID)
	th.AssertEquals(t, "pending", member.Status)

	t.Logf("Updating image's %s member status for tenant %s to 'accepted' ", image.ID, memberTenantID)

	//updating image
	member, err = images.UpdateMember(client, image.ID, memberTenantID, "accepted").Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, member)
	th.AssertEquals(t, "accepted", member.Status)

}

func createTestImage(t *testing.T, client *gophercloud.ServiceClient) images.Image {
	//creating image
	imageName := tools.RandomString("ACCPT", 16)
	containerFormat := "ami"
	createResult := images.Create(client, images.CreateOpts{Name: &imageName,
		ContainerFormat: &containerFormat,
		DiskFormat:      &containerFormat})
	th.AssertNoErr(t, createResult.Err)
	image, err := createResult.Extract()
	th.AssertNoErr(t, err)
	t.Logf("Image %v", image)

	//checking status
	image, err = images.Get(client, image.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, image.Status, images.ImageStatusQueued)
	return *image
}

func deleteImage(t *testing.T, client *gophercloud.ServiceClient, image images.Image) {
	//deteting image
	deleteResult := images.Delete(client, image.ID)
	th.AssertNoErr(t, deleteResult.Err)
}
