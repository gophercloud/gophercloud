// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/storage/v1/containers"
)

// numContainers is the number of containers to create for testing.
var numContainers = 2

func TestContainers(t *testing.T) {
	// Create a new client to execute the HTTP requests. See common.go for newClient body.
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("gophercloud-test-container-", 8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		_, err := containers.Create(client, containers.CreateOpts{
			Name: cNames[i],
		})
		if err != nil {
			t.Error(err)
			return
		}
	}
	// Delete the numContainers containers after function completion.
	defer func() {
		for i := 0; i < len(cNames); i++ {
			err = containers.Delete(client, containers.DeleteOpts{
				Name: cNames[i],
			})
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	// List the numContainer names that were just created. To just list those,
	// the 'prefix' parameter is used.
	lr, err := containers.List(client, containers.ListOpts{
		Full: false,
		Params: map[string]string{
			"prefix": "gophercloud-test-container-",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	// Extract the names from the 'List' response.
	cns, err := containers.ExtractNames(lr)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cns) != len(cNames) {
		t.Errorf("Expected %d names and got %d:\nExpected:%v\nActual:%v", len(cNames), len(cns), cNames, cns)
		return
	}

	// List the info for the numContainer containers that were created.
	lr, err = containers.List(client, containers.ListOpts{
		Full: true,
		Params: map[string]string{
			"prefix": "gophercloud-test-container-",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	// Extract the info from the 'List' response.
	cis, err := containers.ExtractInfo(lr)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cis) != len(cNames) {
		t.Errorf("Expected %d containers and got %d", len(cNames), len(cis))
		return
	}

	// Update one of the numContainer container metadata.
	err = containers.Update(client, containers.UpdateOpts{
		Name:     cNames[0],
		Metadata: metadata,
	})
	if err != nil {
		t.Error(err)
		return
	}
	// After the tests are done, delete the metadata that was set.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		err = containers.Update(client, containers.UpdateOpts{
			Name:     cNames[0],
			Metadata: tempMap,
		})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// Retrieve a container's metadata.
	gr, err := containers.Get(client, containers.GetOpts{
		Name: cNames[0],
	})
	if err != nil {
		t.Error(err)
		return
	}
	// Extract the metadata from the 'Get' response.
	cm := containers.ExtractMetadata(gr)
	for k := range metadata {
		if cm[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
