package volumetypes

import "github.com/gophercloud/gophercloud/v2"

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("types")
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id)
}

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("types")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id)
}

func updateURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("types", id)
}

func extraSpecsListURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("types", id, "extra_specs")
}

func extraSpecsGetURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("types", id, "extra_specs", key)
}

func extraSpecsCreateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("types", id, "extra_specs")
}

func extraSpecUpdateURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("types", id, "extra_specs", key)
}

func extraSpecDeleteURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("types", id, "extra_specs", key)
}

func accessURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("types", id, "os-volume-type-access")
}

func accessActionURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("types", id, "action")
}

func createEncryptionURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("types", id, "encryption")
}

func deleteEncryptionURL(client gophercloud.Client, id, encryptionID string) string {
	return client.ServiceURL("types", id, "encryption", encryptionID)
}

func getEncryptionURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("types", id, "encryption")
}

func getEncryptionSpecURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("types", id, "encryption", key)
}

func updateEncryptionURL(client gophercloud.Client, id, encryptionID string) string {
	return client.ServiceURL("types", id, "encryption", encryptionID)
}
