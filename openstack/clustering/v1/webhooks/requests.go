package webhooks

import (
	"github.com/gophercloud/gophercloud"
)

// TriggerOpts Webhooks request parameters
type TriggerOpts struct {
	// V (Optional)	query	string	The webhook implementation version requested.
	V string `q:"V, omitempty"`

	// params (Optional)	query	object	The query string that forms the inputs to use for the targeted action.
	params string `q:"params, omitempty"`
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
