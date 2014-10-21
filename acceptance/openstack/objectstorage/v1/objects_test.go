// +build acceptance

package v1

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

// numObjects is the number of objects to create for testing.
var numObjects = 2

func TestObjects(t *testing.T) {
	// Create a provider client for executing the HTTP request.
	// See common.go for more information.
	client, err := newClient()
	th.AssertNoErr(err)

	// Make a slice of length numObjects to hold the random object names.
	oNames := make([]string, numObjects)
	for i := 0; i < len(oNames); i++ {
		oNames[i] = tools.RandomString("test-object-", 8)
	}

	// Create a container to hold the test objects.
	cName := tools.RandomString("test-container-", 8)
	res = containers.Create(client, cName, nil)
	th.AssertNoErr(res.Err)

	// Defer deletion of the container until after testing.
	defer func() {
		res = containers.Delete(client, cName)
		th.AssertNoErr(res.Err)
	}()

	// Create a slice of buffers to hold the test object content.
	oContents := make([]*bytes.Buffer, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents[i] = bytes.NewBuffer([]byte(tools.RandomString("", 10)))
		res = objects.Create(client, cName, oNames[i], oContents[i], nil)
		th.AssertNoErr(res.Err)
	}
	// Delete the objects after testing.
	defer func() {
		for i := 0; i < numObjects; i++ {
			res = objects.Delete(client, cName, oNames[i], nil)
		}
	}()

	ons := make([]string, 0, len(oNames))
	err = objects.List(client, cName, &objects.ListOpts{Full: false, Prefix: "test-object-"}).EachPage(func(page pagination.Page) (bool, error) {
		names, err := objects.ExtractNames(page)
		th.AssertNoErr(err)
		ons = append(ons, names...)

		return true, nil
	})
	th.AssertNoErr(err)
	if len(ons) != len(oNames) {
		t.Errorf("Expected %d names and got %d", len(oNames), len(ons))
		return
	}

	ois := make([]objects.Object, 0, len(oNames))
	err = objects.List(client, cName, &objects.ListOpts{Full: true, Prefix: "test-object-"}).EachPage(func(page pagination.Page) (bool, error) {
		info, err := objects.ExtractInfo(page)
		th.AssertNoErr(err)

		ois = append(ois, info...)

		return true, nil
	})
	th.AssertNoErr(err)
	if len(ois) != len(oNames) {
		t.Errorf("Expected %d containers and got %d", len(oNames), len(ois))
		return
	}

	// Copy the contents of one object to another.
	res = objects.Copy(client, cName, oNames[0], &objects.CopyOpts{Destination: cName + "/" + oNames[1]})
	th.AssertNoErr(res.Err)

	// Download one of the objects that was created above.
	o1Content, err := objects.Download(client, cName, oNames[0], nil).ExtractContent()
	th.AssertNoErr(err)

	// Download the another object that was create above.
	o2Content, err := objects.Download(client, cName, oNames[1], nil).ExtractContent()
	th.AssertNoErr(err)

	// Compare the two object's contents to test that the copy worked.
	if string(o2Content) != string(o1Content) {
		t.Errorf("Copy failed. Expected\n%s\nand got\n%s", string(o1Content), string(o2Content))
		return
	}

	// Update an object's metadata.
	res = objects.Update(client, cName, oNames[0], &objects.UpdateOpts{Metadata: metadata})
	th.AssertNoErr(res.Err)

	// Delete the object's metadata after testing.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		res = objects.Update(client, cName, oNames[0], &objects.UpdateOpts{Metadata: tempMap})
		th.AssertNoErr(res.Err)
	}()

	// Retrieve an object's metadata.
	om, err := objects.Get(client, cName, oNames[0], nil).ExtractMetadata()
	th.AssertNoErr(err)
	for k := range metadata {
		if om[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
