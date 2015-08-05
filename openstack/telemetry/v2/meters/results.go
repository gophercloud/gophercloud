package meters

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type MeterListResult struct {
	MeterId    string `json:"meter_id"`
	Name       string `json:"name"`
	ProjectId  string `json:"project_id"`
	ResourceId string `json:"resource_id"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	Unit       string `json:"user"`
	UserId     string `json:"user_id"`
}

type listResult struct {
	gophercloud.Result
}

// Extract interprets any listResult as an array of MeterListResult
func (r listResult) Extract() ([]MeterListResult, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response []MeterListResult

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

type MeterStatisticsResult struct {
	Avg           float32 `json:"avg"`
	Count         int     `json:"count"`
	Duration      float32 `json:"duration"`
	DurationEnd   string  `json:"duration_end"`
	DurationStart string  `json:"duration_start"`
	Max           float32 `json:"max"`
	Min           float32 `json:"min"`
	Period        int     `json:"user_id"`
	PeriodEnd     string  `json:"period_end"`
	PeriodStart   string  `json:"period_start"`
	Sum           float32 `json:"sum"`
	Unit          string  `json:"unit"`
}

type statisticsResult struct {
	gophercloud.Result
}

// Extract interprets any serverResult as a Server, if possible.
func (r statisticsResult) Extract() ([]MeterStatisticsResult, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response []MeterStatisticsResult

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
