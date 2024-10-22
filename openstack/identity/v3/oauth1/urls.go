package oauth1

import "github.com/gophercloud/gophercloud/v2"

func consumersURL(c gophercloud.Client) string {
	return c.ServiceURL("OS-OAUTH1", "consumers")
}

func consumerURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("OS-OAUTH1", "consumers", id)
}

func requestTokenURL(c gophercloud.Client) string {
	return c.ServiceURL("OS-OAUTH1", "request_token")
}

func authorizeTokenURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("OS-OAUTH1", "authorize", id)
}

func createAccessTokenURL(c gophercloud.Client) string {
	return c.ServiceURL("OS-OAUTH1", "access_token")
}

func userAccessTokensURL(c gophercloud.Client, userID string) string {
	return c.ServiceURL("users", userID, "OS-OAUTH1", "access_tokens")
}

func userAccessTokenURL(c gophercloud.Client, userID string, id string) string {
	return c.ServiceURL("users", userID, "OS-OAUTH1", "access_tokens", id)
}

func userAccessTokenRolesURL(c gophercloud.Client, userID string, id string) string {
	return c.ServiceURL("users", userID, "OS-OAUTH1", "access_tokens", id, "roles")
}

func userAccessTokenRoleURL(c gophercloud.Client, userID string, id string, roleID string) string {
	return c.ServiceURL("users", userID, "OS-OAUTH1", "access_tokens", id, "roles", roleID)
}

func authURL(c gophercloud.Client) string {
	return c.ServiceURL("auth", "tokens")
}
