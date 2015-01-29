package services

import (
	"fmt"
	"strings"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToCDNServiceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Marker and Limit are used for pagination.
type ListOpts struct {
	Marker string `q:"marker"`
	Limit  int    `q:"limit"`
}

// ToCDNServiceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToCDNServiceListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// CDN services. It accepts a ListOpts struct, which allows for pagination via
// marker and limit.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToCDNServiceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	createPage := func(r pagination.PageResult) pagination.Page {
		p := ServicePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	pager := pagination.NewPager(c, url, createPage)
	return pager
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToCDNServiceCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// REQUIRED. Specifies the name of the service. The minimum length for name is
	// 3. The maximum length is 256.
	Name string
	// REQUIRED. Specifies a list of domains used by users to access their website.
	Domains []Domain
	// REQUIRED. Specifies a list of origin domains or IP addresses where the
	// original assets are stored.
	Origins []Origin
	// REQUIRED. Specifies the CDN provider flavor ID to use. For a list of
	// flavors, see the operation to list the available flavors. The minimum
	// length for flavor_id is 1. The maximum length is 256.
	FlavorID string
	// OPTIONAL. Specifies the TTL rules for the assets under this service. Supports wildcards for fine-grained control.
	Caching []CacheRule
	// OPTIONAL. Specifies the restrictions that define who can access assets (content from the CDN cache).
	Restrictions []Restriction
}

// ToCDNServiceCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToCDNServiceCreateMap() (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if opts.Name == "" {
		return nil, no("Name")
	}
	s["name"] = opts.Name

	if opts.Domains == nil {
		return nil, no("Domains")
	}
	for _, domain := range opts.Domains {
		if domain.Domain == "" {
			return nil, no("Domains[].Domain")
		}
	}
	s["domains"] = opts.Domains

	if opts.Origins == nil {
		return nil, no("Origins")
	}
	for _, origin := range opts.Origins {
		if origin.Origin == "" {
			return nil, no("Origins[].Origin")
		}
		if origin.Rules == nil && len(opts.Origins) > 1 {
			return nil, no("Origins[].Rules")
		}
		for _, rule := range origin.Rules {
			if rule.Name == "" {
				return nil, no("Origins[].Rules[].Name")
			}
			if rule.RequestURL == "" {
				return nil, no("Origins[].Rules[].RequestURL")
			}
		}
	}
	s["origins"] = opts.Origins

	if opts.FlavorID == "" {
		return nil, no("FlavorID")
	}
	s["flavor_id"] = opts.FlavorID

	if opts.Caching != nil {
		for _, cache := range opts.Caching {
			if cache.Name == "" {
				return nil, no("Caching[].Name")
			}
			if cache.Rules != nil {
				for _, rule := range cache.Rules {
					if rule.Name == "" {
						return nil, no("Caching[].Rules[].Name")
					}
					if rule.RequestURL == "" {
						return nil, no("Caching[].Rules[].RequestURL")
					}
				}
			}
		}
		s["caching"] = opts.Caching
	}

	if opts.Restrictions != nil {
		for _, restriction := range opts.Restrictions {
			if restriction.Name == "" {
				return nil, no("Restrictions[].Name")
			}
			if restriction.Rules != nil {
				for _, rule := range restriction.Rules {
					if rule.Name == "" {
						return nil, no("Restrictions[].Rules[].Name")
					}
				}
			}
		}
		s["restrictions"] = opts.Restrictions
	}

	return s, nil
}

// Create accepts a CreateOpts struct and creates a new CDN service using the
// values provided.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToCDNServiceCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	// Send request to API
	resp, err := perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{202},
	})
	res.Header = resp.HttpResponse.Header
	res.Err = err
	return res
}

// Get retrieves a specific service based on its URL or its unique ID. For
// example, both "96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0" and
// "https://global.cdn.api.rackspacecloud.com/v1.0/services/96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0"
// are valid options for idOrURL.
func Get(c *gophercloud.ServiceClient, idOrURL string) GetResult {
	var url string
	if strings.Contains(idOrURL, "/") {
		url = idOrURL
	} else {
		url = getURL(c, idOrURL)
	}

	var res GetResult
	_, res.Err = perigee.Request("GET", url, perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToCDNServiceUpdateMap() ([]map[string]interface{}, error)
}

// Op represents an update operation.
type Op string

var (
	// Add is a constant used for performing a "add" operation when updating.
	Add Op = "add"
	// Remove is a constant used for performing a "remove" operation when updating.
	Remove Op = "remove"
	// Replace is a constant used for performing a "replace" operation when updating.
	Replace Op = "replace"
)

// Path is a JSON pointer location that indicates which service parameter is being added, replaced,
// or removed.
type Path struct {
	baseElement string
}

func (p Path) renderDash() string {
	return fmt.Sprintf("/%s/-", p.baseElement)
}

func (p Path) renderIndex(index int64) string {
	return fmt.Sprintf("/%s/%d", p.baseElement, index)
}

var (
	// PathDomains indicates that an update operation is to be performed on a Domain.
	PathDomains = Path{baseElement: "domains"}

	// PathOrigins indicates that an update operation is to be performed on an Origin.
	PathOrigins = Path{baseElement: "origins"}

	// PathCaching indicates that an update operation is to be performed on a CacheRule.
	PathCaching = Path{baseElement: "caching"}
)

// UpdateOpts represents the attributes used when updating an existing CDN service.
type UpdateOpts []UpdateOpt

// UpdateOpt represents a single update to an existing service. Multiple updates
// to a service can be submitted at the same time. See UpdateOpts.
type UpdateOpt struct {
	// Specifies the update operation to perform.
	Op Op `json:"op"`
	// Specifies the JSON Pointer location within the service's JSON representation
	// of the service parameter being added, replaced or removed.
	Path string `json:"path"`
	// Specifies the actual value to be added or replaced. It is not required for
	// the remove operation.
	Value map[string]interface{} `json:"value,omitempty"`
}

type value interface {
	toPatchValue() map[string]interface{}
	appropriatePath() Path
}

// Patch represents a single update to an existing Service. Multiple updates to a service can be
// submitted at the same time.
type Patch interface {
	ToCDNServiceUpdateMap() map[string]interface{}
}

// Addition is a Patch that requests the addition of one or more values (Domains, Origins, or
// CacheRules) to a Service. Pass it to the Update function as part of the Patch slice.
type Addition struct {
	Value value
}

// ToCDNServiceUpdateMap converts an Addition into a request body fragment suitable for the
// Update call.
func (a Addition) ToCDNServiceUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"op":    "add",
		"path":  a.Value.appropriatePath().renderDash(),
		"value": a.Value.toPatchValue(),
	}
}

// Replacement is a Patch that alters a specific service parameter (Domain, Origin, or CacheRule)
// in-place by index. Pass it to the Update function as part of the Patch slice.
type Replacement struct {
	Value value
	Index int64
}

// ToCDNServiceUpdateMap converts a Replacement into a request body fragment suitable for the
// Update call.
func (r Replacement) ToCDNServiceUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"op":    "replace",
		"path":  r.Value.appropriatePath().renderIndex(r.Index),
		"value": r.Value.toPatchValue(),
	}
}

// Removal is a Patch that requests the removal of a service parameter (Domain, Origin, or
// CacheRule) by index. Pass it to the Update function as part of the Patch slice.
type Removal struct {
	Path  Path
	Index int64
}

// ToCDNServiceUpdateMap converts a Removal into a request body fragment suitable for the
// Update call.
func (r Removal) ToCDNServiceUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"op":   "remove",
		"path": r.Path.renderIndex(r.Index),
	}
}

// ToCDNServiceUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToCDNServiceUpdateMap() ([]map[string]interface{}, error) {
	s := make([]map[string]interface{}, len(opts))

	for i, opt := range opts {
		if opt.Op != Add && opt.Op != Remove && opt.Op != Replace {
			return nil, fmt.Errorf("Invalid Op: %v", opt.Op)
		}
		if opt.Op == "" {
			return nil, no("Op")
		}
		if opt.Path == "" {
			return nil, no("Path")
		}
		if opt.Op != Remove && opt.Value == nil {
			return nil, no("Value")
		}
		s[i] = map[string]interface{}{
			"op":    opt.Op,
			"path":  opt.Path,
			"value": opt.Value,
		}
	}

	return s, nil
}

// Update accepts a UpdateOpts struct and updates an existing CDN service using
// the values provided. idOrURL can be either the service's URL or its ID. For
// example, both "96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0" and
// "https://global.cdn.api.rackspacecloud.com/v1.0/services/96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0"
// are valid options for idOrURL.
func Update(c *gophercloud.ServiceClient, idOrURL string, opts UpdateOptsBuilder) UpdateResult {
	var url string
	if strings.Contains(idOrURL, "/") {
		url = idOrURL
	} else {
		url = updateURL(c, idOrURL)
	}

	var res UpdateResult
	reqBody, err := opts.ToCDNServiceUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	resp, err := perigee.Request("PATCH", url, perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{202},
	})
	res.Header = resp.HttpResponse.Header
	res.Err = err
	return res
}

// Delete accepts a service's ID or its URL and deletes the CDN service
// associated with it. For example, both "96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0" and
// "https://global.cdn.api.rackspacecloud.com/v1.0/services/96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0"
// are valid options for idOrURL.
func Delete(c *gophercloud.ServiceClient, idOrURL string) DeleteResult {
	var url string
	if strings.Contains(idOrURL, "/") {
		url = idOrURL
	} else {
		url = deleteURL(c, idOrURL)
	}

	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return res
}
