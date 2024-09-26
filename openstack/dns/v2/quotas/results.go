package quotas

import (
	"github.com/gophercloud/gophercloud/v2"
)

// Extract interprets a GetResult, CreateResult or UpdateResult as a Quota.
// An error is returned if the original call or the extraction failed.
func (r Result) Extract() (*Quota, error) {
	var s *Quota
	err := r.ExtractInto(&s)
	return s, err
}

// ListResult is the result of a Create request. Call its Extract method
// to interpret the result as a Zone.
type Result struct {
	gophercloud.Result
}

// Quota represents a quotas on the system.
type Quota struct {
	APIExporterSize  int `json:"api_export_size"`
	RecordsetRecords int `json:"recordset_records"`
	ZoneRecords      int `json:"zone_records"`
	ZoneRecordsets   int `json:"zone_recordsets"`
	Zones            int `json:"zones"`
}
