package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateClusterTemplate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateClusterTemplateSuccessfully(t)

	opts := clustertemplates.CreateOpts{
		Name:                "kubernetes-dev",
		Labels:              map[string]string{},
		FixedSubnet:         "",
		MasterFlavorID:      "",
		NoProxy:             "10.0.0.0/8,172.0.0.0/8,192.0.0.0/8,localhost",
		HTTPSProxy:          "http://10.164.177.169:8080",
		TLSDisabled:         false,
		KeyPairID:           "kp",
		Public:              false,
		HTTPProxy:           "http://10.164.177.169:8080",
		DockerVolumeSize:    3,
		ServerType:          "vm",
		ExternalNetworkID:   "public",
		ImageID:             "Fedora-Atomic-27-20180212.2.x86_64",
		VolumeDriver:        "cinder",
		RegistryEnabled:     false,
		DockerStorageDriver: "devicemapper",
		NetworkDriver:       "flannel",
		FixedNetwork:        "",
		COE:                 "kubernetes",
		FlavorID:            "m1.small",
		MasterLBEnabled:     true,
		DNSNameServer:       "8.8.8.8",
	}

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	res := clustertemplates.Create(sc, opts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, "req-781e9bdc-4163-46eb-91c9-786c53188bbb", requestID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)

	actual.CreatedAt = actual.CreatedAt.UTC()
	th.AssertDeepEquals(t, ExpectedClusterTemplate, *actual)
}
