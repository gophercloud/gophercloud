package containers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ContainerType represents the valid types of containers.
type ContainerType string

const (
	GenericContainer     ContainerType = "generic"
	RSAContainer         ContainerType = "rsa"
	CertificateContainer ContainerType = "certificate"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToContainerListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Limit is the amount of containers to retrieve.
	Limit int `q:"limit"`

	// Name is the name of the container
	Name string `q:"name"`

	// Offset is the index within the list to retrieve.
	Offset int `q:"offset"`
}

// ToContainerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToContainerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List retrieves a list of containers.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToContainerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ContainerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details of a container.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToContainerCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a container.
type CreateOpts struct {
	// Type represents the type of container.
	Type ContainerType `json:"type" required:"true"`

	// Name is the name of the container.
	Name string `json:"name"`

	// SecretRefs is a list of secret refs for the container.
	SecretRefs []SecretRef `json:"secret_refs"`
}

// ToContainerCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToContainerCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new container.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToContainerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// Delete deletes a container.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
