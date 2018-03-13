package buildinfo

import (
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// BuildInfo represents an Build Information (build-info)
type BuildInfo struct {
	API    map[string]interface{} `json:"api"`
	Engine map[string]interface{} `json:"engine"`
}

// ExtractBuildInfo provides access to the list of build info in a page acquired from the ListDetail operation.
func (r commonResult) ExtractBuildInfo() ([]BuildInfo, error) {
	var s struct {
		BuildInfo []BuildInfo `json:"build_info"`
	}
	err := r.ExtractInto(&s)
	return s.BuildInfo, err
}

func (r *BuildInfo) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	_ = json.Unmarshal(b, &data)

	type tmp BuildInfo
	var s struct {
		tmp
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("NO GOOD", err)
		return err
	}
	*r = BuildInfo(s.tmp)

	return nil
}
