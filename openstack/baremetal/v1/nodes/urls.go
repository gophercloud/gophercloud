package nodes

import "github.com/gophercloud/gophercloud/v2"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("nodes")
}

func listURL(client *gophercloud.ServiceClient) string {
	return createURL(client)
}

func listDetailURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("nodes", "detail")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return deleteURL(client, id)
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return deleteURL(client, id)
}

func validateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "validate")
}

func injectNMIURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "management", "inject_nmi")
}

func bootDeviceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "management", "boot_device")
}

func supportedBootDeviceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "management", "boot_device", "supported")
}

func statesResourceURL(client *gophercloud.ServiceClient, id string, state string) string {
	return client.ServiceURL("nodes", id, "states", state)
}

func powerStateURL(client *gophercloud.ServiceClient, id string) string {
	return statesResourceURL(client, id, "power")
}

func provisionStateURL(client *gophercloud.ServiceClient, id string) string {
	return statesResourceURL(client, id, "provision")
}

func raidConfigURL(client *gophercloud.ServiceClient, id string) string {
	return statesResourceURL(client, id, "raid")
}

func biosListSettingsURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "bios")
}

func biosGetSettingURL(client *gophercloud.ServiceClient, id string, setting string) string {
	return client.ServiceURL("nodes", id, "bios", setting)
}

func vendorPassthruMethodsURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "vendor_passthru", "methods")
}

func vendorPassthruCallURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "vendor_passthru")
}

func maintenanceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "maintenance")
}

func inventoryURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "inventory")
}

func firmwareListURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "firmware")
}

func virtualMediaURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "vmedia")
}

func virtualInterfaceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("nodes", id, "vifs")
}

func virtualInterfaceDeleteURL(client *gophercloud.ServiceClient, id string, vifID string) string {
	return client.ServiceURL("nodes", id, "vifs", vifID)
}
