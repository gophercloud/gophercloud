package bgpvpns

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToBGPVPNListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Fields    []string `q:"fields"`
	ProjectID string   `q:"project_id"`
	Networks  []string `q:"networks"`
	Routers   []string `q:"routers"`
	Ports     []string `q:"ports"`
	Limit     int      `q:"limit"`
	Marker    string   `q:"marker"`
}

// ToBGPVPNListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToBGPVPNListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List the BGP VPNs
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	query, err := opts.ToBGPVPNListQuery()
	if err != nil {
		return pagination.Pager{Err: err}
	}
	url += query
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := BGPVPNPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// Get retrieve the specific BGP VPN by its uuid
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToBGPVPNCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a BGP VPN.
type CreateOpts struct {
	Name                string   `json:"name,omitempty"`
	RouteDistinguishers []string `json:"route_distinguishers,omitempty"`
	RouteTargets        []string `json:"route_targets,omitempty"`
	ImportTargets       []string `json:"import_targets,omitempty"`
	ExportTargets       []string `json:"export_targets,omitempty"`
	LocalPref           int      `json:"local_pref,omitempty"`
	VNI                 int      `json:"vni,omitempty"`
	TenantID            string   `json:"tenant_id,omitempty"`
	ProjectID           string   `json:"project_id,omitempty"`
	Type                string   `json:"type,omitempty"`
}

// ToBGPVPNCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToBGPVPNCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bgpvpn")
}

// Create a BGP VPN
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToBGPVPNCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the BGP VPN associated with it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToBGPVPNUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a BGP VPN.
type UpdateOpts struct {
	Name                *string   `json:"name,omitempty"`
	RouteDistinguishers *[]string `json:"route_distinguishers,omitempty"`
	RouteTargets        *[]string `json:"route_targets,omitempty"`
	ImportTargets       *[]string `json:"import_targets,omitempty"`
	ExportTargets       *[]string `json:"export_targets,omitempty"`
	LocalPref           *int      `json:"local_pref,omitempty"`
}

// ToBGPVPNUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToBGPVPNUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bgpvpn")
}

// Update accept a BGP VPN ID and an UpdateOpts and update the BGP VPN
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToBGPVPNUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListNetworkAssociationsOptsBuilder allows extensions to add additional
// parameters to the ListNetworkAssociations request.
type ListNetworkAssociationsOptsBuilder interface {
	ToNetworkAssociationsListQuery() (string, error)
}

// ListNetworkAssociationsOpts allows the filtering and sorting of paginated
// collections through the API.
type ListNetworkAssociationsOpts struct {
	Fields []string `q:"fields"`
	Limit  int      `q:"limit"`
	Marker string   `q:"marker"`
}

// ToNetworkAssociationsListQuery formats a ListNetworkAssociationsOpts into a
// query string.
func (opts ListNetworkAssociationsOpts) ToNetworkAssociationsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// ListNetworkAssociations pages over the network associations of a specified
// BGP VPN.
func ListNetworkAssociations(c *gophercloud.ServiceClient, id string, opts ListNetworkAssociationsOptsBuilder) pagination.Pager {
	url := listNetworkAssociationsURL(c, id)
	query, err := opts.ToNetworkAssociationsListQuery()
	if err != nil {
		return pagination.Pager{Err: err}
	}
	url += query
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := NetworkAssociationPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// CreateNetworkAssociationOptsBuilder allows extensions to add additional
// parameters to the CreateNetworkAssociation request.
type CreateNetworkAssociationOptsBuilder interface {
	ToNetworkAssociationCreateMap() (map[string]interface{}, error)
}

// CreateNetworkAssociationOpts represents options used to create a BGP VPN
// network association.
type CreateNetworkAssociationOpts struct {
	NetworkID string `json:"network_id" required:"true"`
	TenantID  string `json:"tenant_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// ToNetworkAssociationCreateMap builds a request body from
// CreateNetworkAssociationOpts.
func (opts CreateNetworkAssociationOpts) ToNetworkAssociationCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "network_association")
}

// CreateNetworkAssociation creates a new network association for a specified
// BGP VPN.
func CreateNetworkAssociation(ctx context.Context, client *gophercloud.ServiceClient, id string, opts CreateNetworkAssociationOptsBuilder) (r CreateNetworkAssociationResult) {
	b, err := opts.ToNetworkAssociationCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createNetworkAssociationURL(client, id), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetNetworkAssociation retrieves a specific network association by BGP VPN id
// and network association id.
func GetNetworkAssociation(ctx context.Context, c *gophercloud.ServiceClient, bgpVpnID string, id string) (r GetNetworkAssociationResult) {
	resp, err := c.Get(ctx, getNetworkAssociationURL(c, bgpVpnID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteNetworkAssociation deletes a specific network association by BGP VPN id
// and network association id.
func DeleteNetworkAssociation(ctx context.Context, c *gophercloud.ServiceClient, bgpVpnID string, id string) (r DeleteNetworkAssociationResult) {
	resp, err := c.Delete(ctx, deleteNetworkAssociationURL(c, bgpVpnID, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListRouterAssociationsOptsBuilder allows extensions to add additional
// parameters to the ListRouterAssociations request.
type ListRouterAssociationsOptsBuilder interface {
	ToRouterAssociationsListQuery() (string, error)
}

// ListRouterAssociationsOpts allows the filtering and sorting of paginated
// collections through the API.
type ListRouterAssociationsOpts struct {
	Fields []string `q:"fields"`
	Limit  int      `q:"limit"`
	Marker string   `q:"marker"`
}

// ToRouterAssociationsListQuery formats a ListRouterAssociationsOpts into a
// query string.
func (opts ListRouterAssociationsOpts) ToRouterAssociationsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// ListRouterAssociations pages over the router associations of a specified
// BGP VPN.
func ListRouterAssociations(c *gophercloud.ServiceClient, id string, opts ListRouterAssociationsOptsBuilder) pagination.Pager {
	url := listRouterAssociationsURL(c, id)
	query, err := opts.ToRouterAssociationsListQuery()
	if err != nil {
		return pagination.Pager{Err: err}
	}
	url += query
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := RouterAssociationPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// CreateRouterAssociationOptsBuilder allows extensions to add additional
// parameters to the CreateRouterAssociation request.
type CreateRouterAssociationOptsBuilder interface {
	ToRouterAssociationCreateMap() (map[string]interface{}, error)
}

// CreateRouterAssociationOpts represents options used to create a BGP VPN
// router association.
type CreateRouterAssociationOpts struct {
	RouterID             string `json:"router_id" required:"true"`
	AdvertiseExtraRoutes *bool  `json:"advertise_extra_routes,omitempty"`
	TenantID             string `json:"tenant_id,omitempty"`
	ProjectID            string `json:"project_id,omitempty"`
}

// ToRouterAssociationCreateMap builds a request body from
// CreateRouterAssociationOpts.
func (opts CreateRouterAssociationOpts) ToRouterAssociationCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "router_association")
}

// CreateRouterAssociation creates a new router association for a specified
// BGP VPN.
func CreateRouterAssociation(ctx context.Context, client *gophercloud.ServiceClient, id string, opts CreateRouterAssociationOptsBuilder) (r CreateRouterAssociationResult) {
	b, err := opts.ToRouterAssociationCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createRouterAssociationURL(client, id), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetRouterAssociation retrieves a specific router association by BGP VPN id
// and router association id.
func GetRouterAssociation(ctx context.Context, c *gophercloud.ServiceClient, bgpVpnID string, id string) (r GetRouterAssociationResult) {
	resp, err := c.Get(ctx, getRouterAssociationURL(c, bgpVpnID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteRouterAssociation deletes a specific router association by BGP VPN id
// and router association id.
func DeleteRouterAssociation(ctx context.Context, c *gophercloud.ServiceClient, bgpVpnID string, id string) (r DeleteRouterAssociationResult) {
	resp, err := c.Delete(ctx, deleteRouterAssociationURL(c, bgpVpnID, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateRouterAssociationOptsBuilder allows extensions to add additional
// parameters to the UpdateRouterAssociation request.
type UpdateRouterAssociationOptsBuilder interface {
	ToRouterAssociationUpdateMap() (map[string]interface{}, error)
}

// UpdateRouterAssociationOpts represents options used to update a BGP VPN
// router association.
type UpdateRouterAssociationOpts struct {
	AdvertiseExtraRoutes *bool `json:"advertise_extra_routes,omitempty"`
}

// ToRouterAssociationUpdateMap builds a request body from
// UpdateRouterAssociationOpts.
func (opts UpdateRouterAssociationOpts) ToRouterAssociationUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "router_association")
}

// UpdateRouterAssociation updates a router association for a specified BGP VPN.
func UpdateRouterAssociation(ctx context.Context, client *gophercloud.ServiceClient, bgpVpnID string, id string, opts UpdateRouterAssociationOptsBuilder) (r UpdateRouterAssociationResult) {
	b, err := opts.ToRouterAssociationUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateRouterAssociationURL(client, bgpVpnID, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListPortAssociationsOptsBuilder allows extensions to add additional
// parameters to the ListPortAssociations request.
type ListPortAssociationsOptsBuilder interface {
	ToPortAssociationsListQuery() (string, error)
}

// ListPortAssociationsOpts allows the filtering and sorting of paginated
// collections through the API.
type ListPortAssociationsOpts struct {
	Fields []string `q:"fields"`
	Limit  int      `q:"limit"`
	Marker string   `q:"marker"`
}

// ToPortAssociationsListQuery formats a ListPortAssociationsOpts into a
// query string.
func (opts ListPortAssociationsOpts) ToPortAssociationsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// ListPortAssociations pages over the port associations of a specified
// BGP VPN.
func ListPortAssociations(c *gophercloud.ServiceClient, id string, opts ListPortAssociationsOptsBuilder) pagination.Pager {
	url := listPortAssociationsURL(c, id)
	query, err := opts.ToPortAssociationsListQuery()
	if err != nil {
		return pagination.Pager{Err: err}
	}
	url += query
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := PortAssociationPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// PortRoutes represents the routes to be advertised by a BGP VPN port
type PortRoutes struct {
	Type      string `json:"type" required:"true"`
	Prefix    string `json:"prefix,omitempty"`
	BGPVPNID  string `json:"bgpvpn_id,omitempty"`
	LocalPref *int   `json:"local_pref,omitempty"`
}

// CreatePortAssociationOptsBuilder allows extensions to add additional
// parameters to the CreatePortAssociation request.
type CreatePortAssociationOptsBuilder interface {
	ToPortAssociationCreateMap() (map[string]interface{}, error)
}

// CreatePortAssociationOpts represents options used to create a BGP VPN
// port association.
type CreatePortAssociationOpts struct {
	PortID            string       `json:"port_id" required:"true"`
	Routes            []PortRoutes `json:"routes,omitempty"`
	AdvertiseFixedIPs *bool        `json:"advertise_fixed_ips,omitempty"`
	TenantID          string       `json:"tenant_id,omitempty"`
	ProjectID         string       `json:"project_id,omitempty"`
}

// ToPortAssociationCreateMap builds a request body from
// CreatePortAssociationOpts.
func (opts CreatePortAssociationOpts) ToPortAssociationCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "port_association")
}

// CreatePortAssociation creates a new port association for a specified
// BGP VPN.
func CreatePortAssociation(ctx context.Context, client *gophercloud.ServiceClient, id string, opts CreatePortAssociationOptsBuilder) (r CreatePortAssociationResult) {
	b, err := opts.ToPortAssociationCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createPortAssociationURL(client, id), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetPortAssociation retrieves a specific port association by BGP VPN id
// and port association id.
func GetPortAssociation(ctx context.Context, c *gophercloud.ServiceClient, bgpVpnID string, id string) (r GetPortAssociationResult) {
	resp, err := c.Get(ctx, getPortAssociationURL(c, bgpVpnID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeletePortAssociation deletes a specific port association by BGP VPN id
// and port association id.
func DeletePortAssociation(ctx context.Context, c *gophercloud.ServiceClient, bgpVpnID string, id string) (r DeletePortAssociationResult) {
	resp, err := c.Delete(ctx, deletePortAssociationURL(c, bgpVpnID, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdatePortAssociationOptsBuilder allows extensions to add additional
// parameters to the UpdatePortAssociation request.
type UpdatePortAssociationOptsBuilder interface {
	ToPortAssociationUpdateMap() (map[string]interface{}, error)
}

// UpdatePortAssociationOpts represents options used to update a BGP VPN
// port association.
type UpdatePortAssociationOpts struct {
	Routes            *[]PortRoutes `json:"routes,omitempty"`
	AdvertiseFixedIPs *bool         `json:"advertise_fixed_ips,omitempty"`
}

// ToPortAssociationUpdateMap builds a request body from
// UpdatePortAssociationOpts.
func (opts UpdatePortAssociationOpts) ToPortAssociationUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "port_association")
}

// UpdatePortAssociation updates a port association for a specified BGP VPN.
func UpdatePortAssociation(ctx context.Context, client *gophercloud.ServiceClient, bgpVpnID string, id string, opts UpdatePortAssociationOptsBuilder) (r UpdatePortAssociationResult) {
	b, err := opts.ToPortAssociationUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updatePortAssociationURL(client, bgpVpnID, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
