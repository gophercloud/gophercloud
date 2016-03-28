package webhooks

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns all webhooks for a scaling policy.
func List(client *gophercloud.ServiceClient, groupID, policyID string) pagination.Pager {
	url := listURL(client, groupID, policyID)

	createPageFn := func(r pagination.PageResult) pagination.Page {
		return WebhookPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPageFn)
}

// CreateOptsBuilder is the interface responsible for generating the JSON
// for a Create operation.
type CreateOptsBuilder interface {
	ToWebhookCreateMap() ([]map[string]interface{}, error)
}

// CreateOpts is a slice of CreateOpt structs, that allow the user to create
// multiple webhooks in a single operation.
type CreateOpts []CreateOpt

// CreateOpt represents the options to create a webhook.
type CreateOpt struct {
	// Name [required] is a name for the webhook.
	Name string

	// Metadata [optional] is user-provided key-value metadata.
	// Maximum length for keys and values is 256 characters.
	Metadata map[string]string
}

// ToWebhookCreateMap converts a slice of CreateOpt structs into a map for use
// in the request body of a Create operation.
func (opts CreateOpts) ToWebhookCreateMap() ([]map[string]interface{}, error) {
	var webhooks []map[string]interface{}

	for _, o := range opts {
		if o.Name == "" {
			return nil, errors.New("Cannot create a Webhook without a name.")
		}

		hook := make(map[string]interface{})

		hook["name"] = o.Name

		if o.Metadata != nil {
			hook["metadata"] = o.Metadata
		}

		webhooks = append(webhooks, hook)
	}

	return webhooks, nil
}

// Create requests a new webhook be created and associated with the given group
// and scaling policy.
func Create(client *gophercloud.ServiceClient, groupID, policyID string, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToWebhookCreateMap()

	if err != nil {
		res.Err = err
		return res
	}

	resp, err := client.Post(createURL(client, groupID, policyID), reqBody, &res.Body, nil)

	if err != nil {
		res.Err = err
		return res
	}

	pr := pagination.PageResultFromParsed(resp, res.Body)
	return CreateResult{pagination.SinglePageBase(pr)}
}
