package shares

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("shares")
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("shares", "detail")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id)
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id)
}

func updateURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id)
}

func listExportLocationsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "export_locations")
}

func getExportLocationURL(c gophercloud.Client, shareID, id string) string {
	return c.ServiceURL("shares", shareID, "export_locations", id)
}

func grantAccessURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func revokeAccessURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func listAccessRightsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func extendURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func shrinkURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func revertURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func resetStatusURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func forceDeleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func unmanageURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func getMetadataURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "metadata")
}

func getMetadatumURL(c gophercloud.Client, id, key string) string {
	return c.ServiceURL("shares", id, "metadata", key)
}

func setMetadataURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "metadata")
}

func updateMetadataURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("shares", id, "metadata")
}

func deleteMetadatumURL(c gophercloud.Client, id, key string) string {
	return c.ServiceURL("shares", id, "metadata", key)
}
