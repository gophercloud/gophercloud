package drivers

import "github.com/gophercloud/gophercloud/v2"

func driversURL(client gophercloud.Client) string {
	return client.ServiceURL("drivers")
}

func driverDetailsURL(client gophercloud.Client, driverName string) string {
	return client.ServiceURL("drivers", driverName)
}

func driverPropertiesURL(client gophercloud.Client, driverName string) string {
	return client.ServiceURL("drivers", driverName, "properties")
}

func driverDiskPropertiesURL(client gophercloud.Client, driverName string) string {
	return client.ServiceURL("drivers", driverName, "raid", "logical_disk_properties")
}
