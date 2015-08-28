package members

import "github.com/rackspace/gophercloud"

func imageMembersURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL("images", imageID, "members")
}

var listMembersURL = imageMembersURL

var createMemberURL = imageMembersURL

func imageMemberURL(c *gophercloud.ServiceClient, imageID string, memberID string) string {
	return c.ServiceURL("images", imageID, "members", memberID)
}

var getMemberURL = imageMemberURL

var updateMemberURL = imageMemberURL

var deleteMemberURL = imageMemberURL
