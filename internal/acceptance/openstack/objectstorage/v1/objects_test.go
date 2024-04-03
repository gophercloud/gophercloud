//go:build acceptance || objectstorage || objects

package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// numObjects is the number of objects to create for testing.
var numObjects = 2

func TestObjects(t *testing.T) {
	numObjects := numObjects + 1
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Make a slice of length numObjects to hold the random object names.
	oNames := make([]string, numObjects)
	for i := 0; i < len(oNames)-1; i++ {
		oNames[i] = "test-object-" + tools.RandomFunnyString(8)
	}
	oNames[len(oNames)-1] = "test-object-with-/v1/-in-the-name"

	// Create a container to hold the test objects.
	cName := "test-container-" + tools.RandomFunnyStringNoSlash(8)
	opts := containers.CreateOpts{
		TempURLKey: "super-secret",
	}
	header, err := containers.Create(context.TODO(), client, cName, opts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Create object headers: %+v\n", header)

	// Defer deletion of the container until after testing.
	defer func() {
		res := containers.Delete(context.TODO(), client, cName)
		th.AssertNoErr(t, res.Err)
	}()

	// Create a slice of buffers to hold the test object content.
	oContents := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents[i] = tools.RandomFunnyString(10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oContents[i]),
		}
		res := objects.Create(context.TODO(), client, cName, oNames[i], createOpts)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the objects after testing.
	defer func() {
		for i := 0; i < numObjects; i++ {
			res := objects.Delete(context.TODO(), client, cName, oNames[i], nil)
			th.AssertNoErr(t, res.Err)
		}
	}()

	// List all created objects
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

	// Create temporary URL, download its contents and compare with what was originally created.
	// Downloading the URL validates it (this cannot be done in unit tests).
	objURLs := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		objURLs[i], err = objects.CreateTempURL(context.TODO(), client, cName, oNames[i], objects.CreateTempURLOpts{
			Method: http.MethodGet,
			TTL:    180,
		})
		th.AssertNoErr(t, err)

		resp, err := client.ProviderClient.HTTPClient.Get(objURLs[i])
		th.AssertNoErr(t, err)
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			th.AssertNoErr(t, fmt.Errorf("unexpected response code: %d", resp.StatusCode))
		}

		body, err := io.ReadAll(resp.Body)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, oContents[i], string(body))
		resp.Body.Close()

		// custom Temp URL key with a sha256 digest and exact timestamp
		objURLs[i], err = objects.CreateTempURL(context.TODO(), client, cName, oNames[i], objects.CreateTempURLOpts{
			Method:     http.MethodGet,
			Timestamp:  time.Now().UTC().Add(180 * time.Second),
			Digest:     "sha256",
			TempURLKey: opts.TempURLKey,
		})
		th.AssertNoErr(t, err)

		resp, err = client.ProviderClient.HTTPClient.Get(objURLs[i])
		th.AssertNoErr(t, err)
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			th.AssertNoErr(t, fmt.Errorf("unexpected response code: %d", resp.StatusCode))
		}

		body, err = io.ReadAll(resp.Body)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, oContents[i], string(body))
		resp.Body.Close()
	}

	// Copy the contents of one object to another.
	copyOpts := objects.CopyOpts{
		Destination: "/" + cName + "/" + oNames[1],
	}
	copyres := objects.Copy(context.TODO(), client, cName, oNames[0], copyOpts)
	th.AssertNoErr(t, copyres.Err)

	// Download one of the objects that was created above.
	downloadres := objects.Download(context.TODO(), client, cName, oNames[0], nil)
	th.AssertNoErr(t, downloadres.Err)

	o1Content, err := downloadres.ExtractContent()
	th.AssertNoErr(t, err)

	// Download the another object that was create above.
	downloadOpts := objects.DownloadOpts{
		Newest: true,
	}
	downloadres = objects.Download(context.TODO(), client, cName, oNames[1], downloadOpts)
	th.AssertNoErr(t, downloadres.Err)
	o2Content, err := downloadres.ExtractContent()
	th.AssertNoErr(t, err)

	// Compare the two object's contents to test that the copy worked.
	th.AssertEquals(t, string(o2Content), string(o1Content))

	// Update an object's metadata.
	metadata := map[string]string{
		"Gophercloud-Test": "objects",
	}

	disposition := "inline"
	cType := "text/plain"
	updateOpts := &objects.UpdateOpts{
		Metadata:           metadata,
		ContentDisposition: &disposition,
		ContentType:        &cType,
	}
	updateres := objects.Update(context.TODO(), client, cName, oNames[0], updateOpts)
	th.AssertNoErr(t, updateres.Err)

	// Delete the object's metadata after testing.
	defer func() {
		temp := make([]string, len(metadata))
		i := 0
		for k := range metadata {
			temp[i] = k
			i++
		}
		empty := ""
		cType := "application/octet-stream"
		iTrue := true
		updateOpts = &objects.UpdateOpts{
			RemoveMetadata:     temp,
			ContentDisposition: &empty,
			ContentType:        &cType,
			DetectContentType:  &iTrue,
		}
		res := objects.Update(context.TODO(), client, cName, oNames[0], updateOpts)
		th.AssertNoErr(t, res.Err)

		// Retrieve an object's metadata.
		getOpts := objects.GetOpts{
			Newest: true,
		}
		resp := objects.Get(context.TODO(), client, cName, oNames[0], getOpts)
		om, err := resp.ExtractMetadata()
		th.AssertNoErr(t, err)
		if len(om) > 0 {
			t.Errorf("Expected custom metadata to be empty, found: %v", metadata)
		}
		object, err := resp.Extract()
		th.AssertNoErr(t, err)
		th.AssertEquals(t, empty, object.ContentDisposition)
		th.AssertEquals(t, cType, object.ContentType)
	}()

	// Retrieve an object's metadata.
	getOpts := objects.GetOpts{
		Newest: true,
	}
	resp := objects.Get(context.TODO(), client, cName, oNames[0], getOpts)
	om, err := resp.ExtractMetadata()
	th.AssertNoErr(t, err)
	for k := range metadata {
		if om[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}

	object, err := resp.Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, disposition, object.ContentDisposition)
	th.AssertEquals(t, cType, object.ContentType)
}

func TestObjectsListSubdir(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Create a random subdirectory name.
	cSubdir1 := "test-subdir-" + tools.RandomFunnyStringNoSlash(8)
	cSubdir2 := "test-subdir-" + tools.RandomFunnyStringNoSlash(8)

	// Make a slice of length numObjects to hold the random object names.
	oNames1 := make([]string, numObjects)
	for i := 0; i < len(oNames1); i++ {
		oNames1[i] = cSubdir1 + "/test-object-" + tools.RandomFunnyString(8)
	}

	oNames2 := make([]string, numObjects)
	for i := 0; i < len(oNames2); i++ {
		oNames2[i] = cSubdir2 + "/test-object-" + tools.RandomFunnyString(8)
	}

	// Create a container to hold the test objects.
	cName := "test-container-" + tools.RandomFunnyStringNoSlash(8)
	_, err = containers.Create(context.TODO(), client, cName, nil).Extract()
	th.AssertNoErr(t, err)

	// Defer deletion of the container until after testing.
	defer func() {
		t.Logf("Deleting container %s", cName)
		res := containers.Delete(context.TODO(), client, cName)
		th.AssertNoErr(t, res.Err)
	}()

	// Create a slice of buffers to hold the test object content.
	oContents1 := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents1[i] = tools.RandomFunnyString(10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oContents1[i]),
		}
		res := objects.Create(context.TODO(), client, cName, oNames1[i], createOpts)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the objects after testing.
	defer func() {
		for i := 0; i < numObjects; i++ {
			t.Logf("Deleting object %s", oNames1[i])
			res := objects.Delete(context.TODO(), client, cName, oNames1[i], nil)
			th.AssertNoErr(t, res.Err)
		}
	}()

	oContents2 := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents2[i] = tools.RandomFunnyString(10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oContents2[i]),
		}
		res := objects.Create(context.TODO(), client, cName, oNames2[i], createOpts)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the objects after testing.
	defer func() {
		for i := 0; i < numObjects; i++ {
			t.Logf("Deleting object %s", oNames2[i])
			res := objects.Delete(context.TODO(), client, cName, oNames2[i], nil)
			th.AssertNoErr(t, res.Err)
		}
	}()

	listOpts := objects.ListOpts{
		Delimiter: "/",
	}

	allPages, err := objects.List(client, cName, listOpts).AllPages(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	allObjects, err := objects.ExtractNames(allPages)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v\n", allObjects)
	expected := []string{cSubdir1, cSubdir2}
	for _, e := range expected {
		var valid bool
		for _, a := range allObjects {
			if e+"/" == a {
				valid = true
			}
		}
		if !valid {
			t.Fatalf("could not find %s in results", e)
		}
	}

	listOpts = objects.ListOpts{
		Delimiter: "/",
		Prefix:    cSubdir2,
	}

	allPages, err = objects.List(client, cName, listOpts).AllPages(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	allObjects, err = objects.ExtractNames(allPages)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertEquals(t, allObjects[0], cSubdir2+"/")
	t.Logf("%#v\n", allObjects)
}

func TestObjectsBulkDelete(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Create a random subdirectory name.
	cSubdir1 := "test-subdir-" + tools.RandomFunnyString(8)
	cSubdir2 := "test-subdir-" + tools.RandomFunnyString(8)

	// Make a slice of length numObjects to hold the random object names.
	oNames1 := make([]string, numObjects)
	for i := 0; i < len(oNames1); i++ {
		oNames1[i] = cSubdir1 + "/test-object-" + tools.RandomFunnyString(8)
	}

	oNames2 := make([]string, numObjects)
	for i := 0; i < len(oNames2); i++ {
		oNames2[i] = cSubdir2 + "/test-object-" + tools.RandomFunnyString(8)
	}

	// Create a container to hold the test objects.
	cName := "test-container-" + tools.RandomFunnyStringNoSlash(8)
	_, err = containers.Create(context.TODO(), client, cName, nil).Extract()
	th.AssertNoErr(t, err)

	// Defer deletion of the container until after testing.
	defer func() {
		t.Logf("Deleting container %s", cName)
		res := containers.Delete(context.TODO(), client, cName)
		th.AssertNoErr(t, res.Err)
	}()

	// Create a slice of buffers to hold the test object content.
	oContents1 := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents1[i] = tools.RandomFunnyString(10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oContents1[i]),
		}
		res := objects.Create(context.TODO(), client, cName, oNames1[i], createOpts)
		th.AssertNoErr(t, res.Err)
	}

	oContents2 := make([]string, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents2[i] = tools.RandomFunnyString(10)
		createOpts := objects.CreateOpts{
			Content: strings.NewReader(oContents2[i]),
		}
		res := objects.Create(context.TODO(), client, cName, oNames2[i], createOpts)
		th.AssertNoErr(t, res.Err)
	}

	// Delete the objects after testing.
	expectedResp := objects.BulkDeleteResponse{
		ResponseStatus: "200 OK",
		Errors:         [][]string{},
		NumberDeleted:  numObjects * 2,
	}

	resp, err := objects.BulkDelete(context.TODO(), client, cName, append(oNames1, oNames2...)).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *resp, expectedResp)

	// Verify deletion
	listOpts := objects.ListOpts{
		Delimiter: "/",
	}

	allPages, err := objects.List(client, cName, listOpts).AllPages(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	allObjects, err := objects.ExtractNames(allPages)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertEquals(t, len(allObjects), 0)
}
