package bulk

import (
  "github.com/rackspace/gophercloud"

  "github.com/mitchellh/mapstructure"

  )

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.Result
}

type DeleteBody struct {
  NumberNotFound int `mapstructure:"Number Not Found"`
  ResponseStatus string `mapstructure:"Response Status"`
  Errors []string `mapstructure:"Errors"`
  NumberDeleted int `mapstructure:"Number Deleted"`
  ResponseBody string `mapstructure:"Response Body"`
}

func (dr DeleteResult) ExtractBody() (DeleteBody, error) {
  var resp DeleteBody
  err := mapstructure.Decode(dr.Body, &resp)
  return resp, err
}

type ExtractResult struct {
  gophercloud.Result
}
