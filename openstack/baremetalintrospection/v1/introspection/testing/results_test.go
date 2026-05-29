package testing

import (
	"encoding/json"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/v1/introspection"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestHostnameInInventory(t *testing.T) {
	var output introspection.Data
	err := json.Unmarshal([]byte(IntrospectionDataJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal Inventory data: %s", err)
	}

	th.CheckDeepEquals(t, IntrospectionDataRes.Inventory.Hostname, "myawesomehost")
}
