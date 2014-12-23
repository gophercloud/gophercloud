package stackresources

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type Resource struct {
	Links        []gophercloud.Link `mapstructure:"links"`
	LogicalID    string             `mapstructure:"logical_resource_id"`
	Name         string             `mapstructure:"resource_name"`
	PhysicalID   string             `mapstructure:"physical_resource_id"`
	RequiredBy   []interface{}      `mapstructure:"required_by"`
	Status       string             `mapstructure:"resource_status"`
	StatusReason string             `mapstructure:"resource_status_reason"`
	Type         string             `mapstructure:"resource_type"`
	UpdatedTime  time.Time          `mapstructure:"-"`
}

type FindResult struct {
	gophercloud.Result
}

func (r FindResult) Extract() ([]Resource, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Res []Resource `mapstructure:"resources"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	resources := r.Body.(map[string]interface{})["resources"].([]map[string]interface{})

	for i, resource := range resources {
		if date, ok := resource["updated_time"]; ok && date != nil {
			t, err := time.Parse(time.RFC3339, date.(string))
			if err != nil {
				return nil, err
			}
			res.Res[i].UpdatedTime = t
		}
	}

	return res.Res, nil
}

// ResourcePage abstracts the raw results of making a List() request against the API.
// As OpenStack extensions may freely alter the response bodies of structures returned to the client, you may only safely access the
// data provided through the ExtractResources call.
type ResourcePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Server results.
func (page ResourcePage) IsEmpty() (bool, error) {
	resources, err := ExtractResources(page)
	if err != nil {
		return true, err
	}
	return len(resources) == 0, nil
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page ResourcePage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"servers_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

// ExtractResources interprets the results of a single page from a List() call, producing a slice of Resource entities.
func ExtractResources(page pagination.Page) ([]Resource, error) {
	casted := page.(ResourcePage).Body

	var response struct {
		Resources []Resource `mapstructure:"resources"`
	}
	err := mapstructure.Decode(casted, &response)
	return response.Resources, err
}

type GetResult struct {
	gophercloud.Result
}

func (r GetResult) Extract() (*Resource, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Res *Resource `mapstructure:"resource"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	resource := r.Body.(map[string]interface{})["resource"].(map[string]interface{})

	if date, ok := resource["updated_time"]; ok && date != nil {
		t, err := time.Parse(time.RFC3339, date.(string))
		if err != nil {
			return nil, err
		}
		res.Res.UpdatedTime = t
	}

	return res.Res, nil
}

type MetadataResult struct {
	gophercloud.Result
}

func (r MetadataResult) Extract() (map[string]string, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Meta map[string]string `mapstructure:"metadata"`
	}

	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	return res.Meta, nil
}
