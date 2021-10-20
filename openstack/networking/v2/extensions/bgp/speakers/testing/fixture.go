package testing

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speakers"

const ListBGPSpeakerResult = `
{
  "bgp_speakers": [
    {
      "peers": [
        "afacc0e8-6b66-44e4-be53-a1ef16033ceb",
        "acd7c4a1-e243-4fe5-80f9-eba8f143ac1d"
      ],
      "advertise_floating_ip_host_routes": true,
      "name": "gophercloud-testing-speaker",
      "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "local_as": 56789,
      "id": "ab01ade1-ae62-43c9-8a1f-3c24225b96d8",
      "ip_version": 4,
      "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "networks": [
        "acdc6339-7d2d-411f-82bb-e6cc3ad9eb9f"
      ],
      "advertise_tenant_networks": true
    }
  ]
}
`

var BGPSpeaker1 = speakers.BGPSpeaker{
	ID:                            "ab01ade1-ae62-43c9-8a1f-3c24225b96d8",
	Name:                          "gophercloud-testing-speaker",
	TenantID:                      "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	ProjectID:                     "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	AdvertiseFloatingIPHostRoutes: true,
	AdvertiseTenantNetworks:       true,
	IPVersion:                     4,
	LocalAS:                       56789,
	Networks:                      []string{"acdc6339-7d2d-411f-82bb-e6cc3ad9eb9f"},
	Peers: []string{"afacc0e8-6b66-44e4-be53-a1ef16033ceb",
		"acd7c4a1-e243-4fe5-80f9-eba8f143ac1d"},
}

const GetBGPSpeakerResult = `
{
  "bgp_speaker": {
    "peers": [
      "afacc0e8-6b66-44e4-be53-a1ef16033ceb",
      "acd7c4a1-e243-4fe5-80f9-eba8f143ac1d"
    ],
    "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "name": "gophercloud-testing-speaker",
    "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "local_as": 56789,
    "advertise_tenant_networks": true,
    "networks": [
      "acdc6339-7d2d-411f-82bb-e6cc3ad9eb9f"
    ],
    "ip_version": 4,
    "advertise_floating_ip_host_routes": true,
    "id": "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
  }
}
`
