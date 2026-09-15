// Package fakeservice is a hand-written fixture mirroring the Gophercloud
// request/URL/result patterns the impl extractor must understand: a
// BuildRequestBody struct body, a q: query-opts builder, a json:"-" field, a
// manual map[string]any action body, and nested url-func resolution.
package fakeservice

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder is the query-string builder interface.
type ListOptsBuilder interface {
	ToWidgetListQuery() (string, error)
}

// ListOpts exercises q: tags, including a q:"-" field that must be ignored.
type ListOpts struct {
	Name     string `q:"name"`
	Limit    int    `q:"limit"`
	Internal string `q:"-"`
}

func (opts ListOpts) ToWidgetListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List builds a pagination.Pager (the list pattern the extractor detects via
// pagination.NewPager rather than a direct verb call).
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToWidgetListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return WidgetPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is the request-body builder interface.
type CreateOptsBuilder interface {
	ToWidgetCreateMap() (map[string]any, error)
}

// CreateOpts exercises a BuildRequestBody struct body with a json:"-" field
// (Secret) that must be reported as ManualHandled, not a gap.
type CreateOpts struct {
	Name        string `json:"name" required:"true"`
	Description string `json:"description,omitempty"`
	Size        int    `json:"size,omitempty"`
	Secret      string `json:"-"`
}

func (opts CreateOpts) ToWidgetCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "widget")
}

// Create posts a BuildRequestBody-derived body.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToWidgetCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete exercises nested url-func resolution (deleteURL -> resourceURL).
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Reboot exercises a manual map[string]any action body: the single key is the
// action name.
func Reboot(ctx context.Context, c *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]any{"reboot": nil}
	resp, err := c.Post(ctx, actionURL(c, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
