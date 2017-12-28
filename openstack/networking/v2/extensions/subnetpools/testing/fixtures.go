package testing

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
)

const SubnetPoolsListResult = `
{
	  "subnetpools": [
	  		{
	  				"address_scope_id": null,
	  				"created_at": "2017-12-28T07:21:41Z",
	  				"default_prefixlen": 8,
	  				"default_quota": null,
	  				"description": "IPv4",
	  				"id": "d43a57fe-3390-4608-b437-b1307b0adb40",
	  				"ip_version": 4,
	  				"is_default": false,
	  				"max_prefixlen": 32,
	  				"min_prefixlen": 8,
	  				"name": "MyPoolIpv4",
	  				"prefixes": [
	  						"10.10.10.0/24",
	  						"10.11.11.0/24"
	  				],
	  				"project_id": "1e2b9857295a4a3e841809ef492812c5",
	  				"revision_number": 1,
	  				"shared": false,
	  				"tenant_id": "1e2b9857295a4a3e841809ef492812c5",
	  				"updated_at": "2017-12-28T07:21:41Z"
	  		},
	  		{
	  				"address_scope_id": "0bc38e22-be49-4e67-969e-fec3f36508bd",
	  				"created_at": "2017-12-28T07:21:34Z",
	  				"default_prefixlen": 64,
	  				"default_quota": null,
	  				"description": "IPv6",
	  				"id": "832cb7f3-59fe-40cf-8f64-8350ffc03272",
	  				"ip_version": 6,
	  				"is_default": true,
	  				"max_prefixlen": 128,
	  				"min_prefixlen": 64,
	  				"name": "MyPoolIpv6",
	  				"prefixes": [
	  						"fdf7:b13d:dead:beef::/64",
	  						"fd65:86cc:a334:39b7::/64"
	  				],
	  				"project_id": "1e2b9857295a4a3e841809ef492812c5",
	  				"revision_number": 1,
	  				"shared": false,
	  				"tenant_id": "1e2b9857295a4a3e841809ef492812c5",
	  				"updated_at": "2017-12-28T07:21:34Z"
	  		},
	  		{
	  				"address_scope_id": null,
	  				"created_at": "2017-12-28T07:21:27Z",
	  				"default_prefixlen": 64,
	  				"default_quota": 4,
	  				"description": "PublicPool",
	  				"id": "2fe18ae6-58c2-4a85-8bfb-566d6426749b",
	  				"ip_version": 6,
	  				"is_default": false,
	  				"max_prefixlen": 128,
	  				"min_prefixlen": 64,
	  				"name": "PublicIPv6",
	  				"prefixes": [
	  						"2001:db8::a3/64"
	  				],
	  				"project_id": "ceb366d50ad54fe39717df3af60f9945",
	  				"revision_number": 1,
	  				"shared": true,
	  				"tenant_id": "ceb366d50ad54fe39717df3af60f9945",
	  				"updated_at": "2017-12-28T07:21:27Z"
	  		}
	  ]
}
`

var SubnetPool1 = subnetpools.SubnetPool{
	AddressScopeID:   "",
	CreatedAt:        "2017-12-28T07:21:41Z",
	DefaultPrefixLen: 8,
	DefaultQuota:     0,
	Description:      "IPv4",
	ID:               "d43a57fe-3390-4608-b437-b1307b0adb40",
	IPversion:        4,
	IsDefault:        false,
	MaxPrefixLen:     32,
	MinPrefixLen:     8,
	Name:             "MyPoolIpv4",
	Prefixes: []string{
		"10.10.10.0/24",
		"10.11.11.0/24",
	},
	ProjectID:      "1e2b9857295a4a3e841809ef492812c5",
	TenantID:       "1e2b9857295a4a3e841809ef492812c5",
	RevisionNumber: 1,
	Shared:         false,
	UpdatedAt:      "2017-12-28T07:21:41Z",
}

var SubnetPool2 = subnetpools.SubnetPool{
	AddressScopeID:   "0bc38e22-be49-4e67-969e-fec3f36508bd",
	CreatedAt:        "2017-12-28T07:21:34Z",
	DefaultPrefixLen: 64,
	DefaultQuota:     0,
	Description:      "IPv6",
	ID:               "832cb7f3-59fe-40cf-8f64-8350ffc03272",
	IPversion:        6,
	IsDefault:        true,
	MaxPrefixLen:     128,
	MinPrefixLen:     64,
	Name:             "MyPoolIpv6",
	Prefixes: []string{
		"fdf7:b13d:dead:beef::/64",
		"fd65:86cc:a334:39b7::/64",
	},
	ProjectID:      "1e2b9857295a4a3e841809ef492812c5",
	TenantID:       "1e2b9857295a4a3e841809ef492812c5",
	RevisionNumber: 1,
	Shared:         false,
	UpdatedAt:      "2017-12-28T07:21:34Z",
}

var SubnetPool3 = subnetpools.SubnetPool{
	AddressScopeID:   "",
	CreatedAt:        "2017-12-28T07:21:27Z",
	DefaultPrefixLen: 64,
	DefaultQuota:     4,
	Description:      "PublicPool",
	ID:               "2fe18ae6-58c2-4a85-8bfb-566d6426749b",
	IPversion:        6,
	IsDefault:        false,
	MaxPrefixLen:     128,
	MinPrefixLen:     64,
	Name:             "PublicIPv6",
	Prefixes: []string{
		"2001:db8::a3/64",
	},
	ProjectID:      "ceb366d50ad54fe39717df3af60f9945",
	TenantID:       "ceb366d50ad54fe39717df3af60f9945",
	RevisionNumber: 1,
	Shared:         true,
	UpdatedAt:      "2017-12-28T07:21:27Z",
}
