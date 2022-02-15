package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/quotas"
)

// List Output is a sample response to a List call.
const QuotaOutput = `
{
  "api_export_size": 1000,
  "recordset_records": 20,
  "zone_records": 500,
  "zone_recordsets": 500,
  "zones": 100
}
`

var (
	Quota = &quotas.Quota{
		APIExporterSize:  1000,
		RecordsetRecords: 20,
		ZoneRecords:      500,
		ZoneRecordsets:   500,
		Zones:            100,
	}
)

// HandleGetSuccessfully configures the test server to respond to a Get request.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quotas/a86dba58-0043-4cc6-a1bb-69d5e86f3ca3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, QuotaOutput)
	})
}
