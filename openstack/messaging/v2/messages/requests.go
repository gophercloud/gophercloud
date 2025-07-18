package messages

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToMessageListQuery() (string, error)
}

// ListOpts params to be used with List.
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of queues
	Limit int `q:"limit,omitempty"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker,omitempty"`

	// Indicate if the messages can be echoed back to the client that posted them.
	Echo bool `q:"echo,omitempty"`

	// Indicate if the messages list should include the claimed messages.
	IncludeClaimed bool `q:"include_claimed,omitempty"`

	//Indicate if the messages list should include the delayed messages.
	IncludeDelayed bool `q:"include_delayed,omitempty"`
}

func (opts ListOpts) ToMessageListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListMessages lists messages on a specific queue based off queue name.
func List(client *gophercloud.ServiceClient, queueName string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, queueName)
	if opts != nil {
		query, err := opts.ToMessageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pager := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return MessagePage{
			serviceURL:     client.ServiceURL(),
			LinkedPageBase: pagination.LinkedPageBase{PageResult: r},
		}
	})
	return pager
}

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToMessageCreateMap() (map[string]any, error)
}

// BatchCreateOpts is an array of CreateOpts.
type BatchCreateOpts []CreateOpts

// CreateOpts params to be used with Create.
type CreateOpts struct {
	// TTL specifies how long the server waits before marking the message
	// as expired and removing it from the queue.
	TTL int `json:"ttl,omitempty"`

	// Delay specifies how long the message can be claimed.
	Delay int `json:"delay,omitempty"`

	// Body specifies an arbitrary document that constitutes the body of the message being sent.
	Body map[string]any `json:"body" required:"true"`
}

// ToMessageCreateMap constructs a request body from BatchCreateOpts.
func (opts BatchCreateOpts) ToMessageCreateMap() (map[string]any, error) {
	messages := make([]map[string]any, len(opts))
	for i, message := range opts {
		messageMap, err := message.ToMap()
		if err != nil {
			return nil, err
		}
		messages[i] = messageMap
	}
	return map[string]any{"messages": messages}, nil
}

// ToMap constructs a request body from UpdateOpts.
func (opts CreateOpts) ToMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a message on a specific queue based of off queue name.
func Create(ctx context.Context, client *gophercloud.ServiceClient, queueName string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToMessageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client, queueName), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMessagesOptsBuilder allows extensions to add additional parameters to the
// DeleteMessages request.
type DeleteMessagesOptsBuilder interface {
	ToMessagesDeleteQuery() (string, error)
}

// DeleteMessagesOpts params to be used with DeleteMessages.
type DeleteMessagesOpts struct {
	IDs []string `q:"ids,omitempty"`
}

// ToMessagesDeleteQuery formats a DeleteMessagesOpts structure into a query string.
func (opts DeleteMessagesOpts) ToMessagesDeleteQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// DeleteMessages deletes multiple messages based off of ID.
func DeleteMessages(ctx context.Context, client *gophercloud.ServiceClient, queueName string, opts DeleteMessagesOptsBuilder) (r DeleteResult) {
	url := deleteURL(client, queueName)
	if opts != nil {
		query, err := opts.ToMessagesDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		OkCodes: []int{200, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// PopMessagesOptsBuilder allows extensions to add additional parameters to the
// DeleteMessages request.
type PopMessagesOptsBuilder interface {
	ToMessagesPopQuery() (string, error)
}

// PopMessagesOpts params to be used with PopMessages.
type PopMessagesOpts struct {
	Pop int `q:"pop,omitempty"`
}

// ToMessagesPopQuery formats a PopMessagesOpts structure into a query string.
func (opts PopMessagesOpts) ToMessagesPopQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// PopMessages deletes and returns multiple messages based off of number of messages.
func PopMessages(ctx context.Context, client *gophercloud.ServiceClient, queueName string, opts PopMessagesOptsBuilder) (r PopResult) {
	url := deleteURL(client, queueName)
	if opts != nil {
		query, err := opts.ToMessagesPopQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		JSONResponse: &r.Body,
		OkCodes:      []int{200, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMessagesOptsBuilder allows extensions to add additional parameters to the
// GetMessages request.
type GetMessagesOptsBuilder interface {
	ToGetMessagesListQuery() (string, error)
}

// GetMessagesOpts params to be used with GetMessages.
type GetMessagesOpts struct {
	IDs []string `q:"ids,omitempty"`
}

// ToGetMessagesListQuery formats a GetMessagesOpts structure into a query string.
func (opts GetMessagesOpts) ToGetMessagesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// GetMessages requests details on a multiple messages, by IDs.
func GetMessages(ctx context.Context, client *gophercloud.ServiceClient, queueName string, opts GetMessagesOptsBuilder) (r GetMessagesResult) {
	url := getURL(client, queueName)
	if opts != nil {
		query, err := opts.ToGetMessagesListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get requests details on a single message, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, queueName string, messageID string) (r GetResult) {
	resp, err := client.Get(ctx, messageURL(client, queueName, messageID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteOptsBuilder allows extensions to add additional parameters to the
// delete request.
type DeleteOptsBuilder interface {
	ToMessageDeleteQuery() (string, error)
}

// DeleteOpts params to be used with Delete.
type DeleteOpts struct {
	// ClaimID instructs Delete to delete a message that is associated with a claim ID
	ClaimID string `q:"claim_id,omitempty"`
}

// ToMessageDeleteQuery formats a DeleteOpts structure into a query string.
func (opts DeleteOpts) ToMessageDeleteQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Delete deletes a specific message from the queue.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, queueName string, messageID string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := DeleteMessageURL(client, queueName, messageID)
	if opts != nil {
		query, err := opts.ToMessageDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
