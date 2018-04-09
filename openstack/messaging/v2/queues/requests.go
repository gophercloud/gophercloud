package queues

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToQueueListQuery() (string, error)
}

// ListOpts params to be used with List
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of queues
	Limit int `q:"limit,omitempty"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker,omitempty"`

	// Specifies if showing the detailed information when querying queues
	Detailed bool `q:"detailed,omitempty"`
}

// ToQueueListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToQueueListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of queues.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToQueueListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pager := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return QueuePage{pagination.LinkedPageBase{PageResult: r}}

	})
	return pager
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToQueueCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies the queue creation parameters.
type CreateOpts struct {
	// The name of the queue to create.
	QueueName string `json:"queue_name" required:"true"`

	// The target incoming messages will be moved to when a message can’t
	// processed successfully after meet the max claim count is met.
	DeadLetterQueue string `json:"_dead_letter_queue,omitempty"`

	// The new TTL setting for messages when moved to dead letter queue.
	DeadLetterQueueMessagesTTL int `json:"_dead_letter_queue_messages_ttl,omitempty"`

	// The delay of messages defined for a queue. When the messages send to
	// the queue, it will be delayed for some times and means it can not be
	// claimed until the delay expired.
	DefaultMessageDelay int `json:"_default_message_delay,omitempty"`

	// The default TTL of messages defined for a queue, which will effect for
	// any messages posted to the queue.
	DefaultMessageTTL int `json:"_default_message_ttl" required:"true"`

	// The flavor name which can tell Zaqar which storage pool will be used
	// to create the queue.
	Flavor string `json:"_flavor,omitempty"`

	// The max number the message can be claimed.
	MaxClaimCount int `json:"_max_claim_count,omitempty"`

	// The max post size of messages defined for a queue, which will effect
	// for any messages posted to the queue.
	MaxMessagesPostSize int `json:"_max_messages_post_size,omitempty"`

	// Extra is free-form extra key/value pairs to describe the queue.
	Extra map[string]interface{} `json:"-"`
}

// ToQueueCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToQueueCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		for key, value := range opts.Extra {
			b[key] = value
		}

	}
	return b, nil
}

// Create requests the creation of a new queue.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToQueueCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	queueName := b["queue_name"].(string)
	delete(b, "queue_name")

	_, r.Err = client.Put(createURL(client, queueName), b, r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201, 204},
	})
	return
}
