package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// LIST

type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

type ListOpts struct {
	ID   string `q:"id"`
	Name string `q:"name"`
}

func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CREATE

type CreateOptsBuilder interface {
	ToFlavorCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	FlavorProfileId string `json:"flavor_profile_id,required:"true""`
	Enabled         bool   `json:"enabled,default:"true""`
}

func (opts CreateOpts) ToFlavorCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "flavor")
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFlavorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GET

func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UPDATE

type UpdateOptsBuilder interface {
	ToFlavorUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled" default:"true"`
}

func (opts UpdateOpts) ToFlavorUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "flavor")
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToFlavorUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
