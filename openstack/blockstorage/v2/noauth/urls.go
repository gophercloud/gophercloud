package noauth

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func getURL(client *gophercloud.ProviderClient, eo EndpointOpts, clientType string) (string, error) {
	url := ""
	if len(eo.CinderEndpoint) > 0 && clientType == "volumev2" {
		url = gophercloud.NormalizeURL(fmt.Sprintf("%s%s", gophercloud.NormalizeURL(eo.CinderEndpoint), client.IdentityBase))
	} else {
		return "", fmt.Errorf("Pass proper EndPointOpt")
	}
	return url, nil
}
