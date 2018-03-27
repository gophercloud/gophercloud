package webhooks

import (
	"github.com/gophercloud/gophercloud"
)

// TriggerOpts represents options used for triggering an action
type TriggerOpts struct {
	V      string `q:"V"`
	Params string `q:"params, omitempty"`
}

// TriggerOptsBuilder Query string builder interface for webhooks
type TriggerOptsBuilder interface {
	ToTriggerQuery() (string, error)
}

// Query string builder for webhooks
func (opts TriggerOpts) ToTriggerQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Trigger an action represented by a webhook.
func Trigger(client *gophercloud.ServiceClient, id string, opts TriggerOptsBuilder) (r TriggerResult) {
	url := triggerURL(client, id)
	if opts != nil {
		query, err := opts.ToTriggerQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Post(url, nil, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	return
}
