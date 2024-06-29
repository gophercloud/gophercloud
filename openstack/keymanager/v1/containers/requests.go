package containers

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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
func List(client gophercloud.Client, opts ListOptsBuilder) pagination.Pager {
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
func Get(ctx context.Context, client gophercloud.Client, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToContainerCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a container.
type CreateOpts struct {
	// Type represents the type of container.
	Type ContainerType `json:"type" required:"true"`

	// Name is the name of the container.
	Name string `json:"name"`

	// SecretRefs is a list of secret refs for the container.
	SecretRefs []SecretRef `json:"secret_refs,omitempty"`
}

// ToContainerCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToContainerCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new container.
func Create(ctx context.Context, client gophercloud.Client, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToContainerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a container.
func Delete(ctx context.Context, client gophercloud.Client, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListConsumersOptsBuilder allows extensions to add additional parameters to
// the ListConsumers request
type ListConsumersOptsBuilder interface {
	ToContainerListConsumersQuery() (string, error)
}

// ListConsumersOpts provides options to filter the List results.
type ListConsumersOpts struct {
	// Limit is the amount of consumers to retrieve.
	Limit int `q:"limit"`

	// Offset is the index within the list to retrieve.
	Offset int `q:"offset"`
}

// ToContainerListConsumersQuery formats a ListConsumersOpts into a query
// string.
func (opts ListOpts) ToContainerListConsumersQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListConsumers retrieves a list of consumers from a container.
func ListConsumers(client gophercloud.Client, containerID string, opts ListConsumersOptsBuilder) pagination.Pager {
	url := listConsumersURL(client, containerID)
	if opts != nil {
		query, err := opts.ToContainerListConsumersQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ConsumerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateConsumerOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateConsumerOptsBuilder interface {
	ToContainerConsumerCreateMap() (map[string]any, error)
}

// CreateConsumerOpts provides options used to create a container.
type CreateConsumerOpts struct {
	// Name is the name of the consumer.
	Name string `json:"name"`

	// URL is the URL to the consumer resource.
	URL string `json:"URL"`
}

// ToContainerConsumerCreateMap formats a CreateConsumerOpts into a create
// request.
func (opts CreateConsumerOpts) ToContainerConsumerCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// CreateConsumer creates a new consumer.
func CreateConsumer(ctx context.Context, client gophercloud.Client, containerID string, opts CreateConsumerOptsBuilder) (r CreateConsumerResult) {
	b, err := opts.ToContainerConsumerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createConsumerURL(client, containerID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteConsumerOptsBuilder allows extensions to add additional parameters to
// the Delete request.
type DeleteConsumerOptsBuilder interface {
	ToContainerConsumerDeleteMap() (map[string]any, error)
}

// DeleteConsumerOpts represents options used for deleting a consumer.
type DeleteConsumerOpts struct {
	// Name is the name of the consumer.
	Name string `json:"name"`

	// URL is the URL to the consumer resource.
	URL string `json:"URL"`
}

// ToContainerConsumerDeleteMap formats a DeleteConsumerOpts into a create
// request.
func (opts DeleteConsumerOpts) ToContainerConsumerDeleteMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// DeleteConsumer deletes a consumer.
func DeleteConsumer(ctx context.Context, client gophercloud.Client, containerID string, opts DeleteConsumerOptsBuilder) (r DeleteConsumerResult) {
	url := deleteConsumerURL(client, containerID)

	b, err := opts.ToContainerConsumerDeleteMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Request(ctx, "DELETE", url, &gophercloud.RequestOpts{
		JSONBody:     b,
		JSONResponse: &r.Body,
		OkCodes:      []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SecretRefBuilder allows extensions to add additional parameters to the
// Create request.
type SecretRefBuilder interface {
	ToContainerSecretRefMap() (map[string]any, error)
}

// ToContainerSecretRefMap formats a SecretRefBuilder into a create
// request.
func (opts SecretRef) ToContainerSecretRefMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// CreateSecret creates a new consumer.
func CreateSecretRef(ctx context.Context, client gophercloud.Client, containerID string, opts SecretRefBuilder) (r CreateSecretRefResult) {
	b, err := opts.ToContainerSecretRefMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createSecretRefURL(client, containerID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteSecret deletes a consumer.
func DeleteSecretRef(ctx context.Context, client gophercloud.Client, containerID string, opts SecretRefBuilder) (r DeleteSecretRefResult) {
	url := deleteSecretRefURL(client, containerID)

	b, err := opts.ToContainerSecretRefMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Request(ctx, "DELETE", url, &gophercloud.RequestOpts{
		JSONBody: b,
		OkCodes:  []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
