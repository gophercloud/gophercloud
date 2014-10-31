package lb

import (
	"errors"
	"strconv"

	"github.com/racker/perigee"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToLBListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API.
type ListOpts struct {
	ChangesSince string `q:"changes-since"`
	Status       Status `q:"status"`
	NodeAddr     string `q:"nodeaddress"`
	Marker       string `q:"marker"`
	Limit        int    `q:"limit"`
}

// ToLBListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLBListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToLBListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return LBPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type enabledState *bool

var (
	iTrue  = true
	iFalse = false

	Enabled  enabledState = &iTrue
	Disabled enabledState = &iFalse
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToLBCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Required - name of the load balancer to create. The name must be 128
	// characters or fewer in length, and all UTF-8 characters are valid.
	Name string

	// Optional - nodes to be added.
	Nodes []Node

	// Required - protocol of the service that is being load balanced.
	Protocol Protocol

	// Optional - enables or disables Half-Closed support for the load balancer.
	// Half-Closed support provides the ability for one end of the connection to
	// terminate its output, while still receiving data from the other end. Only
	// available for TCP/TCP_CLIENT_FIRST protocols.
	HalfClosed enabledState

	// Optional - the type of virtual IPs you want associated with the load
	// balancer.
	VIPs []VIP

	// Optional - the access list management feature allows fine-grained network
	// access controls to be applied to the load balancer virtual IP address.
	AccessList string

	// Optional - algorithm that defines how traffic should be directed between
	// back-end nodes.
	Algorithm Algorithm

	// Optional - current connection logging configuration.
	ConnectionLogging *ConnectionLogging

	// Optional - specifies a limit on the number of connections per IP address
	// to help mitigate malicious or abusive traffic to your applications.
	//??? ConnThrottle string

	//??? HealthMonitor string

	// Optional - arbitrary information that can be associated with each LB.
	Metadata map[string]interface{}

	// Optional - port number for the service you are load balancing.
	Port int

	// Optional - the timeout value for the load balancer and communications with
	// its nodes. Defaults to 30 seconds with a maximum of 120 seconds.
	Timeout int

	// Optional - specifies whether multiple requests from clients are directed
	// to the same node.
	//??? SessionPersistence

	// Optional - enables or disables HTTP to HTTPS redirection for the load
	// balancer. When enabled, any HTTP request returns status code 301 (Moved
	// Permanently), and the requester is redirected to the requested URL via the
	// HTTPS protocol on port 443. For example, http://example.com/page.html
	// would be redirected to https://example.com/page.html. Only available for
	// HTTPS protocol (port=443), or HTTP protocol with a properly configured SSL
	// termination (secureTrafficOnly=true, securePort=443).
	HTTPSRedirect enabledState
}

var (
	errNameRequired    = errors.New("Name is a required attribute")
	errTimeoutExceeded = errors.New("Timeout must be less than 120")
)

// ToLBCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToLBCreateMap() (map[string]interface{}, error) {
	lb := make(map[string]interface{})

	if opts.Name == "" {
		return lb, errNameRequired
	}
	if opts.Timeout > 120 {
		return lb, errTimeoutExceeded
	}

	lb["name"] = opts.Name

	if len(opts.Nodes) > 0 {
		nodes := []map[string]interface{}{}
		for _, n := range opts.Nodes {
			nodes = append(nodes, map[string]interface{}{
				"address":   n.Address,
				"port":      n.Port,
				"condition": n.Condition,
			})
		}
		lb["nodes"] = nodes
	}

	if opts.Protocol != "" {
		lb["protocol"] = opts.Protocol
	}
	if opts.HalfClosed != nil {
		lb["halfClosed"] = opts.HalfClosed
	}

	if len(opts.VIPs) > 0 {

		lb["virtualIps"] = opts.VIPs
	}

	// if opts.AccessList != "" {
	// 	lb["accessList"] = opts.AccessList
	// }
	if opts.Algorithm != "" {
		lb["algorithm"] = opts.Algorithm
	}
	if opts.ConnectionLogging != nil {
		lb["connectionLogging"] = &opts.ConnectionLogging
	}
	// if opts.ConnThrottle != "" {
	// 	lb["connectionThrottle"] = opts.ConnThrottle
	// }
	// if opts.HealthMonitor != "" {
	// 	lb["healthMonitor"] = opts.HealthMonitor
	// }
	if len(opts.Metadata) != 0 {
		lb["metadata"] = opts.Metadata
	}
	if opts.Port > 0 {
		lb["port"] = opts.Port
	}
	if opts.Timeout > 0 {
		lb["timeout"] = opts.Timeout
	}
	// if opts.SessionPersistence != "" {
	// 	lb["sessionPersistence"] = opts.SessionPersistence
	// }
	if opts.HTTPSRedirect != nil {
		lb["httpsRedirect"] = &opts.HTTPSRedirect
	}

	return map[string]interface{}{"loadBalancer": lb}, nil
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToLBCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", rootURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

func Get(c *gophercloud.ServiceClient, id int) GetResult {
	var res GetResult

	_, res.Err = perigee.Request("GET", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

func BulkDelete(c *gophercloud.ServiceClient, ids []int) DeleteResult {
	var res DeleteResult

	url := rootURL(c)
	for k, v := range ids {
		if k == 0 {
			url += "?"
		} else {
			url += "&"
		}
		url += "id=" + strconv.Itoa(v)
	}

	_, res.Err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return res
}

func Delete(c *gophercloud.ServiceClient, id int) DeleteResult {
	var res DeleteResult

	_, res.Err = perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return res
}
