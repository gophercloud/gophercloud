// +build acceptance

package openstack

import (
	"bytes"
	"github.com/rackspace/gophercloud/acceptance/tools"
	storage "github.com/rackspace/gophercloud/openstack/storage/v1"
	"github.com/rackspace/gophercloud/openstack/storage/v1/accounts"
	"github.com/rackspace/gophercloud/openstack/storage/v1/containers"
	"github.com/rackspace/gophercloud/openstack/storage/v1/objects"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"strings"
	"testing"
)

var metadata = map[string]string{"gopher": "cloud"}
var numContainers = 2
var numObjects = 2

func newClient() (*storage.Client, error) {
	ao, err := utils.AuthOptions()
	if err != nil {
		return nil, err
	}

	client, err := utils.NewClient(ao, utils.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
		Type:   "object-store",
	})
	if err != nil {
		return nil, err
	}

	return storage.NewClient(client.Endpoint, client.Authority, client.Options), nil
}

func TestAccount(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	err = accounts.Update(client, accounts.UpdateOpts{
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
		err = accounts.Update(client, accounts.UpdateOpts{
			Metadata: tempMap,
		})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	gr, err := accounts.Get(client, accounts.GetOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	am := accounts.ExtractMetadata(gr)
	for k := range metadata {
		if am[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}

func TestContainers(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("test-container-", 8)
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
		t.Errorf("Expected %d names and got %d", len(cNames), len(cns))
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

func TestObjects(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Error(err)
		return
	}

	oNames := make([]string, numObjects)
	for i := 0; i < len(oNames); i++ {
		oNames[i] = tools.RandomString("test-object-", 8)
	}

	cName := tools.RandomString("test-container-", 8)
	_, err = containers.Create(client, containers.CreateOpts{
		Name: cName,
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err = containers.Delete(client, containers.DeleteOpts{
			Name: cName,
		})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	oContents := make([]*bytes.Buffer, numObjects)
	for i := 0; i < numObjects; i++ {
		oContents[i] = bytes.NewBuffer([]byte(tools.RandomString("", 10)))
		err = objects.Create(client, objects.CreateOpts{
			Container: cName,
			Name:      oNames[i],
			Content:   oContents[i],
		})
		if err != nil {
			t.Error(err)
			return
		}
	}
	defer func() {
		for i := 0; i < numObjects; i++ {
			err = objects.Delete(client, objects.DeleteOpts{
				Container: cName,
				Name:      oNames[i],
			})
		}
	}()

	lr, err := objects.List(client, objects.ListOpts{
		Full:      false,
		Container: cName,
	})
	if err != nil {
		t.Error(err)
		return
	}
	ons, err := objects.ExtractNames(lr)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ons) != len(oNames) {
		t.Errorf("Expected %d names and got %d", len(oNames), len(ons))
		return
	}

	lr, err = objects.List(client, objects.ListOpts{
		Full:      true,
		Container: cName,
	})
	if err != nil {
		t.Error(err)
		return
	}
	ois, err := objects.ExtractInfo(lr)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ois) != len(oNames) {
		t.Errorf("Expected %d containers and got %d", len(oNames), len(ois))
		return
	}

	err = objects.Copy(client, objects.CopyOpts{
		Container:    cName,
		Name:         oNames[0],
		NewContainer: cName,
		NewName:      oNames[1],
	})
	if err != nil {
		t.Error(err)
		return
	}

	dr, err := objects.Download(client, objects.DownloadOpts{
		Container: cName,
		Name:      oNames[1],
	})
	if err != nil {
		t.Error(err)
		return
	}
	o2Content, err := objects.ExtractContent(dr)
	if err != nil {
		t.Error(err)
	}
	dr, err = objects.Download(client, objects.DownloadOpts{
		Container: cName,
		Name:      oNames[0],
	})
	if err != nil {
		t.Error(err)
		return
	}
	o1Content, err := objects.ExtractContent(dr)
	if err != nil {
		t.Error(err)
		return
	}
	if string(o2Content) != string(o1Content) {
		t.Errorf("Copy failed. Expected\n%s\nand got\n%s", string(o1Content), string(o2Content))
		return
	}

	err = objects.Update(client, objects.UpdateOpts{
		Container: cName,
		Name:      oNames[0],
		Metadata:  metadata,
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
		err = objects.Update(client, objects.UpdateOpts{
			Container: cName,
			Name:      oNames[0],
			Metadata:  tempMap,
		})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	gr, err := objects.Get(client, objects.GetOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	om := objects.ExtractMetadata(gr)
	for k := range metadata {
		if om[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
