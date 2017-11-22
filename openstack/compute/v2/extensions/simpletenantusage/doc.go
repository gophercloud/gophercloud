/*
Package simpletenantusage provides information and interaction with the
SimpleTenantUsage extension for the OpenStack Compute service.

Example to get usage for an entire region
	page, err := simpletenantusage.Get(computeClient).AllPages()
	if err != nil {
		return err
	}
	usage, err := simpletenantusage.ExtractSimpleTenantUsage(page)
	if err != nil {
		return err
	}

Example to get usage for a particular tenant.
	page, err := simpletenantusage.GetTenant(computeClient, tenantID).AllPages()
	if err != nil {
		return err
	}
	usage, err := simpletenantusage.ExtractSimpleTenantUsages(page)
	if err != nil {
		return err
	}

*/
package simpletenantusage
