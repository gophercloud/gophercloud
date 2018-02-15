/*
Package rbacpolicies contains functionality for working with Neutron RBAC Policies.
Role-Based Access Control (RBAC) policy framework enables both operators
and users to grant access to resources for specific projects.

Sharing an object with a specific project is accomplished by creating a
policy entry that permits the target project the access_as_shared action
on that object.

To make a network available as an external network for specific projects
rather than all projects, use the access_as_external action.
If a network is marked as external during creation, it now implicitly creates
a wildcard RBAC policy granting everyone access to preserve previous behavior
before this feature was added.

Example to Create a RBAC Policy

	createOpts := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.ActionAccessShared,
		ObjectType:   "network",
                TargetTenant: "6e547a3bcfe44702889fdeff3c3520c3",
                ObjectID:     "240d22bf-bd17-4238-9758-25f72610ecdc"
	}

	rbacPolicy, err := rbacpolicies.Create(rbacClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package rbacpolicies
