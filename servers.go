// TODO(sfalvo): Remove Rackspace-specific Server structure fields and refactor them into a provider-specific access method.
// Be sure to update godocs accordingly.

package gophercloud

import (
	"github.com/racker/perigee"
)

// genericServersProvider structures provide the implementation for generic OpenStack-compatible
// CloudServersProvider interfaces.
type genericServersProvider struct {
	// endpoint refers to the provider's API endpoint base URL.  This will be used to construct
	// and issue queries.
	endpoint string

	// Test context (if any) in which to issue requests.
	context *Context

	// access associates this API provider with a set of credentials,
	// which may be automatically renewed if they near expiration.
	access AccessProvider
}

// See the CloudServersProvider interface for details.
func (gcp *genericServersProvider) ListServers() ([]Server, error) {
	var ss []Server

	url := gcp.endpoint + "/servers"
	err := perigee.Get(url, perigee.Options{
		CustomClient: gcp.context.httpClient,
		Results:      &struct{ Servers *[]Server }{&ss},
		MoreHeaders: map[string]string{
			"X-Auth-Token": gcp.access.AuthToken(),
		},
	})
	return ss, err
}

// RaxBandwidth provides measurement of server bandwidth consumed over a given audit interval.
type RaxBandwidth struct {
	AuditPeriodEnd    string `json:"audit_period_end"`
	AuditPeriodStart  string `json:"audit_period_start"`
	BandwidthInbound  int64  `json:"bandwidth_inbound"`
	BandwidthOutbound int64  `json:"bandwidth_outbound"`
	Interface         string `json:"interface"`
}

// A VersionedAddress denotes either an IPv4 or IPv6 (depending on version indicated)
// address.
type VersionedAddress struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
}

// An AddressSet provides a set of public and private IP addresses for a resource.
// Each address has a version to identify if IPv4 or IPv6.
type AddressSet struct {
	Public  []VersionedAddress `json:"public"`
	Private []VersionedAddress `json:"private"`
}

// Server records represent (virtual) hardware instances (not configurations) accessible by the user.
//
// The AccessIPv4 / AccessIPv6 fields provides IP addresses for the server in the IPv4 or IPv6 format, respectively.
//
// Addresses provides addresses for any attached isolated networks.
// The version field indicates whether the IP address is version 4 or 6.
//
// Created tells when the server entity was created.
//
// The Flavor field includes the flavor ID and flavor links.
//
// The compute provisioning algorithm has an anti-affinity property that
// attempts to spread customer VMs across hosts.
// Under certain situations,
// VMs from the same customer might be placed on the same host.
// The HostId field represents the host your server runs on and
// can be used to determine this scenario if it is relevant to your application.
// Note that HostId is unique only per account; it is not globally unique.
// 
// Id provides the server's unique identifier.
// This field must be treated opaquely.
//
// Image indicates which image is installed on the server.
//
// Links provides one or more means of accessing the server.
//
// Metadata provides a small key-value store for application-specific information.
//
// Name provides a human-readable name for the server.
//
// Progress indicates how far along it is towards being provisioned.
// 100 represents complete, while 0 represents just beginning.
//
// Status provides an indication of what the server's doing at the moment.
// A server will be in ACTIVE state if it's ready for use.
//
// OsDcfDiskConfig indicates the server's boot volume configuration.
// Valid values are:
//     AUTO
//     ----
//     The server is built with a single partition the size of the target flavor disk.
//     The file system is automatically adjusted to fit the entire partition.
//     This keeps things simple and automated.
//     AUTO is valid only for images and servers with a single partition that use the EXT3 file system.
//     This is the default setting for applicable Rackspace base images.
//
//     MANUAL
//     ------
//     The server is built using whatever partition scheme and file system is in the source image.
//     If the target flavor disk is larger,
//     the remaining disk space is left unpartitioned.
//     This enables images to have non-EXT3 file systems, multiple partitions, and so on,
//     and enables you to manage the disk configuration.
//
// RaxBandwidth provides measures of the server's inbound and outbound bandwidth per interface.
//
// OsExtStsPowerState provides an indication of the server's power.
// This field appears to be a set of flag bits:
//
//           ... 4  3   2   1   0
//         +--//--+---+---+---+---+
//         | .... | 0 | S | 0 | I |
//         +--//--+---+---+---+---+
//                      |       |
//                      |       +---  0=Instance is down.
//                      |             1=Instance is up.
//                      |
//                      +-----------  0=Server is switched ON.
//                                    1=Server is switched OFF.
//                                    (note reverse logic.)
//
// Unused bits should be ignored when read, and written as 0 for future compatibility.
//
// OsExtStsTaskState and OsExtStsVmState work together
// to provide visibility in the provisioning process for the instance.
// Consult Rackspace documentation at
// http://docs.rackspace.com/servers/api/v2/cs-devguide/content/ch_extensions.html#ext_status
// for more details.  It's too lengthy to include here.
type Server struct {
	AccessIPv4         string         `json:"accessIPv4"`
	AccessIPv6         string         `json:"accessIPv6"`
	Addresses          AddressSet     `json:"addresses"`
	Created            string         `json:"created"`
	Flavor             FlavorLink     `json:"flavor"`
	HostId             string         `json:"hostId"`
	Id                 string         `json:"id"`
	Image              ImageLink      `json:"image"`
	Links              []Link         `json:"links"`
	Metadata           interface{}    `json:"metadata"`
	Name               string         `json:"name"`
	Progress           int            `json:"progress"`
	Status             string         `json:"status"`
	TenantId           string         `json:"tenant_id"`
	Updated            string         `json:"updated"`
	UserId             string         `json:"user_id"`
	OsDcfDiskConfig    string         `json:"OS-DCF:diskConfig"`
	RaxBandwidth       []RaxBandwidth `json:"rax-bandwidth:bandwidth"`
	OsExtStsPowerState int            `json:"OS-EXT-STS:power_state"`
	OsExtStsTaskState  string         `json:"OS-EXT-STS:task_state"`
	OsExtStsVmState    string         `json:"OS-EXT-STS:vm_state"`
}
