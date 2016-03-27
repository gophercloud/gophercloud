package webhooks

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type webhookResult struct {
	gophercloud.Result
}

// Webhook represents a webhook associted with a scaling policy.
type Webhook struct {
	// UUID for the webhook.
	ID string `mapstructure:"id" json:"id"`

	// Name of the webhook.
	Name string `mapstructure:"name" json:"name"`

	// Links associated with the webhook, including the capability URL.
	Links []gophercloud.Link `mapstructure:"links" json:"links"`

	// Metadata associated with the webhook.
	Metadata map[string]string `mapstructure:"metadata" json:"metadata"`
}

// WebhookPage is the page returned by a pager when traversing over a collection
// of webhooks.
type WebhookPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a page contains no Webhook results.
func (page WebhookPage) IsEmpty() (bool, error) {
	hooks, err := ExtractWebhooks(page)

	if err != nil {
		return true, err
	}

	return len(hooks) == 0, nil
}

// ExtractWebhooks interprets the results of a single page from a List() call,
// producing a slice of Webhooks.
func ExtractWebhooks(page pagination.Page) ([]Webhook, error) {
	casted := page.(WebhookPage).Body

	var response struct {
		Webhooks []Webhook `mapstructure:"webhooks"`
	}

	err := mapstructure.Decode(casted, &response)

	if err != nil {
		return nil, err
	}

	return response.Webhooks, err
}
