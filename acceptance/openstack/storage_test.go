// +build acceptance

package openstack

import (
	"bytes"
	"github.com/rackspace/gophercloud/openstack/storage"
	"github.com/rackspace/gophercloud/openstack/storage/accounts"
	"github.com/rackspace/gophercloud/openstack/storage/containers"
	"github.com/rackspace/gophercloud/openstack/storage/objects"
	"os"
	"strings"
	"testing"
)

var objectStorage = "object-store"
var metadata = map[string]string{"gopher": "cloud"}
var numContainers = 2
var numObjects = 2

func TestAccount(t *testing.T) {
	ts, err := setupForList(objectStorage)
	if err != nil {
		t.Error(err)
		return
	}

	region := os.Getenv("OS_REGION_NAME")
	for _, ep := range ts.eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := storage.NewClient(ep.PublicURL, ts.a, ts.o)

		err := accounts.Update(client, accounts.UpdateOpts{
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
		am := accounts.GetMetadata(gr)
		for k := range metadata {
			if am[k] != metadata[strings.Title(k)] {
				t.Errorf("Expected custom metadata with key: %s", k)
				return
			}
		}
	}
}

func TestContainers(t *testing.T) {
	ts, err := setupForList(objectStorage)
	if err != nil {
		t.Error(err)
		return
	}

	region := os.Getenv("OS_REGION_NAME")
	for _, ep := range ts.eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := storage.NewClient(ep.PublicURL, ts.a, ts.o)

		cNames := make([]string, numContainers)
		for i := 0; i < numContainers; i++ {
			cNames[i] = randomString("test-container-", 8)
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
		cns, err := containers.GetNames(lr)
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
		cis, err := containers.GetInfo(lr)
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
		cm := containers.GetMetadata(gr)
		for k := range metadata {
			if cm[k] != metadata[strings.Title(k)] {
				t.Errorf("Expected custom metadata with key: %s", k)
				return
			}
		}
	}
}

func TestObjects(t *testing.T) {
	ts, err := setupForList(objectStorage)
	if err != nil {
		t.Error(err)
		return
	}

	region := os.Getenv("OS_REGION_NAME")

	for _, ep := range ts.eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := storage.NewClient(ep.PublicURL, ts.a, ts.o)

		oNames := make([]string, numObjects)
		for i := 0; i < len(oNames); i++ {
			oNames[i] = randomString("test-object-", 8)
		}

		cName := randomString("test-container-", 8)
		_, err := containers.Create(client, containers.CreateOpts{
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
			oContents[i] = bytes.NewBuffer([]byte(randomString("", 10)))
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
		ons := objects.GetNames(lr)
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
		ois, err := objects.GetInfo(lr)
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
		o2Content := objects.GetContent(dr)

		dr, err = objects.Download(client, objects.DownloadOpts{
			Container: cName,
			Name:      oNames[0],
		})
		if err != nil {
			t.Error(err)
			return
		}
		o1Content := objects.GetContent(dr)

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
		om := objects.GetMetadata(gr)
		for k := range metadata {
			if om[k] != metadata[strings.Title(k)] {
				t.Errorf("Expected custom metadata with key: %s", k)
				return
			}
		}
	}
}
