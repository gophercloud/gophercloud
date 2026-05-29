package testing

import "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"

const ListBGPVPNsResult = `
{
  "bgpvpns": [
    {
      "export_targets": [
        "64512:1666"
      ],
      "name": "",
      "routers": [],
      "route_distinguishers": [
        "64512:1777",
        "64512:1888",
        "64512:1999"
      ],
      "tenant_id": "b7549121395844bea941bb92feb3fad9",
      "project_id": "b7549121395844bea941bb92feb3fad9",
      "import_targets": [
        "64512:1555"
      ],
      "route_targets": [
        "64512:1444"
      ],
      "type": "l3",
      "id": "0f9d472a-908f-40f5-8574-b4e8a63ccbf0",
      "networks": [],
      "local_pref": null,
      "vni": 1000
    }
  ]
}
`

var BGPVPN = bgpvpns.BGPVPN{
	ID:                  "0f9d472a-908f-40f5-8574-b4e8a63ccbf0",
	Name:                "",
	RouteDistinguishers: []string{"64512:1777", "64512:1888", "64512:1999"},
	RouteTargets:        []string{"64512:1444"},
	ImportTargets:       []string{"64512:1555"},
	ExportTargets:       []string{"64512:1666"},
	LocalPref:           nil,
	VNI:                 1000,
	TenantID:            "b7549121395844bea941bb92feb3fad9",
	ProjectID:           "b7549121395844bea941bb92feb3fad9",
	Type:                "l3",
	Routers:             []string{},
	Networks:            []string{},
}

const GetBGPVPNResult = `
{
    "bgpvpn": {
       "id": "460ac411-3dfb-45bb-8116-ed1a7233d143",
       "name": "foo",
       "route_targets": ["64512:1444"],
       "export_targets": [],
       "import_targets": [],
       "type": "l3",
       "tenant_id": "f94ea398564d49dfb0d542f086c68ce7",
       "project_id": "f94ea398564d49dfb0d542f086c68ce7",
       "routers": [],
       "route_distinguishers": [],
       "networks": [
         "a4f2b8df-cb42-4893-a333-d0b5c36ade17"
       ],
       "local_pref": null,
       "vni": 1000
    }
}
`

var GetBGPVPN = bgpvpns.BGPVPN{
	ID:                  "460ac411-3dfb-45bb-8116-ed1a7233d143",
	Name:                "foo",
	RouteDistinguishers: []string{},
	RouteTargets:        []string{"64512:1444"},
	ImportTargets:       []string{},
	ExportTargets:       []string{},
	LocalPref:           nil,
	VNI:                 1000,
	TenantID:            "f94ea398564d49dfb0d542f086c68ce7",
	ProjectID:           "f94ea398564d49dfb0d542f086c68ce7",
	Type:                "l3",
	Routers:             []string{},
	Networks:            []string{"a4f2b8df-cb42-4893-a333-d0b5c36ade17"},
}

const CreateRequest = `
{
  "bgpvpn": {
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "route_targets": ["64512:1444"],
    "import_targets": ["64512:1555"],
    "export_targets": ["64512:1666"],
    "route_distinguishers": ["64512:1777", "64512:1888", "64512:1999"],
    "type": "l3",
    "vni": 1000
  }
}
`

const CreateResponse = `
{
  "bgpvpn": {
    "export_targets": [
      "64512:1666"
    ],
    "name": "",
    "routers": [],
    "route_distinguishers": [
      "64512:1777",
      "64512:1888",
      "64512:1999"
    ],
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "import_targets": [
      "64512:1555"
    ],
    "route_targets": [
      "64512:1444"
    ],
    "type": "l3",
    "id": "0f9d472a-908f-40f5-8574-b4e8a63ccbf0",
    "networks": [],
    "local_pref": null,
    "vni": 1000
  }
}
`

var CreateBGPVPN = bgpvpns.BGPVPN{
	ID: "0f9d472a-908f-40f5-8574-b4e8a63ccbf0",
	RouteDistinguishers: []string{
		"64512:1777",
		"64512:1888",
		"64512:1999",
	},
	RouteTargets:  []string{"64512:1444"},
	ImportTargets: []string{"64512:1555"},
	ExportTargets: []string{"64512:1666"},
	LocalPref:     nil,
	VNI:           1000,
	TenantID:      "b7549121395844bea941bb92feb3fad9",
	ProjectID:     "b7549121395844bea941bb92feb3fad9",
	Type:          "l3",
	Routers:       []string{},
	Networks:      []string{},
}

const UpdateBGPVPNRequest = `
{
    "bgpvpn": {
       "name": "foo",
       "route_targets": ["64512:1444"],
       "export_targets": [],
       "import_targets": []
    }
}
`

const UpdateBGPVPNResponse = `
{
  "bgpvpn": {
    "export_targets": [],
    "name": "foo",
    "routers": [],
    "route_distinguishers": [
      "12345:1234"
    ],
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "import_targets": [],
    "route_targets": ["64512:1444"],
    "type": "l3",
    "id": "4d627abf-06dd-45ab-920b-8e61422bb984",
    "networks": [],
    "local_pref": null,
    "vni": 1000
  }
}
`

const ListNetworkAssociationsResult = `
{
  "network_associations": [
    {
      "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
      "network_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
      "tenant_id": "b7549121395844bea941bb92feb3fad9",
      "project_id": "b7549121395844bea941bb92feb3fad9"
    }
  ]
}
`

var NetworkAssociation = bgpvpns.NetworkAssociation{
	ID:        "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	NetworkID: "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:  "b7549121395844bea941bb92feb3fad9",
	ProjectID: "b7549121395844bea941bb92feb3fad9",
}

const GetNetworkAssociationResult = `
{
  "network_association": {
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
    "network_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9"
  }
}
`

var GetNetworkAssociation = bgpvpns.NetworkAssociation{
	ID:        "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	NetworkID: "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:  "b7549121395844bea941bb92feb3fad9",
	ProjectID: "b7549121395844bea941bb92feb3fad9",
}

const CreateNetworkAssociationRequest = `
{
  "network_association": {
    "network_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd"
  }
}
`
const CreateNetworkAssociationResponse = `
{
  "network_association": {
    "network_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
  }
}
`

var CreateNetworkAssociation = bgpvpns.NetworkAssociation{
	ID:        "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	NetworkID: "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:  "b7549121395844bea941bb92feb3fad9",
	ProjectID: "b7549121395844bea941bb92feb3fad9",
}

const ListRouterAssociationsResult = `
{
  "router_associations": [
    {
      "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
      "router_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
      "tenant_id": "b7549121395844bea941bb92feb3fad9",
      "project_id": "b7549121395844bea941bb92feb3fad9"
    }
  ]
}
`

var RouterAssociation = bgpvpns.RouterAssociation{
	ID:        "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	RouterID:  "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:  "b7549121395844bea941bb92feb3fad9",
	ProjectID: "b7549121395844bea941bb92feb3fad9",
}

const GetRouterAssociationResult = `
{
  "router_association": {
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
    "router_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9"
  }
}
`

var GetRouterAssociation = bgpvpns.RouterAssociation{
	ID:        "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	RouterID:  "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:  "b7549121395844bea941bb92feb3fad9",
	ProjectID: "b7549121395844bea941bb92feb3fad9",
}

const CreateRouterAssociationRequest = `
{
  "router_association": {
    "router_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd"
  }
}
`
const CreateRouterAssociationResponse = `
{
  "router_association": {
    "router_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
    "advertise_extra_routes": true
  }
}
`

var CreateRouterAssociation = bgpvpns.RouterAssociation{
	ID:                   "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	RouterID:             "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:             "b7549121395844bea941bb92feb3fad9",
	ProjectID:            "b7549121395844bea941bb92feb3fad9",
	AdvertiseExtraRoutes: true,
}

const UpdateRouterAssociationRequest = `
{
  "router_association": {
    "advertise_extra_routes": false
  }
}
`
const UpdateRouterAssociationResponse = `
{
  "router_association": {
    "router_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
  }
}
`

var UpdateRouterAssociation = bgpvpns.RouterAssociation{
	ID:                   "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	RouterID:             "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:             "b7549121395844bea941bb92feb3fad9",
	ProjectID:            "b7549121395844bea941bb92feb3fad9",
	AdvertiseExtraRoutes: false,
}

const ListPortAssociationsResult = `
{
  "port_associations": [
    {
      "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
      "port_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
      "tenant_id": "b7549121395844bea941bb92feb3fad9",
      "project_id": "b7549121395844bea941bb92feb3fad9",
      "advertise_fixed_ips": true
    }
  ]
}
`

var PortAssociation = bgpvpns.PortAssociation{
	ID:                "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	PortID:            "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:          "b7549121395844bea941bb92feb3fad9",
	ProjectID:         "b7549121395844bea941bb92feb3fad9",
	AdvertiseFixedIPs: true,
}

const GetPortAssociationResult = `
{
  "port_association": {
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
    "port_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "advertise_fixed_ips": true
  }
}
`

var GetPortAssociation = bgpvpns.PortAssociation{
	ID:                "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	PortID:            "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:          "b7549121395844bea941bb92feb3fad9",
	ProjectID:         "b7549121395844bea941bb92feb3fad9",
	AdvertiseFixedIPs: true,
}

const CreatePortAssociationRequest = `
{
  "port_association": {
    "port_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd"
  }
}
`
const CreatePortAssociationResponse = `
{
  "port_association": {
    "port_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
    "advertise_fixed_ips": true
  }
}
`

var CreatePortAssociation = bgpvpns.PortAssociation{
	ID:                "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	PortID:            "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:          "b7549121395844bea941bb92feb3fad9",
	ProjectID:         "b7549121395844bea941bb92feb3fad9",
	AdvertiseFixedIPs: true,
}

const UpdatePortAssociationRequest = `
{
  "port_association": {
    "advertise_fixed_ips": false
  }
}
`
const UpdatePortAssociationResponse = `
{
  "port_association": {
    "port_id": "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
    "tenant_id": "b7549121395844bea941bb92feb3fad9",
    "project_id": "b7549121395844bea941bb92feb3fad9",
    "id": "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
    "advertise_fixed_ips": false
  }
}
`

var UpdatePortAssociation = bgpvpns.PortAssociation{
	ID:                "73238ca1-e05d-4c7a-b4d4-70407b4b8730",
	PortID:            "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	TenantID:          "b7549121395844bea941bb92feb3fad9",
	ProjectID:         "b7549121395844bea941bb92feb3fad9",
	AdvertiseFixedIPs: false,
}
