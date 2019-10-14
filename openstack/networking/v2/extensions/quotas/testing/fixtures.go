package testing

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"

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
        "subnetpool": -1
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
}
