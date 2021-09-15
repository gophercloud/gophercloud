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
    "advertise_floating_ip_host_routes": false,
    "advertise_tenant_networks": true,
    "ip_version": 6,
    "local_as": "2000",
    "name": "gophercloud-testing-bgp-speaker"
  }
}
`

const CreateResponse = `
{
  "bgp_speaker": {
    "peers": [],
    "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "name": "gophercloud-testing-bgp-speaker",
    "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "local_as": 2000,
    "advertise_tenant_networks": true,
    "networks": [],
    "ip_version": 6,
    "advertise_floating_ip_host_routes": false,
    "id": "26e98af2-4dc7-452a-91b0-65ee45f3e7c1"
  }
}
`

const UpdateBGPSpeakerRequest = `
{
  "bgp_speaker": {
    "advertise_floating_ip_host_routes": true,
    "advertise_tenant_networks": false,
    "name": "testing-bgp-speaker"
  }
}
`

const UpdateBGPSpeakerResponse = `
{
  "bgp_speaker": {
    "peers": [],
    "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "name": "testing-bgp-speaker",
    "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "local_as": 2000,
    "advertise_tenant_networks": false,
    "networks": [],
    "ip_version": 4,
    "advertise_floating_ip_host_routes": true,
    "id": "d25d0036-7f17-49d7-8d02-4bf9dd49d5a9"
  }
}
`

const GetAdvertisedRoutesResult = `
{
  "advertised_routes": [
    {
      "next_hop": "172.17.128.212",
      "destination": "172.17.129.192/27"
    },
    {
      "next_hop": "172.17.128.218",
      "destination": "172.17.129.0/27"
    },
    {
      "next_hop": "172.17.128.231",
      "destination": "172.17.129.160/27"
    }
  ]
}
`

const AddBGPPeerRequest = `
{
  "bgp_peer_id": "f5884c7c-71d5-43a3-88b4-1742e97674aa"
}
`

const AddRemoveBGPPeerJSON = `
{
  "bgp_peer_id": "f5884c7c-71d5-43a3-88b4-1742e97674aa"
}
`

const AddRemoveGatewayNetworkJSON = `
{
  "network_id": "ac13bb26-6219-49c3-a880-08847f6830b7"
}
`
