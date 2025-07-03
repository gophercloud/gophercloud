package portgroups

import (
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ResourceLink represents a link with href and rel attributes
type ResourceLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// PortGroup represents a port group in the baremetal service
// https://docs.openstack.org/api-ref/baremetal/#portgroups-portgroups
type PortGroup struct {
	// Human-readable identifier for the Portgroup resource. May be undefined.
	Name string `json:"name"`

	// The UUID for the resource.
	UUID string `json:"uuid"`

	// Physical hardware address of this Portgroup, typically the hardware MAC address.
	Address string `json:"address,omitempty"`

	// UUID of the Node this resource belongs to.
	NodeUUID string `json:"node_uuid"`

	// Indicates whether ports that are members of this portgroup can be used as
	// stand-alone ports.
	StandalonePortsSupported bool `json:"standalone_ports_supported"`

	// Internal metadata set and stored by the Portgroup. This field is read-only.
	InternalInfo map[string]any `json:"internal_info"`

	// A set of one or more arbitrary metadata key and value pairs.
	Extra map[string]any `json:"extra"`

	// Mode of the port group. For possible values, refer to
	// https://www.kernel.org/doc/Documentation/networking/bonding.txt
	Mode string `json:"mode"`

	// Key/value properties related to the port group's configuration.
	Properties map[string]any `json:"properties"`

	// The UTC date and time when the resource was created, ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`

	// The UTC date and time when the resource was updated, ISO 8601 format.
	// May be "null".
	UpdatedAt time.Time `json:"updated_at"`

	// A list of relative links. Includes the self and bookmark links.
	Links []ResourceLink `json:"links"`

	// Links to the collection of ports belonging to this portgroup.
	Ports []ResourceLink `json:"ports"`
}

type portgroupsResult struct {
	gophercloud.Result
}

func (r portgroupsResult) Extract() (*PortGroup, error) {
	var s PortGroup
	err := r.ExtractInto(&s)
	return &s, err
}

func (r portgroupsResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "")
}

func ExtractPortGroupsInto(r pagination.Page, v any) error {
	return r.(PortGroupsPage).ExtractIntoSlicePtr(v, "portgroups")
}

// PortGroupsPage abstracts the raw results of making a List() request against
// the API.
type PortGroupsPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no PortGroup results.
func (r PortGroupsPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	s, err := ExtractPortGroups(r)
	return len(s) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r PortGroupsPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"portgroups_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractPortGroups interprets the results of a single page from a List() call,
// producing a slice of PortGroup entities.
func ExtractPortGroups(r pagination.Page) ([]PortGroup, error) {
	var s []PortGroup
	err := ExtractPortGroupsInto(r, &s)
	return s, err
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a PortGroup.
type GetResult struct {
	portgroupsResult
}

// CreateResult is the response from a Create operation.
type CreateResult struct {
	portgroupsResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
