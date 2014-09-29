// +build acceptance

package v1

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/storage/v1/containers"
	"github.com/rackspace/gophercloud/openstack/storage/v1/objects"
	"github.com/rackspace/gophercloud/pagination"
)

// numObjects is the number of objects to create for testing.
var numObjects = 2

func TestObjects(t *testing.T) {
	// Create a provider client for executing the HTTP request.
	// See common.go for more information.
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	// Make a slice of length numObjects to hold the random object names.
	oNames := make([]string, numObjects)
	for i := 0; i < len(oNames); i++ {
		oNames[i] = tools.RandomString("test-object-", 8)
	}

	// Create a container to hold the test objects.
	cName := tools.RandomString("test-container-", 8)
	_, err = containers.Create(client, cName, containers.CreateOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	// Defer deletion of the container until after testing.
	defer func() {
		err = containers.Delete(client, cName)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// Create a slice of buffers to hold the test object content.
	oContents := make([]*bytes.Buffer, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents[i] = bytes.NewBuffer([]byte(tools.RandomString("", 10)))
		err = objects.Create(client, cName, oNames[i], oContents[i], objects.CreateOpts{})
		if err != nil {
			t.Error(err)
			return
		}
	}
	// Delete the objects after testing.
	defer func() {
		for i := 0; i < numObjects; i++ {
			err = objects.Delete(client, cName, oNames[i], objects.DeleteOpts{})
		}
	}()

	pager := objects.List(client, cName, objects.ListOpts{Full: false})
	ons := make([]string, 0, len(oNames))
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		names, err := objects.ExtractNames(page)
		if err != nil {
			return false, err
		}
		ons = append(ons, names...)

		return true, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if len(ons) != len(oNames) {
		t.Errorf("Expected %d names and got %d", len(oNames), len(ons))
		return
	}

	pager = objects.List(client, cName, objects.ListOpts{Full: true})
	ois := make([]objects.Object, 0, len(oNames))
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		info, err := objects.ExtractInfo(page)
		if err != nil {
			return false, nil
		}

		ois = append(ois, info...)

		return true, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if len(ois) != len(oNames) {
		t.Errorf("Expected %d containers and got %d", len(oNames), len(ois))
		return
	}

	// Copy the contents of one object to another.
	err = objects.Copy(client, cName, oNames[0], objects.CopyOpts{Destination: cName + "/" + oNames[1]})
	if err != nil {
		t.Error(err)
		return
	}

	// Download one of the objects that was created above.
	o1Content, err := objects.Download(client, cName, oNames[0], objects.DownloadOpts{}).ExtractContent()
	if err != nil {
		t.Error(err)
		return
	}
	// Download the another object that was create above.
	o2Content, err := objects.Download(client, cName, oNames[1], objects.DownloadOpts{}).ExtractContent()
	if err != nil {
		t.Error(err)
		return
	}
	// Compare the two object's contents to test that the copy worked.
	if string(o2Content) != string(o1Content) {
		t.Errorf("Copy failed. Expected\n%s\nand got\n%s", string(o1Content), string(o2Content))
		return
	}

	// Update an object's metadata.
	err = objects.Update(client, cName, oNames[0], objects.UpdateOpts{Metadata: metadata})
	if err != nil {
		t.Error(err)
		return
	}
	// Delete the object's metadata after testing.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		err = objects.Update(client, cName, oNames[0], objects.UpdateOpts{Metadata: tempMap})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// Retrieve an object's metadata.
	om, err := objects.Get(client, cName, oNames[0], objects.GetOpts{}).ExtractMetadata()
	if err != nil {
		t.Error(err)
		return
	}
	for k := range metadata {
		if om[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
