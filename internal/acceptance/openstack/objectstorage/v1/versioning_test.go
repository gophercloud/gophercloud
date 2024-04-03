//go:build acceptance || objectstorage || versioning

package v1

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestObjectsVersioning(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ussuri")

	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Make a slice of length numObjects to hold the random object names.
	oNames := make([]string, numObjects)
	for i := 0; i < len(oNames); i++ {
		oNames[i] = tools.RandomString("test-object-", 8)
	}

	// Create a container to hold the test objects.
	cName := tools.RandomString("test-container-", 8)
	opts := containers.CreateOpts{
		VersionsEnabled: true,
	}
	header, err := containers.Create(context.TODO(), client, cName, opts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Create container headers: %+v\n", header)

	// Defer deletion of the container until after testing.
	defer func() {
		res := containers.Delete(context.TODO(), client, cName)
		th.AssertNoErr(t, res.Err)
	}()

	// ensure versioning is enabled
	get, err := containers.Get(context.TODO(), client, cName, nil).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Get container headers: %+v\n", get)
	th.AssertEquals(t, true, get.VersionsEnabled)

	// Create a slice of buffers to hold the test object content.
	oContents := make([]string, numObjects)
	oContentVersionIDs := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents[i] = tools.RandomString("", 10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oContents[i]),
		}
		obj, err := objects.Create(context.TODO(), client, cName, oNames[i], createOpts).Extract()
		th.AssertNoErr(t, err)
		oContentVersionIDs[i] = obj.ObjectVersionID
	}
	oNewContents := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oNewContents[i] = tools.RandomString("", 10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oNewContents[i]),
		}
		_, err := objects.Create(context.TODO(), client, cName, oNames[i], createOpts).Extract()
		th.AssertNoErr(t, err)
	}
	// Delete the objects after testing two times.
	defer func() {
		// disable object versioning
		opts := containers.UpdateOpts{
			VersionsEnabled: new(bool),
		}
		header, err := containers.Update(context.TODO(), client, cName, opts).Extract()
		th.AssertNoErr(t, err)

		t.Logf("Update container headers: %+v\n", header)

		// ensure versioning is disabled
		get, err := containers.Get(context.TODO(), client, cName, nil).Extract()
		th.AssertNoErr(t, err)
		t.Logf("Get container headers: %+v\n", get)
		th.AssertEquals(t, false, get.VersionsEnabled)

		// delete all object versions before deleting the container
		currentVersionIDs := make([]string, numObjects)
		for i := 0; i < numObjects; i++ {
			opts := objects.DeleteOpts{
				ObjectVersionID: oContentVersionIDs[i],
			}
			obj, err := objects.Delete(context.TODO(), client, cName, oNames[i], opts).Extract()
			th.AssertNoErr(t, err)
			currentVersionIDs[i] = obj.ObjectCurrentVersionID
		}
		for i := 0; i < numObjects; i++ {
			opts := objects.DeleteOpts{
				ObjectVersionID: currentVersionIDs[i],
			}
			res := objects.Delete(context.TODO(), client, cName, oNames[i], opts)
			th.AssertNoErr(t, res.Err)
		}
	}()

	// List created objects
	listOpts := objects.ListOpts{
		Prefix: "test-object-",
	}

	allPages, err := objects.List(client, cName, listOpts).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list objects: %v", err)
	}

	ons, err := objects.ExtractNames(allPages)
	if err != nil {
		t.Fatalf("Unable to extract objects: %v", err)
	}
	th.AssertEquals(t, len(ons), len(oNames))

	ois, err := objects.ExtractInfo(allPages)
	if err != nil {
		t.Fatalf("Unable to extract object info: %v", err)
	}
	th.AssertEquals(t, len(ois), len(oNames))

	// List all created objects
	listOpts = objects.ListOpts{
		Prefix:   "test-object-",
		Versions: true,
	}

	allPages, err = objects.List(client, cName, listOpts).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list objects: %v", err)
	}

	ons, err = objects.ExtractNames(allPages)
	if err != nil {
		t.Fatalf("Unable to extract objects: %v", err)
	}
	th.AssertEquals(t, len(ons), 2*len(oNames))

	ois, err = objects.ExtractInfo(allPages)
	if err != nil {
		t.Fatalf("Unable to extract object info: %v", err)
	}
	th.AssertEquals(t, len(ois), 2*len(oNames))

	// ensure proper versioning attributes are set
	for i, obj := range ois {
		if i%2 == 0 {
			th.AssertEquals(t, true, obj.IsLatest)
		} else {
			th.AssertEquals(t, false, obj.IsLatest)
		}
		if obj.VersionID == "" {
			t.Fatalf("Unexpected empty version_id for the %s object", obj.Name)
		}
	}

	// Download one of the objects that was created above.
	downloadres := objects.Download(context.TODO(), client, cName, oNames[0], nil)
	th.AssertNoErr(t, downloadres.Err)

	o1Content, err := downloadres.ExtractContent()
	th.AssertNoErr(t, err)

	// Compare the two object's contents to test that the copy worked.
	th.AssertEquals(t, oNewContents[0], string(o1Content))

	// Download the another object that was create above.
	downloadOpts := objects.DownloadOpts{
		Newest: true,
	}
	downloadres = objects.Download(context.TODO(), client, cName, oNames[1], downloadOpts)
	th.AssertNoErr(t, downloadres.Err)
	o2Content, err := downloadres.ExtractContent()
	th.AssertNoErr(t, err)

	// Compare the two object's contents to test that the copy worked.
	th.AssertEquals(t, oNewContents[1], string(o2Content))
}
