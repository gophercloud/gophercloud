package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const ClusterTemplateResponse = `
{
	"insecure_registry": null,
	"links": [
	{
	  "href": "http://65.61.151.130:9511/v1/baymodels/472807c2-f175-4946-9765-149701a5aba7",
	  "rel": "self"
	},
	{
	  "href": "http://65.61.151.130:9511/baymodels/472807c2-f175-4946-9765-149701a5aba7",
	  "rel": "bookmark"
	}
	],
	"http_proxy": null,
	"updated_at": null,
	"fixed_subnet": null,
	"master_flavor_id": null,
	"uuid": "472807c2-f175-4946-9765-149701a5aba7",
	"no_proxy": null,
	"https_proxy": null,
	"tls_disabled": false,
	"keypair_id": "testkey",
	"public": false,
	"labels": {},
	"docker_volume_size": 5,
	"server_type": "vm",
	"external_network_id": "public",
	"cluster_distro": "fedora-atomic",
	"image_id": "fedora-atomic-latest",
	"volume_driver": null,
	"registry_enabled": false,
	"docker_storage_driver": null,
	"apiserver_port": null,
	"name": "kubernetes-dev",
	"created_at": "2016-08-10T13:47:01+00:00",
	"created_at2": "2015-02-10T14:26:14Z",
	"network_driver": "flannel",
	"fixed_network": null,
	"coe": "kubernetes",
	"flavor_id": "m1.small",
	"master_lb_enabled": false,
	"dns_nameserver": "8.8.8.8"
}`

var ExpectedClusterTemplate = clustertemplates.ClusterTemplate{
	InsecureRegistry: "",
	Links: []gophercloud.Link{
		{Href: "http://65.61.151.130:9511/v1/baymodels/472807c2-f175-4946-9765-149701a5aba7", Rel: "self"},
		{Href: "http://65.61.151.130:9511/baymodels/472807c2-f175-4946-9765-149701a5aba7", Rel: "bookmark"},
	},
	HTTPProxy:           "",
	UpdatedAt:           time.Time{},
	FixedSubnet:         "",
	MasterFlavorID:      "",
	UUID:                "472807c2-f175-4946-9765-149701a5aba7",
	NoProxy:             "",
	HTTPSProxy:          "",
	TLSDisabled:         false,
	KeyPairID:           "testkey",
	Public:              false,
	Labels:              map[string]string{},
	DockerVolumeSize:    5,
	ServerType:          "vm",
	ExternalNetworkID:   "public",
	ClusterDistro:       "fedora-atomic",
	ImageID:             "fedora-atomic-latest",
	VolumeDriver:        "",
	RegistryEnabled:     false,
	DockerStorageDriver: "",
	APIServerPort:       "",
	Name:                "kubernetes-dev",
	CreatedAt:           time.Date(2016, 8, 10, 13, 47, 01, 0, time.UTC),
	NetworkDriver:       "flannel",
	FixedNetwork:        "",
	COE:                 "kubernetes",
	FlavorID:            "m1.small",
	MasterLBEnabled:     false,
	DNSNameServer:       "8.8.8.8",
}

const ClusterTemplateResponse_EmptyTime = `
{
	"insecure_registry": null,
	"links": [
	{
	  "href": "http://65.61.151.130:9511/v1/baymodels/472807c2-f175-4946-9765-149701a5aba7",
	  "rel": "self"
	},
	{
	  "href": "http://65.61.151.130:9511/baymodels/472807c2-f175-4946-9765-149701a5aba7",
	  "rel": "bookmark"
	}
	],
	"http_proxy": null,
	"updated_at": null,
	"fixed_subnet": null,
	"master_flavor_id": null,
	"uuid": "472807c2-f175-4946-9765-149701a5aba7",
	"no_proxy": null,
	"https_proxy": null,
	"tls_disabled": false,
	"keypair_id": "testkey",
	"public": false,
	"labels": {},
	"docker_volume_size": 5,
	"server_type": "vm",
	"external_network_id": "public",
	"cluster_distro": "fedora-atomic",
	"image_id": "fedora-atomic-latest",
	"volume_driver": null,
	"registry_enabled": false,
	"docker_storage_driver": null,
	"apiserver_port": null,
	"name": "kubernetes-dev",
	"created_at": null,
	"network_driver": "flannel",
	"fixed_network": null,
	"coe": "kubernetes",
	"flavor_id": "m1.small",
	"master_lb_enabled": false,
	"dns_nameserver": "8.8.8.8"
}`

var ExpectedClusterTemplate_EmptyTime = clustertemplates.ClusterTemplate{
	InsecureRegistry: "",
	Links: []gophercloud.Link{
		{Href: "http://65.61.151.130:9511/v1/baymodels/472807c2-f175-4946-9765-149701a5aba7", Rel: "self"},
		{Href: "http://65.61.151.130:9511/baymodels/472807c2-f175-4946-9765-149701a5aba7", Rel: "bookmark"},
	},
	HTTPProxy:           "",
	UpdatedAt:           time.Time{},
	FixedSubnet:         "",
	MasterFlavorID:      "",
	UUID:                "472807c2-f175-4946-9765-149701a5aba7",
	NoProxy:             "",
	HTTPSProxy:          "",
	TLSDisabled:         false,
	KeyPairID:           "testkey",
	Public:              false,
	Labels:              map[string]string{},
	DockerVolumeSize:    5,
	ServerType:          "vm",
	ExternalNetworkID:   "public",
	ClusterDistro:       "fedora-atomic",
	ImageID:             "fedora-atomic-latest",
	VolumeDriver:        "",
	RegistryEnabled:     false,
	DockerStorageDriver: "",
	APIServerPort:       "",
	Name:                "kubernetes-dev",
	CreatedAt:           time.Time{},
	NetworkDriver:       "flannel",
	FixedNetwork:        "",
	COE:                 "kubernetes",
	FlavorID:            "m1.small",
	MasterLBEnabled:     false,
	DNSNameServer:       "8.8.8.8",
}

func HandleCreateClusterTemplateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clustertemplates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", "req-781e9bdc-4163-46eb-91c9-786c53188bbb")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ClusterTemplateResponse)
	})
}

func HandleCreateClusterTemplateEmptyTimeSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clustertemplates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ClusterTemplateResponse_EmptyTime)
	})
}
