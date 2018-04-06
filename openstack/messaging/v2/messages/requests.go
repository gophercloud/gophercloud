package messages

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToMessageCreateMap() (map[string]interface{}, error)
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
	Body map[string]interface{} `json:"body" required:"true"`
}

// ToMessageCreateMap constructs a request body from BatchCreateOpts.
func (opts BatchCreateOpts) ToMessageCreateMap() (map[string]interface{}, error) {
	messages := make([]map[string]interface{}, len(opts))
	for i, message := range opts {
		messageMap, err := message.ToMap()
		if err != nil {
			return nil, err
		}
		messages[i] = messageMap
	}
	return map[string]interface{}{"messages": messages}, nil
}

// ToMap constructs a request body from UpdateOpts.
func (opts CreateOpts) ToMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a message on a specific queue based of off queue name.
func Create(client *gophercloud.ServiceClient, queueName string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToMessageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client, queueName), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}
