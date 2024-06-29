package nodes

import "github.com/gophercloud/gophercloud/v2"

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("nodes")
}

func listURL(client gophercloud.Client) string {
	return createURL(client)
}

func listDetailURL(client gophercloud.Client) string {
	return client.ServiceURL("nodes", "detail")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id)
}

func getURL(client gophercloud.Client, id string) string {
	return deleteURL(client, id)
}

func updateURL(client gophercloud.Client, id string) string {
	return deleteURL(client, id)
}

func validateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "validate")
}

func injectNMIURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "management", "inject_nmi")
}

func bootDeviceURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "management", "boot_device")
}

func supportedBootDeviceURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "management", "boot_device", "supported")
}

func statesResourceURL(client gophercloud.Client, id string, state string) string {
	return client.ServiceURL("nodes", id, "states", state)
}

func powerStateURL(client gophercloud.Client, id string) string {
	return statesResourceURL(client, id, "power")
}

func provisionStateURL(client gophercloud.Client, id string) string {
	return statesResourceURL(client, id, "provision")
}

func raidConfigURL(client gophercloud.Client, id string) string {
	return statesResourceURL(client, id, "raid")
}

func biosListSettingsURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "bios")
}

func biosGetSettingURL(client gophercloud.Client, id string, setting string) string {
	return client.ServiceURL("nodes", id, "bios", setting)
}

func vendorPassthruMethodsURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "vendor_passthru", "methods")
}

func vendorPassthruCallURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "vendor_passthru")
}

func maintenanceURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "maintenance")
}

func inventoryURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "inventory")
}

func firmwareListURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "firmware")
}

func virtualMediaURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("nodes", id, "vmedia")
}
