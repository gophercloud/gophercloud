package testing

import "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/quotas"

const GetResponseRaw_1 = `
{
    "quota": {
        "loadbalancer": 15,
        "listener": 30,
        "member": -1,
        "pool": 15,
        "healthmonitor": 30,
        "l7policy": 100,
        "l7rule": -1
    }
}
`

const GetResponseRaw_2 = `
{
    "quota": {
        "load_balancer": 15,
        "listener": 30,
        "member": -1,
        "pool": 15,
        "health_monitor": 30,
        "l7policy": 100,
        "l7rule": -1
    }
}
`

var GetResponse = quotas.Quota{
	Loadbalancer:  15,
	Listener:      30,
	Member:        -1,
	Pool:          15,
	Healthmonitor: 30,
	L7Policy:      100,
	L7Rule:        -1,
}
