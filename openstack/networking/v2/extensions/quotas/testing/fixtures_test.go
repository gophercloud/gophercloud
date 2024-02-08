package testing

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
)

const GetResponseRaw = `
{
    "quota": {
        "floatingip": 15,
        "network": 20,
        "port": 25,
        "rbac_policy": -1,
        "router": 30,
        "security_group": 35,
        "security_group_rule": 40,
        "subnet": 45,
        "subnetpool": -1,
        "trunk": 50
    }
}
`

// GetDetailedResponseRaw is a sample response to a Get call with the detailed option.
//
// One "reserved" property is returned as a string to reflect a buggy behaviour
// of Neutron.
//
// cf. https://bugs.launchpad.net/neutron/+bug/1918565
const GetDetailedResponseRaw = `
{
   "quota" : {
      "floatingip" : {
          "used": 0,
          "limit": 15,
          "reserved": 0
      },
      "network" : {
          "used": 0,
          "limit": 20,
          "reserved": 0
      },
      "port" : {
          "used": 0,
          "limit": 25,
          "reserved": "0"
      },
      "rbac_policy" : {
          "used": 0,
          "limit": -1,
          "reserved": 0
      },
      "router" : {
          "used": 0,
          "limit": 30,
          "reserved": 0
      },
      "security_group" : {
          "used": 0,
          "limit": 35,
          "reserved": 0
      },
      "security_group_rule" : {
          "used": 0,
          "limit": 40,
          "reserved": 0
      },
      "subnet" : {
          "used": 0,
          "limit": 45,
          "reserved": 0
      },
      "subnetpool" : {
          "used": 0,
          "limit": -1,
          "reserved": 0
      },
      "trunk" : {
          "used": 0,
          "limit": 50,
          "reserved": 0
      }
   }
}
`

var GetResponse = quotas.Quota{
	FloatingIP:        15,
	Network:           20,
	Port:              25,
	RBACPolicy:        -1,
	Router:            30,
	SecurityGroup:     35,
	SecurityGroupRule: 40,
	Subnet:            45,
	SubnetPool:        -1,
	Trunk:             50,
}

// GetDetailResponse is the first result in ListOutput.
var GetDetailResponse = quotas.QuotaDetailSet{
	FloatingIP:        quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 15},
	Network:           quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 20},
	Port:              quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 25},
	RBACPolicy:        quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: -1},
	Router:            quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 30},
	SecurityGroup:     quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 35},
	SecurityGroupRule: quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 40},
	Subnet:            quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 45},
	SubnetPool:        quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: -1},
	Trunk:             quotas.QuotaDetail{Used: 0, Reserved: 0, Limit: 50},
}

const UpdateRequestResponseRaw = `
{
    "quota": {
        "floatingip": 0,
        "network": -1,
        "port": 5,
        "rbac_policy": 10,
        "router": 15,
        "security_group": 20,
        "security_group_rule": -1,
        "subnet": 25,
        "subnetpool": 0,
        "trunk": 5
    }
}
`

var UpdateResponse = quotas.Quota{
	FloatingIP:        0,
	Network:           -1,
	Port:              5,
	RBACPolicy:        10,
	Router:            15,
	SecurityGroup:     20,
	SecurityGroupRule: -1,
	Subnet:            25,
	SubnetPool:        0,
	Trunk:             5,
}
