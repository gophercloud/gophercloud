package impliedroles

import (
	"net/url"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToImpliedRoleListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Name filters the response by role name.
	Name string `q:"name"`

	// Filters filters the response by custom filters such as
	// 'name__contains=foo'
	Filters map[string]string `q:"-"`
}

// ToImpliedRoleListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImpliedRoleListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()
	for k, v := range opts.Filters {
		i := strings.Index(k, "__")
		if i > 0 && i < len(k)-2 {
			params.Add(k, v)
		} else {
			return "", InvalidListFilter{FilterName: k}
		}
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

// List enumerates the ApplicationCredentials to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToImpliedRoleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ImpliedRolePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Create Implied roles.
func Create(client *gophercloud.ServiceClient, pirorRoleId string, impliedRoleID string) (r createImpliedRoleResult) {

	resp, err := client.Put(createURL(client, pirorRoleId, impliedRoleID), nil, r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a implied roles.
func Delete(client *gophercloud.ServiceClient, pirorRoleId string, impliesRoleId string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, pirorRoleId, impliesRoleId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
