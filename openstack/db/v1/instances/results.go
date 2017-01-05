package instances

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/openstack/db/v1/flavors"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/pagination"
)

// Volume represents information about an attached volume for a database instance.
type Volume struct {
	// The size in GB of the volume
	Size int

	Used float64
}

// Fault describes the fault reason in more detail when a database instance has errored
type Fault struct {
	// Indicates the time when the fault occured
	Created time.Time `json:"-"`
	// A message describing the fault reason
	Message string
	// More details about the fault, for example a stack trace. Only filled
	// in for admin users.
	Details string
}

// Instance represents a remote MySQL instance.
type Instance struct {
	// Indicates the datetime that the instance was created
	Created time.Time `json:"created"`

	// Indicates the most recent datetime that the instance was updated.
	Updated time.Time `json:"updated"`

	// Indicates the hardware flavor the instance uses.
	Flavor flavors.Flavor

	// A DNS-resolvable hostname associated with the database instance (rather
	// than an IPv4 address). Since the hostname always resolves to the correct
	// IP address of the database instance, this relieves the user from the task
	// of maintaining the mapping. Note that although the IP address may likely
	// change on resizing, migrating, and so forth, the hostname always resolves
	// to the correct database instance.
	Hostname string

	// Indicates the unique identifier for the instance resource.
	ID string

	// Exposes various links that reference the instance resource.
	Links []gophercloud.Link

	// The human-readable name of the instance.
	Name string

	// The build status of the instance.
	Status string

	// Fault information (only available when the instance has errored)
	Fault Fault

	// Information about the attached volume of the instance.
	Volume Volume

	// Indicates how the instance stores data.
	Datastore datastores.DatastorePartial
}

func (s *Fault) UnmarshalJSON(b []byte) error {
	type tmp Fault
	var p *struct {
		tmp
		Created string `json:"created"`
	}
	err := json.Unmarshal(b, &p)
	if err != nil {
		return err
	}
	*s = Fault(p.tmp)

	if p.Created != "" {
		s.Created, err = time.Parse(gophercloud.RFC3339NoZ, p.Created)
		if err != nil {
			return err
		}
	}

	return err
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract will extract an Instance from various result structs.
func (r commonResult) Extract() (*Instance, error) {
	var s struct {
		Instance *Instance `json:"instance"`
	}
	err := r.ExtractInto(&s)
	return s.Instance, err
}

// InstancePage represents a single page of a paginated instance collection.
type InstancePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page InstancePage) IsEmpty() (bool, error) {
	instances, err := ExtractInstances(page)
	return len(instances) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page InstancePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"instances_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractInstances will convert a generic pagination struct into a more
// relevant slice of Instance structs.
func ExtractInstances(r pagination.Page) ([]Instance, error) {
	var s struct {
		Instances []Instance `json:"instances"`
	}
	err := (r.(InstancePage)).ExtractInto(&s)
	return s.Instances, err
}

// EnableRootUserResult represents the result of an operation to enable the root user.
type EnableRootUserResult struct {
	gophercloud.Result
}

// Extract will extract root user information from a UserRootResult.
func (r EnableRootUserResult) Extract() (*users.User, error) {
	var s struct {
		User *users.User `json:"user"`
	}
	err := r.ExtractInto(&s)
	return s.User, err
}

// ActionResult represents the result of action requests, such as: restarting
// an instance service, resizing its memory allocation, and resizing its
// attached volume size.
type ActionResult struct {
	gophercloud.ErrResult
}

// IsRootEnabledResult is the result of a call to IsRootEnabled. To see if
// root is enabled, call the type's Extract method.
type IsRootEnabledResult struct {
	gophercloud.Result
}

// Extract is used to extract the data from a IsRootEnabledResult.
func (r IsRootEnabledResult) Extract() (bool, error) {
	return r.Body.(map[string]interface{})["rootEnabled"] == true, r.Err
}
