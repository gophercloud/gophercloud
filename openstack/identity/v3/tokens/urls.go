package tokens

import "github.com/gophercloud/gophercloud/v2"

func tokenURL(c gophercloud.Client) string {
	return c.ServiceURL("auth", "tokens")
}
