package v1

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	th "github.com/gophercloud/gophercloud/testhelper"
)

var Client *gophercloud.ServiceClient

func NewClient() (*gophercloud.ServiceClient, error) {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return openstack.NewContainerOrchestrationV1(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

func Setup(t *testing.T) {
	client, err := NewClient()
	th.AssertNoErr(t, err)
	Client = client
}

func Teardown() {
	Client = nil
}
