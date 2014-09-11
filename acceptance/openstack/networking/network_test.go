// +build acceptance networking

package networking

import (
	"fmt"
	"os"
	"testing"

	networking "github.com/rackspace/gophercloud/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/utils"
)

var Client *networking.Client

func NewClient() (*networking.Client, error) {
	opts, err := utils.AuthOptions()
	if err != nil {
		return nil, err
	}

	res, err := identity.Authenticate(opts)
	if err != nil {
		return nil, err
	}

	catalog, err := identity.GetServiceCatalog(res)
	if err != nil {
		return nil, err
	}

	entries, err := catalog.CatalogEntries()
	if err != nil {
		return nil, err
	}

	var endpoints []identity.Endpoint
	for _, entry := range entries {
		if entry.Type == "network" {
			endpoints = entry.Endpoints
		}
	}

	region := os.Getenv("OS_REGION_NAME")

	url := ""
	for _, endpoint := range endpoints {
		if endpoint.Region == "region" {
			url = endpoint.PublicURL
		}
	}

	client := networking.NewClient(url, res, opts)
	return client, nil
}

type SuiteTester struct {
	suite.Suite
	Client *networking.Client
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

func TestListApiVersions(t *testing.T) {
	networks.ApiVersions()
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
