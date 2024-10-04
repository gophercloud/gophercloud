package members

import "github.com/gophercloud/gophercloud/v2"

func imageMembersURL(c gophercloud.Client, imageID string) string {
	return c.ServiceURL("images", imageID, "members")
}

func listMembersURL(c gophercloud.Client, imageID string) string {
	return imageMembersURL(c, imageID)
}

func createMemberURL(c gophercloud.Client, imageID string) string {
	return imageMembersURL(c, imageID)
}

func imageMemberURL(c gophercloud.Client, imageID string, memberID string) string {
	return c.ServiceURL("images", imageID, "members", memberID)
}

func getMemberURL(c gophercloud.Client, imageID string, memberID string) string {
	return imageMemberURL(c, imageID, memberID)
}

func updateMemberURL(c gophercloud.Client, imageID string, memberID string) string {
	return imageMemberURL(c, imageID, memberID)
}

func deleteMemberURL(c gophercloud.Client, imageID string, memberID string) string {
	return imageMemberURL(c, imageID, memberID)
}
