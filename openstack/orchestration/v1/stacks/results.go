package stacks

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type CreatedStack struct {
	ID    string             `mapstructure:"id"`
	Links []gophercloud.Link `mapstructure:"links"`
}

type CreateResult struct {
	gophercloud.Result
}

func (r CreateResult) Extract() (*CreatedStack, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Stack *CreatedStack `mapstructure:"stack"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	return res.Stack, nil
}

type AdoptResult struct {
	CreateResult
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

type ListedStack struct {
	CreationTime time.Time          `mapstructure:"-"`
	Description  string             `mapstructure:"description"`
	ID           string             `mapstructure:"id"`
	Links        []gophercloud.Link `mapstructure:"links"`
	Name         string             `mapstructure:"stack_name"`
	Status       string             `mapstructure:"stack_status"`
	StatusReason string             `mapstructure:"stack_status_reason"`
	UpdatedTime  time.Time          `mapstructure:"-"`
}

// ExtractStacks extracts and returns a slice of Stacks. It is used while iterating
// over a stacks.List call.
func ExtractStacks(page pagination.Page) ([]ListedStack, error) {
	var res struct {
		Stacks []ListedStack `mapstructure:"stacks"`
	}

	err := mapstructure.Decode(page.(StackPage).Body, &res)
	if err != nil {
		return nil, err
	}

	rawStacks := (((page.(StackPage).Body).(map[string]interface{}))["stacks"]).([]interface{})
	for i := range rawStacks {
		creationTime, err := time.Parse(time.RFC3339, ((rawStacks[i]).(map[string]interface{}))["creation_time"].(string))
		if err != nil {
			return res.Stacks, err
		}
		res.Stacks[i].CreationTime = creationTime

		updatedTime, err := time.Parse(time.RFC3339, ((rawStacks[i]).(map[string]interface{}))["updated_time"].(string))
		if err != nil {
			return res.Stacks, err
		}
		res.Stacks[i].UpdatedTime = updatedTime
	}

	return res.Stacks, nil
}

type RetrievedStack struct {
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

func (r GetResult) Extract() (*RetrievedStack, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Stack *RetrievedStack `mapstructure:"stack"`
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

type PreviewedStack struct {
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

func (r PreviewResult) Extract() (*PreviewedStack, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Stack *PreviewedStack `mapstructure:"stack"`
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

type AbandonedStack struct {
}

type AbandonResult struct {
	gophercloud.Result
}
