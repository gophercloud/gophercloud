package testing

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speaker"

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

var BGPSpeaker1 = speaker.BGPSpeaker{
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

const CreateRequest = `
{
  "bgp_speaker": {
    "advertise_floating_ip_host_routes": true,
    "advertise_tenant_networks": true,
    "ip_version": 4,
    "local_as": "2000",
    "name": "gophercloud-testing-bgp-speaker"
  }
}
`

const CreateResponse = `
{
  "bgp_speaker": {
    "peers": [],
    "project_id": "bb18eab692114b45aed901f880508a5a",
    "name": "gophercloud-testing-bgp-speaker",
    "tenant_id": "bb18eab692114b45aed901f880508a5a",
    "local_as": 2000,
    "advertise_tenant_networks": false,
    "networks": [],
    "ip_version": 6,
    "advertise_floating_ip_host_routes": false,
    "id": "26e98af2-4dc7-452a-91b0-65ee45f3e7c1"
  }
}
`
