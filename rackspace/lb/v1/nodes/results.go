package nodes

// Node represents a back-end device, usually a virtual machine, that can
// handle traffic. It is assigned traffic based on its parent load balancer.
type Node struct {
	// The IP address or CIDR for this back-end node.
	Address string

	// The unique ID for this node.
	ID int

	// The port on which traffic is sent and received.
	Port int

	// The node's status.
	Status Status

	// The node's condition.
	Condition string

	// The priority at which this node will receive traffic if a weighted
	// algorithm is used by its parent load balancer. Ranges from 1 to 100.
	Weight int

	// Type of node.
	Type Type
}

// Type indicates whether the node is of a PRIMARY or SECONDARY nature.
type Type string

const (
	// PRIMARY nodes are in the normal rotation to receive traffic from the load
	// balancer.
	PRIMARY Type = "PRIMARY"

	// SECONDARY nodes are only in the rotation to receive traffic from the load
	// balancer when all the primary nodes fail. This provides a failover feature
	// that automatically routes traffic to the secondary node in the event that
	// the primary node is disabled or in a failing state. Note that active
	// health monitoring must be enabled on the load balancer to enable the
	// failover feature to the secondary node.
	SECONDARY Type = "SECONDARY"
)

type Condition string

const (
	// ENABLED indicates that the node is permitted to accept new connections.
	ENABLED Condition = "ENABLED"

	// DISABLED indicates that the node is not permitted to accept any new
	// connections regardless of session persistence configuration. Existing
	// connections are forcibly terminated.
	DISABLED Condition = "DISABLED"

	// DRAINING indicates that the node is allowed to service existing
	// established connections and connections that are being directed to it as a
	// result of the session persistence configuration.
	DRAINING
)

// Status indicates whether the node can accept service traffic. If a node is
// not listening on its port or does not meet the conditions of the defined
// active health check for the load balancer, then the load balancer does not
// forward connections and its status is listed as OFFLINE
type Status string

const (
	// ONLINE indicates that the node is healthy and capable of receiving traffic
	// from the load balancer.
	ONLINE Status = "ONLINE"

	// OFFLINE indicates that the node is not in a position to receive service
	// traffic. It is usually switched into this state when a health check is not
	// satisfied with the node's response time.
	OFFLINE Status = "OFFLINE"
)
