package roles

import (
	"testing"

	"github.com/rackspace/gophercloud"
)

func TestRoleAssignmentsURL(t *testing.T) {
	client := gophercloud.ServiceClient{Endpoint: "http://localhost:5000/v3/"}
	url := roleAssignmentsURL(&client)
	if url != "http://localhost:5000/v3/role_assignments" {
		t.Errorf("Unexpected list URL generated: [%s]", url)
	}
}
