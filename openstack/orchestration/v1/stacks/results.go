package stacks

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateStack struct {
	ID    string             `mapstructure:"id"`
	Links []gophercloud.Link `mapstructure:"links"`
}

type CreateResult struct {
	gophercloud.Result
}

func (r CreateResult) Extract() (*CreateStack, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Stack *CreateStack `json:"stack"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	return res.Stack, nil
}

type AdoptResult struct {
	gophercloud.Result
}

// StackPage is a pagination.Pager that is returned from a call to the List function.
type StackPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Stacks.
func (r StackPage) IsEmpty() (bool, error) {
	stacks, err := ExtractStacks(r)
	if err != nil {
		return true, err
	}
	return len(stacks) == 0, nil
}

type ListStack struct {
	CreationTime time.Time          `mapstructure:"-"`
	Description  string             `mapstructure:"description"`
	ID           string             `mapstructure:"id"`
	Links        []gophercloud.Link `mapstructure:"links"`
	Name         string             `mapstructure:"stack_name"`
	Status       string             `mapstructure:"stack_status"`
	StausReason  string             `mapstructure:"stack_status_reason"`
	UpdatedTime  time.Time          `mapstructure:"-"`
}

// ExtractStacks extracts and returns a slice of Stacks. It is used while iterating
// over a stacks.List call.
func ExtractStacks(page pagination.Page) ([]ListStack, error) {
	var res struct {
		Stacks []ListStack `json:"stacks"`
	}

	err := mapstructure.Decode(page.(StackPage).Body, &res)
	return res.Stacks, err
}

type GetStack struct {
	Capabilities        []interface{}       `mapstructure:"capabilities"`
	CreationTime        time.Time           `mapstructure:"-"`
	Description         string              `mapstructure:"description"`
	DisableRollback     bool                `mapstructure:"disable_rollback"`
	ID                  string              `mapstructure:"id"`
	Links               []gophercloud.Link  `mapstructure:"links"`
	NotificationTopics  []interface{}       `mapstructure:"notification_topics"`
	Outputs             []map[string]string `mapstructure:"outputs"`
	Parameters          map[string]string   `mapstructure:"parameters"`
	Name                string              `mapstructure:"stack_name"`
	Status              string              `mapstructure:"stack_status"`
	StausReason         string              `mapstructure:"stack_status_reason"`
	TemplateDescription string              `mapstructure:"template_description"`
	Timeout             int                 `mapstructure:"timeout_mins"`
	UpdatedTime         time.Time           `mapstructure:"-"`
}

type GetResult struct {
	gophercloud.Result
}

func (r GetResult) Extract() (*GetStack, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Stack *GetStack `json:"stack"`
	}

	config := &mapstructure.DecoderConfig{
		Result:           &res,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(r.Body); err != nil {
		return nil, err
	}

	b := r.Body.(map[string]interface{})["stack"].(map[string]interface{})

	if date, ok := b["creation_time"]; ok && date != nil {
		t, err := time.Parse(time.RFC3339, date.(string))
		if err != nil {
			return nil, err
		}
		res.Stack.CreationTime = t
	}

	if date, ok := b["updated_time"]; ok && date != nil {
		t, err := time.Parse(time.RFC3339, date.(string))
		if err != nil {
			return nil, err
		}
		res.Stack.UpdatedTime = t
	}

	return res.Stack, err
}

type UpdateResult struct {
	gophercloud.ErrResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type PreviewStack struct {
	Capabilities        []interface{}       `mapstructure:"capabilities"`
	CreationTime        time.Time           `mapstructure:"-"`
	Description         string              `mapstructure:"description"`
	DisableRollback     bool                `mapstructure:"disable_rollback"`
	ID                  string              `mapstructure:"id"`
	Links               []gophercloud.Link  `mapstructure:"links"`
	Name                string              `mapstructure:"stack_name"`
	NotificationTopics  []interface{}       `mapstructure:"notification_topics"`
	Parameters          map[string]string   `mapstructure:"parameters"`
	Resources           []map[string]string `mapstructure:"resources"`
	Status              string              `mapstructure:"stack_status"`
	StausReason         string              `mapstructure:"stack_status_reason"`
	TemplateDescription string              `mapstructure:"template_description"`
	Timeout             int                 `mapstructure:"timeout_mins"`
	UpdatedTime         time.Time           `mapstructure:"-"`
}

type PreviewResult struct {
	gophercloud.Result
}

func (r PreviewResult) Extract() (*PreviewStack, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Stack *PreviewStack `json:"stack"`
	}

	config := &mapstructure.DecoderConfig{
		Result:           &res,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(r.Body); err != nil {
		return nil, err
	}

	b := r.Body.(map[string]interface{})["stack"].(map[string]interface{})

	if date, ok := b["creation_time"]; ok && date != nil {
		t, err := time.Parse(time.RFC3339, date.(string))
		if err != nil {
			return nil, err
		}
		res.Stack.CreationTime = t
	}

	if date, ok := b["updated_time"]; ok && date != nil {
		t, err := time.Parse(time.RFC3339, date.(string))
		if err != nil {
			return nil, err
		}
		res.Stack.UpdatedTime = t
	}

	return res.Stack, err
}

type AbandonStack struct {
}

type AbandonResult struct {
	gophercloud.Result
}
