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
		"href": "http://10.63.176.154:9511/v1/clustertemplates/79c0f9e5-93b8-4719-8fab-063afc67bffe",
		"rel": "self"
	},
	{
		"href": "http://10.63.176.154:9511/clustertemplates/79c0f9e5-93b8-4719-8fab-063afc67bffe",
		"rel": "bookmark"
	}
	],
	"http_proxy": "http://10.164.177.169:8080",
	"updated_at": null,
	"floating_ip_enabled": true,
	"fixed_subnet": null,
	"master_flavor_id": null,
	"user_id": "c48d66144e9c4a54ae2b164b85cfefe3",
	"uuid": "79c0f9e5-93b8-4719-8fab-063afc67bffe",
	"no_proxy": "10.0.0.0/8,172.0.0.0/8,192.0.0.0/8,localhost",
	"https_proxy": "http://10.164.177.169:8080",
	"tls_disabled": false,
	"keypair_id": "kp",
	"project_id": "76bd201dbc1641729904ab190d3390c6",
	"public": false,
	"labels": null,
	"docker_volume_size": 3,
	"server_type": "vm",
	"external_network_id": "public",
	"cluster_distro": "fedora-atomic",
	"image_id": "Fedora-Atomic-27-20180212.2.x86_64",
	"volume_driver": "cinder",
	"registry_enabled": false,
	"docker_storage_driver": "devicemapper",
	"apiserver_port": null,
	"name": "kubernetes-dev",
	"created_at": "2018-06-27T16:52:21+00:00",
	"network_driver": "flannel",
	"fixed_network": null,
	"coe": "kubernetes",
	"flavor_id": "m1.small",
	"master_lb_enabled": true,
	"dns_nameserver": "8.8.8.8"
}`

var ExpectedClusterTemplate = clustertemplates.ClusterTemplate{
	InsecureRegistry: "",
	Links: []gophercloud.Link{
		{Href: "http://10.63.176.154:9511/v1/clustertemplates/79c0f9e5-93b8-4719-8fab-063afc67bffe", Rel: "self"},
		{Href: "http://10.63.176.154:9511/clustertemplates/79c0f9e5-93b8-4719-8fab-063afc67bffe", Rel: "bookmark"},
	},
	HTTPProxy:           "http://10.164.177.169:8080",
	UpdatedAt:           time.Time{},
	FloatingIPEnabled:   true,
	FixedSubnet:         "",
	MasterFlavorID:      "",
	UserID:              "c48d66144e9c4a54ae2b164b85cfefe3",
	UUID:                "79c0f9e5-93b8-4719-8fab-063afc67bffe",
	NoProxy:             "10.0.0.0/8,172.0.0.0/8,192.0.0.0/8,localhost",
	HTTPSProxy:          "http://10.164.177.169:8080",
	TLSDisabled:         false,
	KeyPairID:           "kp",
	ProjectID:           "76bd201dbc1641729904ab190d3390c6",
	Public:              false,
	Labels:              map[string]string(nil),
	DockerVolumeSize:    3,
	ServerType:          "vm",
	ExternalNetworkID:   "public",
	ClusterDistro:       "fedora-atomic",
	ImageID:             "Fedora-Atomic-27-20180212.2.x86_64",
	VolumeDriver:        "cinder",
	RegistryEnabled:     false,
	DockerStorageDriver: "devicemapper",
	APIServerPort:       "",
	Name:                "kubernetes-dev",
	CreatedAt:           time.Date(2018, 6, 27, 16, 52, 21, 0, time.UTC),
	NetworkDriver:       "flannel",
	FixedNetwork:        "",
	COE:                 "kubernetes",
	FlavorID:            "m1.small",
	MasterLBEnabled:     true,
	DNSNameServer:       "8.8.8.8",
}

func HandleCreateClusterTemplateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clustertemplates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("OpenStack-API-Minimum-Version", "container-infra 1.1")
		w.Header().Add("OpenStack-API-Maximum-Version", "container-infra 1.6")
		w.Header().Add("OpenStack-API-Version", "container-infra 1.1")
		w.Header().Add("X-OpenStack-Request-Id", "req-781e9bdc-4163-46eb-91c9-786c53188bbb")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, ClusterTemplateResponse)
	})
}

func HandleDeleteClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clustertemplates/6dc6d336e3fc4c0a951b5698cd1236ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("OpenStack-API-Minimum-Version", "container-infra 1.1")
		w.Header().Add("OpenStack-API-Maximum-Version", "container-infra 1.6")
		w.Header().Add("OpenStack-API-Version", "container-infra 1.1")
		w.Header().Add("X-OpenStack-Request-Id", "req-781e9bdc-4163-46eb-91c9-786c53188bbb")
		w.WriteHeader(http.StatusNoContent)
	})
}
