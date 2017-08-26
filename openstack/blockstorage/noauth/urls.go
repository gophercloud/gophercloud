package noauth

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func getURL(client *gophercloud.ProviderClient, endpoint string) string {
	url := gophercloud.NormalizeURL(fmt.Sprintf("%s%s", gophercloud.NormalizeURL(endpoint), client.IdentityBase))
	return url
}
