package replicas

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

const (
	invalidMarker = "-1"
)

// Replica contains all information associated with an OpenStack Share Replica.
type Replica struct {
	// ID of the share replica
	ID string `json:"id"`
	// The availability zone of the share replica.
	AvailabilityZone string `json:"availability_zone"`
	// Indicates whether existing access rules will be cast to read/only.
	CastRulesToReadonly bool `json:"cast_rules_to_readonly"`
	// The host name of the share replica.
	Host string `json:"host"`
	// The UUID of the share to which a share replica belongs.
	ShareID string `json:"share_id"`
	// The UUID of the share network where the resource is exported to.
	ShareNetworkID string `json:"share_network_id"`
	// The UUID of the share server.
	ShareServerID string `json:"share_server_id"`
	// The share replica status.
	Status string `json:"status"`
	// The share replica state.
	State string `json:"replica_state"`
	// Timestamp when the replica was created.
	CreatedAt time.Time `json:"-"`
	// Timestamp when the replica was updated.
	UpdatedAt time.Time `json:"-"`
}

func (r *Replica) UnmarshalJSON(b []byte) error {
	type tmp Replica
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Replica(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Replica object from the commonResult.
func (r commonResult) Extract() (*Replica, error) {
	var s struct {
		Replica *Replica `json:"share_replica"`
	}
	err := r.ExtractInto(&s)
	return s.Replica, err
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// ReplicaPage is a pagination.pager that is returned from a call to the List function.
type ReplicaPage struct {
	pagination.MarkerPageBase
}

// NextPageURL generates the URL for the page of results after this one.
func (r ReplicaPage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == invalidMarker {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("offset", mark)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// LastMarker returns the last offset in a ListResult.
func (r ReplicaPage) LastMarker() (string, error) {
	replicas, err := ExtractReplicas(r)
	if err != nil {
		return invalidMarker, err
	}
	if len(replicas) == 0 {
		return invalidMarker, nil
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return invalidMarker, err
	}
	queryParams := u.Query()
	offset := queryParams.Get("offset")
	limit := queryParams.Get("limit")

	// Limit is not present, only one page required
	if limit == "" {
		return invalidMarker, nil
	}

	iOffset := 0
	if offset != "" {
		iOffset, err = strconv.Atoi(offset)
		if err != nil {
			return invalidMarker, err
		}
	}
	iLimit, err := strconv.Atoi(limit)
	if err != nil {
		return invalidMarker, err
	}
	iOffset = iOffset + iLimit
	offset = strconv.Itoa(iOffset)

	return offset, nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface.
func (r ReplicaPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	replicas, err := ExtractReplicas(r)
	return len(replicas) == 0, err
}

// ExtractReplicas extracts and returns Replicas. It is used while iterating
// over a replicas.List or replicas.ListDetail calls.
func ExtractReplicas(r pagination.Page) ([]Replica, error) {
	var s []Replica
	err := ExtractReplicasInto(r, &s)
	return s, err
}

// ExtractReplicasInto similar to ExtractReplicas but operates on a `list` of
// replicas.
func ExtractReplicasInto(r pagination.Page, v any) error {
	return r.(ReplicaPage).Result.ExtractIntoSlicePtr(v, "share_replicas")
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// ListExportLocationsResult contains the result body and error from a
// ListExportLocations request.
type ListExportLocationsResult struct {
	gophercloud.Result
}

// GetExportLocationResult contains the result body and error from a
// GetExportLocation request.
type GetExportLocationResult struct {
	gophercloud.Result
}

// ExportLocation contains all information associated with a share export location
type ExportLocation struct {
	// The share replica export location UUID.
	ID string `json:"id"`
	// The export location path that should be used for mount operation.
	Path string `json:"path"`
	// The UUID of the share instance that this export location belongs to.
	ShareInstanceID string `json:"share_instance_id"`
	// Defines purpose of an export location. If set to true, then it is
	// expected to be used for service needs and by administrators only. If
	// it is set to false, then this export location can be used by end users.
	IsAdminOnly bool `json:"is_admin_only"`
	// Drivers may use this field to identify which export locations are
	// most efficient and should be used preferentially by clients.
	// By default it is set to false value. New in version 2.14.
	Preferred bool `json:"preferred"`
	// The availability zone of the share replica.
	AvailabilityZone string `json:"availability_zone"`
	// The share replica state.
	State string `json:"replica_state"`
	// Timestamp when the export location was created.
	CreatedAt time.Time `json:"-"`
	// Timestamp when the export location was updated.
	UpdatedAt time.Time `json:"-"`
}

func (r *ExportLocation) UnmarshalJSON(b []byte) error {
	type tmp ExportLocation
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ExportLocation(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

// Extract will get the Export Locations from the ListExportLocationsResult
func (r ListExportLocationsResult) Extract() ([]ExportLocation, error) {
	var s struct {
		ExportLocations []ExportLocation `json:"export_locations"`
	}
	err := r.ExtractInto(&s)
	return s.ExportLocations, err
}

// Extract will get the Export Location from the GetExportLocationResult
func (r GetExportLocationResult) Extract() (*ExportLocation, error) {
	var s struct {
		ExportLocation *ExportLocation `json:"export_location"`
	}
	err := r.ExtractInto(&s)
	return s.ExportLocation, err
}

// PromoteResult contains the error from an Promote request.
type PromoteResult struct {
	gophercloud.ErrResult
}

// ResyncResult contains the error from a Resync request.
type ResyncResult struct {
	gophercloud.ErrResult
}

// ResetStatusResult contains the error from a ResetStatus request.
type ResetStatusResult struct {
	gophercloud.ErrResult
}

// ResetStateResult contains the error from a ResetState request.
type ResetStateResult struct {
	gophercloud.ErrResult
}

// ForceDeleteResult contains the error from a ForceDelete request.
type ForceDeleteResult struct {
	gophercloud.ErrResult
}
