// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/pagination"
)

// numContainers is the number of containers to create for testing.
var numContainers = 2

func TestContainers(t *testing.T) {
	// Create a new client to execute the HTTP requests. See common.go for newClient body.
	client, err := newClient()
	if err != nil {
		t.Error(err)
	}

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("gophercloud-test-container-", 8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		_, err := containers.Create(client, cNames[i], nil).ExtractHeaders()
		if err != nil {
			t.Error(err)
		}
	}
	// Delete the numContainers containers after function completion.
	defer func() {
		for i := 0; i < len(cNames); i++ {
			_, err = containers.Delete(client, cNames[i]).ExtractHeaders()
			if err != nil {
				t.Error(err)
			}
		}
	}()

	// List the numContainer names that were just created. To just list those,
	// the 'prefix' parameter is used.
	err = containers.List(client, &containers.ListOpts{Full: true, Prefix: "gophercloud-test-container-"}).EachPage(func(page pagination.Page) (bool, error) {
		containerList, err := containers.ExtractInfo(page)
		if err != nil {
			t.Error(err)
		}
		for _, n := range containerList {
			t.Logf("Container: Name [%s] Count [%d] Bytes [%d]",
				n.Name, n.Count, n.Bytes)
		}

		return true, nil
	})
	if err != nil {
		t.Error(err)
	}

	// List the info for the numContainer containers that were created.
	err = containers.List(client, &containers.ListOpts{Full: false, Prefix: "gophercloud-test-container-"}).EachPage(func(page pagination.Page) (bool, error) {
		containerList, err := containers.ExtractNames(page)
		if err != nil {
			return false, err
		}
		for _, n := range containerList {
			t.Logf("Container: Name [%s]", n)
		}

		return true, nil
	})
	if err != nil {
		t.Error(err)
	}

	// Update one of the numContainer container metadata.
	_, err = containers.Update(client, cNames[0], &containers.UpdateOpts{Metadata: metadata}).ExtractHeaders()
	if err != nil {
		t.Error(err)
	}
	// After the tests are done, delete the metadata that was set.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		_, err = containers.Update(client, cNames[0], &containers.UpdateOpts{Metadata: tempMap}).ExtractHeaders()
		if err != nil {
			t.Error(err)
		}
	}()

	// Retrieve a container's metadata.
	cm, err := containers.Get(client, cNames[0]).ExtractMetadata()
	if err != nil {
		t.Error(err)
	}
	for k := range metadata {
		if cm[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
		}
	}
}
