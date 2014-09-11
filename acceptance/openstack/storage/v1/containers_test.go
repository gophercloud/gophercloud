// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud/acceptance/tools"
	"github.com/rackspace/gophercloud/openstack/storage/v1/containers"
)

var numContainers = 2

func TestContainers(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("gophercloud-test-container-", 8)
	}

	for i := 0; i < len(cNames); i++ {
		_, err := containers.Create(client, containers.CreateOpts{
			Name: cNames[i],
		})
		if err != nil {
			t.Error(err)
			return
		}
	}
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
	cns, err := containers.ExtractNames(lr)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cns) != len(cNames) {
		t.Errorf("Expected %d names and got %d:\nExpected:%v\nActual:%v", len(cNames), len(cns), cNames, cns)
		return
	}

	lr, err = containers.List(client, containers.ListOpts{
		Full: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	cis, err := containers.ExtractInfo(lr)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cis) != len(cNames) {
		t.Errorf("Expected %d containers and got %d", len(cNames), len(cis))
		return
	}
	err = containers.Update(client, containers.UpdateOpts{
		Name:     cNames[0],
		Metadata: metadata,
	})
	if err != nil {
		t.Error(err)
		return
	}
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

	gr, err := containers.Get(client, containers.GetOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	cm := containers.ExtractMetadata(gr)
	for k := range metadata {
		if cm[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
