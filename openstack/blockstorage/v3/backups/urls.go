package backups

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("backups")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id)
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id)
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("backups")
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("backups", "detail")
}

func updateURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id)
}

func restoreURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id, "restore")
}

func exportURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id, "export_record")
}

func importURL(c gophercloud.Client) string {
	return c.ServiceURL("backups", "import_record")
}

func resetStatusURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id, "action")
}

func forceDeleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("backups", id, "action")
}
