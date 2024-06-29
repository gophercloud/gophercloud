package sharetypes

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("types")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id)
}

func listURL(c gophercloud.Client) string {
	return createURL(c)
}

func getDefaultURL(c gophercloud.Client) string {
	return c.ServiceURL("types", "default")
}

func getExtraSpecsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id, "extra_specs")
}

func setExtraSpecsURL(c gophercloud.Client, id string) string {
	return getExtraSpecsURL(c, id)
}

func unsetExtraSpecsURL(c gophercloud.Client, id string, key string) string {
	return c.ServiceURL("types", id, "extra_specs", key)
}

func showAccessURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id, "share_type_access")
}

func addAccessURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id, "action")
}

func removeAccessURL(c gophercloud.Client, id string) string {
	return addAccessURL(c, id)
}
