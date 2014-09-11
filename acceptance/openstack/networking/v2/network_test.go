// +build acceptance networking

package v2

import (
	"fmt"
	"testing"

	networking "github.com/rackspace/gophercloud/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/utils"
)

var Client *networking.Client

func NewClient() (*networking.Client, error) {
	provider := gophercloud.AuthenticatedClient(utils.AuthOptions())

	url, err := provider.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}

	return &gophercloud.ServiceClient{Provider: provider, Endpoint: url}, nil
}

func Setup() {
	client, err = NewClient()
	if err != nil {
		fmt.Println("Client failed to load")
		return
	}

	Client = client
}

func Teardown() {
	Client = nil
}

func TestListAPIVersions(t *testing.T) {
	Setup()
	defer Teardown()

	res, err := networks.APIVersions(Client)
	if err != nil {
		t.Fatalf("Failed to list API versions")
	}

	err = gophercloud.EachPage(res, func(page gophercloud.Collection) bool {
		t.Logf("--- Page ---")
		for _, v := range networks.ToAPIVersions(page) {
			t.Logf("API version: ID [%s] Status [%s]", v.ID, v.Status)
		}
		return true
	})
	if err != nil {
		t.Fatalf("Unexpected error while iterating API versions")
	}
}

func TestGetApiInfo(t *testing.T) {
	networks.ApiInfo()
}

func TestListExts(t *testing.T) {
	networks.Extensions()
}

func TestGetExt(t *testing.T) {
	networks.Extension()
}

func TestListNetworks(t *testing.T) {
	networks.List()
}

func TestGetNetwork(t *testing.T) {
	networks.Get()
}

func TestCreateNetwork(t *testing.T) {
	networks.Create()
}

func TestCreateMultipleNetworks(t *testing.T) {
	networks.CreateMany()
}

func TestUpdateNetwork(t *testing.T) {
	networks.Update()
}

func TestDeleteNetwork(t *testing.T) {
	networks.Delete()
}
