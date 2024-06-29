package datastores

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client) string {
	return c.ServiceURL("datastores")
}

func resourceURL(c gophercloud.Client, dsID string) string {
	return c.ServiceURL("datastores", dsID)
}

func versionsURL(c gophercloud.Client, dsID string) string {
	return c.ServiceURL("datastores", dsID, "versions")
}

func versionURL(c gophercloud.Client, dsID, versionID string) string {
	return c.ServiceURL("datastores", dsID, "versions", versionID)
}
