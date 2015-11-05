package volumeactions

import (
	"github.com/rackspace/gophercloud"
)

type AttachOptsBuilder interface {
	ToVolumeAttachMap() (map[string]interface{}, error)
}

type AttachOpts struct {
	// The mountpoint of this volume
	MountPoint string
	// The nova instance ID, can't set simultaneously with HostName
	InstanceUUID string
	// The hostname of baremetal host, can't set simultaneously with InstanceUUID
	HostName string
	// Mount mode of this volume
	Mode string
}

func (opts AttachOpts) ToVolumeAttachMap() (map[string]interface{}, error) {
	v := make(map[string]interface{})

	if opts.MountPoint != "" {
		v["mountpoint"] = opts.MountPoint
	}
	if opts.Mode != "" {
		v["mode"] = opts.Mode
	}
	if opts.InstanceUUID != "" {
		v["instance_uuid"] = opts.InstanceUUID
	}
	if opts.HostName != "" {
		v["host_name"] = opts.HostName
	}

	return map[string]interface{}{"os-attach": v}, nil
}

func Attach(client *gophercloud.ServiceClient, id string, opts AttachOptsBuilder) AttachResult {
	var res AttachResult

	reqBody, err := opts.ToVolumeAttachMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(attachURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	return res
}

func Detach(client *gophercloud.ServiceClient, id string) DetachResult {
	var res DetachResult

	v := make(map[string]interface{})
	reqBody := map[string]interface{}{"os-detach": v}

	_, res.Err = client.Post(detachURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	return res
}

func Reserve(client *gophercloud.ServiceClient, id string) ReserveResult {
	var res ReserveResult

	v := make(map[string]interface{})
	reqBody := map[string]interface{}{"os-reserve": v}

	_, res.Err = client.Post(reserveURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	return res
}

func Unreserve(client *gophercloud.ServiceClient, id string) UnreserveResult {
	var res UnreserveResult

	v := make(map[string]interface{})
	reqBody := map[string]interface{}{"os-unreserve": v}

	_, res.Err = client.Post(unreserveURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	return res
}

type ConnectorOptsBuilder interface {
	ToConnectorMap() (map[string]interface{}, error)
}

type ConnectorOpts struct {
	IP        string
	Host      string
	Initiator string
	Wwpns     string
	Wwnns     string
	Multipath bool
	Platform  string
	OSType    string
}

func (opts ConnectorOpts) ToConnectorMap() (map[string]interface{}, error) {
	v := make(map[string]interface{})

	if opts.IP != "" {
		v["ip"] = opts.IP
	}
	if opts.Host != "" {
		v["host"] = opts.Host
	}
	if opts.Initiator != "" {
		v["initiator"] = opts.Initiator
	}
	if opts.Wwpns != "" {
		v["wwpns"] = opts.Wwpns
	}
	if opts.Wwnns != "" {
		v["wwnns"] = opts.Wwnns
	}

	v["multipath"] = opts.Multipath

	if opts.Platform != "" {
		v["platform"] = opts.Platform
	}
	if opts.OSType != "" {
		v["os_type"] = opts.OSType
	}

	return map[string]interface{}{"connector": v}, nil
}

func InitializeConnection(client *gophercloud.ServiceClient, id string, opts *ConnectorOpts) InitializeConnectionResult {
	var res InitializeConnectionResult

	connctorMap, err := opts.ToConnectorMap()
	if err != nil {
		res.Err = err
		return res
	}

	reqBody := map[string]interface{}{"os-initialize_connection": connctorMap}

	_, res.Err = client.Post(initializeConnectionURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	return res
}

func TerminateConnection(client *gophercloud.ServiceClient, id string, opts *ConnectorOpts) TerminateConnectionResult {
	var res TerminateConnectionResult

	connctorMap, err := opts.ToConnectorMap()
	if err != nil {
		res.Err = err
		return res
	}

	reqBody := map[string]interface{}{"os-terminate_connection": connctorMap}

	_, res.Err = client.Post(teminateConnectionURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	return res
}
