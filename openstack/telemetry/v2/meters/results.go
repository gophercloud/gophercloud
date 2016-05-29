package meters

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type Meter struct {
	MeterId    string `mapstructure:"meter_id"`
	Name       string `json:"name"`
	ProjectId  string `mapstructure:"project_id"`
	ResourceId string `mapstructure:"resource_id"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	Unit       string `json:"user"`
	UserId     string `mapstructure:"user_id"`
}

type ListResult struct {
	gophercloud.Result
}

// Extract interprets any ListResult as an array of Meter
func (r ListResult) Extract() ([]Meter, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response []Meter

	config := &mapstructure.DecoderConfig{
		DecodeHook: toMapFromString,
		Result:     &response,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(r.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type Statistics struct {
	Avg           float32 `json:"avg"`
	Count         int     `json:"count"`
	Duration      float32 `json:"duration"`
	DurationEnd   string  `mapstructure:"duration_end"`
	DurationStart string  `mapstructure:"duration_start"`
	Max           float32 `json:"max"`
	Min           float32 `json:"min"`
	Period        int     `json:"user_id"`
	PeriodEnd     string  `mapstructure:"period_end"`
	PeriodStart   string  `mapstructure:"period_start"`
	Sum           float32 `json:"sum"`
	Unit          string  `json:"unit"`
}

type StatisticsResult struct {
	gophercloud.Result
}

// Extract interprets any serverResult as a Server, if possible.
func (r StatisticsResult) Extract() ([]Statistics, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response []Statistics

	config := &mapstructure.DecoderConfig{
		DecodeHook: toMapFromString,
		Result:     &response,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(r.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func toMapFromString(from reflect.Kind, to reflect.Kind, data interface{}) (interface{}, error) {
	if (from == reflect.String) && (to == reflect.Map) {
		return map[string]interface{}{}, nil
	}
	return data, nil
}
