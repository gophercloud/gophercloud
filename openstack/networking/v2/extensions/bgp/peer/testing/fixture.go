package testing

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/peer"

const ListBGPPeersResult = `
{
  "bgp_peers": [
    {
      "auth_type": "none",
      "remote_as": 4321,
      "name": "testing-peer-1",
      "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "peer_ip": "1.2.3.4",
      "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "id": "afacc0e8-6b66-44e4-be53-a1ef16033ceb"
    },
    {
      "auth_type": "none",
      "remote_as": 4321,
      "name": "testing-peer-2",
      "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "peer_ip": "5.6.7.8",
      "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "id": "acd7c4a1-e243-4fe5-80f9-eba8f143ac1d"
    }
  ]
}
`

var BGPPeer1 = peer.BGPPeer{
	AuthType:  "none",
	ID:        "afacc0e8-6b66-44e4-be53-a1ef16033ceb",
	Name:      "testing-peer-1",
	TenantID:  "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	PeerIP:    "1.2.3.4",
	ProjectID: "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	RemoteAS:  4321,
}

var BGPPeer2 = peer.BGPPeer{
	AuthType:  "none",
	ID:        "acd7c4a1-e243-4fe5-80f9-eba8f143ac1d",
	Name:      "testing-peer-2",
	TenantID:  "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	PeerIP:    "5.6.7.8",
	ProjectID: "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	RemoteAS:  4321,
}
