package configurations

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client) string {
	return c.ServiceURL("configurations")
}

func resourceURL(c gophercloud.Client, configID string) string {
	return c.ServiceURL("configurations", configID)
}

func instancesURL(c gophercloud.Client, configID string) string {
	return c.ServiceURL("configurations", configID, "instances")
}

func listDSParamsURL(c gophercloud.Client, datastoreID, versionID string) string {
	return c.ServiceURL("datastores", datastoreID, "versions", versionID, "parameters")
}

func getDSParamURL(c gophercloud.Client, datastoreID, versionID, paramID string) string {
	return c.ServiceURL("datastores", datastoreID, "versions", versionID, "parameters", paramID)
}

func listGlobalParamsURL(c gophercloud.Client, versionID string) string {
	return c.ServiceURL("datastores", "versions", versionID, "parameters")
}

func getGlobalParamURL(c gophercloud.Client, versionID, paramID string) string {
	return c.ServiceURL("datastores", "versions", versionID, "parameters", paramID)
}
