package bgpvpns

import (
	"net/url"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

const (
	invalidMarker = "-1"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a BGP VPN resource.
func (r commonResult) Extract() (*BGPVPN, error) {
	var s BGPVPN
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "bgpvpn")
}

// BGPVPN represents an MPLS network with which Neutron routers and/or networks
// may be associated
type BGPVPN struct {
	// The ID of the BGP VPN.
	ID string `json:"id"`

	// The user meaningful name of the BGP VPN.
	Name string `json:"name"`

	// Selection of the type of VPN and the technology behind it. Allowed
	// values are l2 or l3.
	Type string `json:"type"`

	// Indicates whether this BGP VPN is shared across tenants.
	Shared bool `json:"shared"`

	// List of route distinguisher strings. If this parameter is specified,
	// one of these RDs will be used to advertise VPN routes.
	RouteDistinguishers []string `json:"route_distinguishers"`

	// Route Targets that will be both imported and used for export.
	RouteTargets []string `json:"route_targets"`

	// Additional Route Targets that will be imported.
	ImportTargets []string `json:"import_targets"`

	// Additional Route Targets that will be used for export.
	ExportTargets []string `json:"export_targets"`

	// This read-only list of network IDs reflects the associations defined
	// by Network association API resources.
	Networks []string `json:"networks"`

	// This read-only list of router IDs reflects the associations defined
	// by Router association API resources.
	Routers []string `json:"routers"`

	// This read-only list of port IDs reflects the associations defined by
	// Port association API resources (only present if the
	// bgpvpn-routes-control API extension is enabled).
	Ports []string `json:"ports"`

	// The default BGP LOCAL_PREF of routes that will be advertised to the
	// BGPVPN (unless overridden per-route).
	LocalPref *int `json:"local_pref"`

	// The globally-assigned VXLAN vni for the BGP VPN.
	VNI int `json:"vni"`

	// The ID of the project.
	TenantID string `json:"tenant_id"`

	// The ID of the project.
	ProjectID string `json:"project_id"`
}

// BGPVPNPage is the page returned by a pager when traversing over a
// collection of BGP VPNs.
type BGPVPNPage struct {
	pagination.MarkerPageBase
}

// NextPageURL generates the URL for the page of results after this one.
func (r BGPVPNPage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == invalidMarker {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// LastMarker returns the last offset in a ListResult.
func (r BGPVPNPage) LastMarker() (string, error) {
	results, err := ExtractBGPVPNs(r)
	if err != nil {
		return invalidMarker, err
	}
	if len(results) == 0 {
		return invalidMarker, nil
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return invalidMarker, err
	}
	queryParams := u.Query()
	limit := queryParams.Get("limit")

	// Limit is not present, only one page required
	if limit == "" {
		return invalidMarker, nil
	}

	return results[len(results)-1].ID, nil
}

// IsEmpty checks whether a BGPPage struct is empty.
func (r BGPVPNPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractBGPVPNs(r)
	return len(is) == 0, err
}

// ExtractBGPVPNs accepts a Page struct, specifically a BGPVPNPage struct,
// and extracts the elements into a slice of BGPVPN structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBGPVPNs(r pagination.Page) ([]BGPVPN, error) {
	var s []BGPVPN
	err := ExtractBGPVPNsInto(r, &s)
	return s, err
}

func ExtractBGPVPNsInto(r pagination.Page, v any) error {
	return r.(BGPVPNPage).Result.ExtractIntoSlicePtr(v, "bgpvpns")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a BGPVPN.
type GetResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to intepret it as a BGPVPN.
type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a BGPVPN.
type UpdateResult struct {
	commonResult
}

type commonNetworkAssociationResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a BGP VPN resource.
func (r commonNetworkAssociationResult) Extract() (*NetworkAssociation, error) {
	var s NetworkAssociation
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonNetworkAssociationResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "network_association")
}

// NetworkAssociation represents a BGP VPN network association object.
type NetworkAssociation struct {
	ID        string `json:"id"`
	NetworkID string `json:"network_id"`
	TenantID  string `json:"tenant_id"`
	ProjectID string `json:"project_id"`
}

// NetworkAssociationPage is the page returned by a pager when traversing over a
// collection of network associations.
type NetworkAssociationPage struct {
	pagination.MarkerPageBase
}

// NextPageURL generates the URL for the page of results after this one.
func (r NetworkAssociationPage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == invalidMarker {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// LastMarker returns the last offset in a ListResult.
func (r NetworkAssociationPage) LastMarker() (string, error) {
	results, err := ExtractNetworkAssociations(r)
	if err != nil {
		return invalidMarker, err
	}
	if len(results) == 0 {
		return invalidMarker, nil
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return invalidMarker, err
	}
	queryParams := u.Query()
	limit := queryParams.Get("limit")

	// Limit is not present, only one page required
	if limit == "" {
		return invalidMarker, nil
	}

	return results[len(results)-1].ID, nil
}

// IsEmpty checks whether a NetworkAssociationPage struct is empty.
func (r NetworkAssociationPage) IsEmpty() (bool, error) {
	is, err := ExtractNetworkAssociations(r)
	return len(is) == 0, err
}

// ExtractNetworkAssociations accepts a Page struct, specifically a NetworkAssociationPage struct,
// and extracts the elements into a slice of NetworkAssociation structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractNetworkAssociations(r pagination.Page) ([]NetworkAssociation, error) {
	var s []NetworkAssociation
	err := ExtractNetworkAssociationsInto(r, &s)
	return s, err
}

func ExtractNetworkAssociationsInto(r pagination.Page, v interface{}) error {
	return r.(NetworkAssociationPage).Result.ExtractIntoSlicePtr(v, "network_associations")
}

// CreateNetworkAssociationResult represents the result of a create operation. Call its Extract
// method to interpret it as a NetworkAssociation.
type CreateNetworkAssociationResult struct {
	commonNetworkAssociationResult
}

// GetNetworkAssociationResult represents the result of a get operation. Call its Extract
// method to interpret it as a NetworkAssociation.
type GetNetworkAssociationResult struct {
	commonNetworkAssociationResult
}

// DeleteNetworkAssociationResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteNetworkAssociationResult struct {
	gophercloud.ErrResult
}

type commonRouterAssociationResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a BGP VPN resource.
func (r commonRouterAssociationResult) Extract() (*RouterAssociation, error) {
	var s RouterAssociation
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonRouterAssociationResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "router_association")
}

// RouterAssociation represents a BGP VPN router association object.
type RouterAssociation struct {
	ID                   string `json:"id"`
	RouterID             string `json:"router_id"`
	TenantID             string `json:"tenant_id"`
	ProjectID            string `json:"project_id"`
	AdvertiseExtraRoutes bool   `json:"advertise_extra_routes"`
}

// RouterAssociationPage is the page returned by a pager when traversing over a
// collection of router associations.
type RouterAssociationPage struct {
	pagination.MarkerPageBase
}

// NextPageURL generates the URL for the page of results after this one.
func (r RouterAssociationPage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == invalidMarker {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// LastMarker returns the last offset in a ListResult.
func (r RouterAssociationPage) LastMarker() (string, error) {
	results, err := ExtractRouterAssociations(r)
	if err != nil {
		return invalidMarker, err
	}
	if len(results) == 0 {
		return invalidMarker, nil
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return invalidMarker, err
	}
	queryParams := u.Query()
	limit := queryParams.Get("limit")

	// Limit is not present, only one page required
	if limit == "" {
		return invalidMarker, nil
	}

	return results[len(results)-1].ID, nil
}

// IsEmpty checks whether a RouterAssociationPage struct is empty.
func (r RouterAssociationPage) IsEmpty() (bool, error) {
	is, err := ExtractRouterAssociations(r)
	return len(is) == 0, err
}

// ExtractRouterAssociations accepts a Page struct, specifically a RouterAssociationPage struct,
// and extracts the elements into a slice of RouterAssociation structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractRouterAssociations(r pagination.Page) ([]RouterAssociation, error) {
	var s []RouterAssociation
	err := ExtractRouterAssociationsInto(r, &s)
	return s, err
}

func ExtractRouterAssociationsInto(r pagination.Page, v interface{}) error {
	return r.(RouterAssociationPage).Result.ExtractIntoSlicePtr(v, "router_associations")
}

// CreateRouterAssociationResult represents the result of a create operation. Call its Extract
// method to interpret it as a RouterAssociation.
type CreateRouterAssociationResult struct {
	commonRouterAssociationResult
}

// GetRouterAssociationResult represents the result of a get operation. Call its Extract
// method to interpret it as a RouterAssociation.
type GetRouterAssociationResult struct {
	commonRouterAssociationResult
}

// DeleteRouterAssociationResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteRouterAssociationResult struct {
	gophercloud.ErrResult
}

// UpdateRouterAssociationResult represents the result of an update operation. Call its Extract
// method to interpret it as a RouterAssociation.
type UpdateRouterAssociationResult struct {
	commonRouterAssociationResult
}

type commonPortAssociationResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a BGP VPN resource.
func (r commonPortAssociationResult) Extract() (*PortAssociation, error) {
	var s PortAssociation
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonPortAssociationResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "port_association")
}

// PortAssociation represents a BGP VPN port association object.
type PortAssociation struct {
	ID                string       `json:"id"`
	PortID            string       `json:"port_id"`
	TenantID          string       `json:"tenant_id"`
	ProjectID         string       `json:"project_id"`
	Routes            []PortRoutes `json:"routes"`
	AdvertiseFixedIPs bool         `json:"advertise_fixed_ips"`
}

// PortAssociationPage is the page returned by a pager when traversing over a
// collection of port associations.
type PortAssociationPage struct {
	pagination.MarkerPageBase
}

// NextPageURL generates the URL for the page of results after this one.
func (r PortAssociationPage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == invalidMarker {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// LastMarker returns the last offset in a ListResult.
func (r PortAssociationPage) LastMarker() (string, error) {
	results, err := ExtractPortAssociations(r)
	if err != nil {
		return invalidMarker, err
	}
	if len(results) == 0 {
		return invalidMarker, nil
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return invalidMarker, err
	}
	queryParams := u.Query()
	limit := queryParams.Get("limit")

	// Limit is not present, only one page required
	if limit == "" {
		return invalidMarker, nil
	}

	return results[len(results)-1].ID, nil
}

// IsEmpty checks whether a PortAssociationPage struct is empty.
func (r PortAssociationPage) IsEmpty() (bool, error) {
	is, err := ExtractPortAssociations(r)
	return len(is) == 0, err
}

// ExtractPortAssociations accepts a Page struct, specifically a PortAssociationPage struct,
// and extracts the elements into a slice of PortAssociation structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPortAssociations(r pagination.Page) ([]PortAssociation, error) {
	var s []PortAssociation
	err := ExtractPortAssociationsInto(r, &s)
	return s, err
}

func ExtractPortAssociationsInto(r pagination.Page, v interface{}) error {
	return r.(PortAssociationPage).Result.ExtractIntoSlicePtr(v, "port_associations")
}

// CreatePortAssociationResult represents the result of a create operation. Call its Extract
// method to interpret it as a PortAssociation.
type CreatePortAssociationResult struct {
	commonPortAssociationResult
}

// GetPortAssociationResult represents the result of a get operation. Call its Extract
// method to interpret it as a PortAssociation.
type GetPortAssociationResult struct {
	commonPortAssociationResult
}

// DeletePortAssociationResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeletePortAssociationResult struct {
	gophercloud.ErrResult
}

// UpdatePortAssociationResult represents the result of an update operation. Call its Extract
// method to interpret it as a PortAssociation.
type UpdatePortAssociationResult struct {
	commonPortAssociationResult
}
