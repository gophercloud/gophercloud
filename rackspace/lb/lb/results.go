package lb

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type Protocol string

// The constants below represent all the compatible load balancer protocols.
const (
	// DNSTCP is a protocol that works with IPv6 and allows your DNS server to
	// receive traffic using TCP port 53.
	DNSTCP = "DNS_TCP"

	// DNSUDP is a protocol that works with IPv6 and allows your DNS server to
	// receive traffic using UDP port 53.
	DNSUDP = "DNS_UDP"

	// TCP is one of the core protocols of the Internet Protocol Suite. It
	// provides a reliable, ordered delivery of a stream of bytes from one
	// program on a computer to another program on another computer. Applications
	// that require an ordered and reliable delivery of packets use this protocol.
	TCP = "TCP"

	// TCPCLIENTFIRST is a protocol similar to TCP, but is more efficient when a
	// client is expected to write the data first.
	TCPCLIENTFIRST = "TCP_CLIENT_FIRST"

	// UDP provides a datagram service that emphasizes speed over reliability. It
	// works well with applications that provide security through other measures.
	UDP = "UDP"

	// UDPSTREAM is a protocol designed to stream media over networks and is
	// built on top of UDP.
	UDPSTREAM = "UDP_STREAM"
)

// Algorithm defines how traffic should be directed between back-end nodes.
type Algorithm string

const (
	// LC directs traffic to the node with the lowest number of connections.
	LC = "LEAST_CONNECTIONS"

	// RAND directs traffic to nodes at random.
	RAND = "RANDOM"

	// RR directs traffic to each of the nodes in turn.
	RR = "ROUND_ROBIN"

	// WLC directs traffic to a node based on the number of concurrent
	// connections and its weight.
	WLC = "WEIGHTED_LEAST_CONNECTIONS"

	// WRR directs traffic to a node according to the RR algorithm, but with
	// different proportions of traffic being directed to the back-end nodes.
	// Weights must be defined as part of the node configuration.
	WRR = "WEIGHTED_ROUND_ROBIN"
)

type Status string

const (
	// ACTIVE indicates that the LB is configured properly and ready to serve
	// traffic to incoming requests via the configured virtual IPs.
	ACTIVE = "ACTIVE"

	// BUILD indicates that the LB is being provisioned for the first time and
	// configuration is being applied to bring the service online. The service
	// cannot yet serve incoming requests.
	BUILD = "BUILD"

	// PENDINGUPDATE indicates that the LB is online but configuration changes
	// are being applied to update the service based on a previous request.
	PENDINGUPDATE = "PENDING_UPDATE"

	// PENDINGDELETE indicates that the LB is online but configuration changes
	// are being applied to begin deletion of the service based on a previous
	// request.
	PENDINGDELETE = "PENDING_DELETE"

	// SUSPENDED indicates that the LB has been taken offline and disabled.
	SUSPENDED = "SUSPENDED"

	// ERROR indicates that the system encountered an error when attempting to
	// configure the load balancer.
	ERROR = "ERROR"

	// DELETED indicates that the LB has been deleted.
	DELETED = "DELETED"
)

type Datetime struct {
	Time string
}

type VIP struct {
	Address string
	ID      int
	Type    string
	Version string `mapstructure:"ipVersion"`
}

type LoadBalancer struct {
	// Human-readable name for the load balancer.
	Name string

	// The unique ID for the load balancer.
	ID int

	// Represents the service protocol being load balanced.
	Protocol Protocol

	// Defines how traffic should be directed between back-end nodes. The default
	// algorithm is RANDOM.
	Algorithm Algorithm

	// The current status of the load balancer.
	Status Status

	// The number of load balancer nodes.
	NodeCount int `mapstructure:"nodeCount"`

	// Slice of virtual IPs associated with this load balancer.
	VIPs []VIP `mapstructure:"virtualIps"`

	// Datetime when the LB was created.
	Created Datetime

	// Datetime when the LB was created.
	Updated Datetime

	Port int
}

// LBPage is the page returned by a pager when traversing over a collection of
// LBs.
type LBPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (p LBPage) IsEmpty() (bool, error) {
	is, err := ExtractLBs(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractLBs accepts a Page struct, specifically a LBPage struct, and extracts
// the elements into a slice of LoadBalancer structs. In other words, a generic
// collection is mapped into a relevant slice.
func ExtractLBs(page pagination.Page) ([]LoadBalancer, error) {
	var resp struct {
		LBs []LoadBalancer `mapstructure:"loadBalancers" json:"loadBalancers"`
	}

	err := mapstructure.Decode(page.(LBPage).Body, &resp)

	return resp.LBs, err
}
