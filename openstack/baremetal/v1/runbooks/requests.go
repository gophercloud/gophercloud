package runbooks

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToRunbookListQuery() (string, error)
}

// ListOpts allows the sorting and pagination of runbook collections.
type ListOpts struct {
	// One or more fields to be returned in the response.
	Fields []string `q:"fields" format:"comma-separated"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`

	// Sorts the response by the requested sort direction.
	// Valid value is asc (ascending) or desc (descending). Default is asc.
	SortDir string `q:"sort_dir"`

	// Sorts the response by this attribute value. Default is id.
	SortKey string `q:"sort_key"`

	// Detail returns every runbook field, not only the UUID and name.
	Detail bool `q:"detail"`
}

// ToRunbookListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRunbookListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list runbooks accessible to you.
// This requires the client to be set to microversion 1.92 or later.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToRunbookListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RunbookPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get requests the details of a runbook by UUID or name.
// This requires the client to be set to microversion 1.92 or later.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Step is a single step run by a runbook.
type Step struct {
	// Interface is the driver interface that runs the step.
	Interface nodes.StepInterface `json:"interface" required:"true"`

	// Step is the name of the step method on the driver interface.
	Step string `json:"step" required:"true"`

	// Args are the keyword arguments passed to the step method.
	Args map[string]any `json:"args"`

	// Order is the non-negative position of the step in the runbook.
	Order int `json:"order"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRunbookCreateMap() (map[string]any, error)
}

// CreateOpts specifies runbook creation parameters.
type CreateOpts struct {
	// Name is the unique name of the runbook. Before microversion 1.112
	// it must be a valid trait name, such as CUSTOM_EXAMPLE.
	Name string `json:"name" required:"true"`

	// Steps are the steps the runbook runs. At least one is required.
	Steps []Step `json:"steps" required:"true"`

	// UUID is the UUID for the runbook. Ironic generates one if unset.
	UUID string `json:"uuid,omitempty"`

	// DisableRamdisk runs the steps without booting the IPA ramdisk.
	DisableRamdisk *bool `json:"disable_ramdisk,omitempty"`

	// Extra is a set of arbitrary metadata key and value pairs.
	Extra map[string]any `json:"extra,omitempty"`

	// Public makes the runbook available to every project.
	// It must not be true when Owner is set.
	Public *bool `json:"public,omitempty"`

	// Owner is the project that owns the runbook.
	Owner string `json:"owner,omitempty"`

	// Description is a free-text description of the runbook.
	// This requires the client to be set to microversion 1.112 or later.
	Description string `json:"description,omitempty"`
}

// ToRunbookCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToRunbookCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create requests the creation of a runbook.
// This requires the client to be set to microversion 1.92 or later.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToRunbookCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, rootURL(client), reqBody, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Patch is a single JSON patch operation for an Update request.
type Patch interface {
	ToRunbookUpdateMap() map[string]any
}

// UpdateOpts is a slice of Patches used to update a runbook.
type UpdateOpts []Patch

// UpdateOp is a JSON patch operation type.
type UpdateOp string

const (
	ReplaceOp UpdateOp = "replace"
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
)

// UpdateOperation is a JSON patch operation on a runbook field.
type UpdateOperation struct {
	// Op is the patch operation to apply.
	Op UpdateOp `json:"op" required:"true"`

	// Path is the JSON pointer to the field, such as /name.
	Path string `json:"path" required:"true"`

	// Value is the new value. It is ignored by remove operations.
	Value any `json:"value,omitempty"`
}

// ToRunbookUpdateMap assembles a JSON patch operation.
func (opts UpdateOperation) ToRunbookUpdateMap() map[string]any {
	return map[string]any{
		"op":    opts.Op,
		"path":  opts.Path,
		"value": opts.Value,
	}
}

// Update requests the update of a runbook by UUID or name.
// This requires the client to be set to microversion 1.92 or later.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	body := make([]map[string]any, len(opts))
	for i, patch := range opts {
		body[i] = patch.ToRunbookUpdateMap()
	}

	resp, err := client.Patch(ctx, resourceURL(client, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete requests the deletion of a runbook by UUID or name.
// This requires the client to be set to microversion 1.92 or later.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListTraits requests the traits of a runbook. A node can run the runbook
// only if it has one of these traits.
// This requires the client to be set to microversion 1.112 or later.
func ListTraits(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ListTraitsResult) {
	resp, err := client.Get(ctx, traitsURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SetTraits replaces every trait of a runbook with traits.
// This requires the client to be set to microversion 1.112 or later.
func SetTraits(ctx context.Context, client *gophercloud.ServiceClient, id string, traits []string) (r SetTraitsResult) {
	if traits == nil {
		traits = []string{}
	}
	body := map[string]any{"traits": traits}

	resp, err := client.Put(ctx, traitsURL(client, id), body, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AddTrait adds a single trait to a runbook.
// This requires the client to be set to microversion 1.112 or later.
func AddTrait(ctx context.Context, client *gophercloud.ServiceClient, id string, trait string) (r AddTraitResult) {
	resp, err := client.Put(ctx, traitURL(client, id, trait), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteTraits removes every trait from a runbook.
// This requires the client to be set to microversion 1.112 or later.
func DeleteTraits(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteTraitsResult) {
	resp, err := client.Delete(ctx, traitsURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteTrait removes a single trait from a runbook. Removing a trait the
// runbook does not have succeeds.
// This requires the client to be set to microversion 1.112 or later.
func DeleteTrait(ctx context.Context, client *gophercloud.ServiceClient, id string, trait string) (r DeleteTraitResult) {
	resp, err := client.Delete(ctx, traitURL(client, id, trait), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
