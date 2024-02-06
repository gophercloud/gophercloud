package testing

import "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/quotas"

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

const UpdateRequestResponseRaw_1 = `
{
    "quota": {
        "loadbalancer": 20,
        "listener": 40,
        "member": 200,
        "pool": 20,
        "healthmonitor": -1,
        "l7policy": 50,
        "l7rule": 100
    }
}
`

const UpdateRequestResponseRaw_2 = `
{
    "quota": {
        "load_balancer": 20,
        "listener": 40,
        "member": 200,
        "pool": 20,
        "health_monitor": -1,
        "l7policy": 50,
        "l7rule": 100
    }
}
`

var UpdateResponse = quotas.Quota{
	Loadbalancer:  20,
	Listener:      40,
	Member:        200,
	Pool:          20,
	Healthmonitor: -1,
	L7Policy:      50,
	L7Rule:        100,
}
