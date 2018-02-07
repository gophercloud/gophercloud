package testing

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbac"
)

// CreateRequest is the structure of request body to create rbac-policy.
const CreateRequest = `
{
    "rbac_policy": {
        "action": "access_as_shared",
        "object_type": "network",
        "target_tenant": "6e547a3bcfe44702889fdeff3c3520c3",
        "object_id": "240d22bf-bd17-4238-9758-25f72610ecdc"
    }
}`

// CreateResponse is the structure of response body of rbac-policy create.
const CreateResponse = `
{
    "rbac_policy": {
        "target_tenant": "6e547a3bcfe44702889fdeff3c3520c3",
        "tenant_id": "3de27ce0a2a54cc6ae06dc62dd0ec832",
        "object_type": "network",
        "object_id": "240d22bf-bd17-4238-9758-25f72610ecdc",
        "action": "access_as_shared",
        "project_id": "3de27ce0a2a54cc6ae06dc62dd0ec832",
        "id": "2cf7523a-93b5-4e69-9360-6c6bf986bb7c"
    }
}`

var rbac1 = rbac.Rbac{
	ID:           "2cf7523a-93b5-4e69-9360-6c6bf986bb7c",
	Action:       "access_as_shared",
	ObjectID:     "240d22bf-bd17-4238-9758-25f72610ecdc",
	ObjectType:   "network",
	TenantID:     "3de27ce0a2a54cc6ae06dc62dd0ec832",
	TargetTenant: "6e547a3bcfe44702889fdeff3c3520c3",
	ProjectID:    "3de27ce0a2a54cc6ae06dc62dd0ec832",
}
