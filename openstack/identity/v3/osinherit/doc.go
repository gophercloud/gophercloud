/*
Package osinherit enables projects to inherit role assignments from
either their owning domain or projects that are higher in the hierarchy.

Example to Assign a Inherited Role to a User to a Domain

	domainID := "a99e9b4e620e4db09a2dfb6e42a01e66"
	userID := "9df1a02f5eb2416a9781e8b0c022d3ae"
	roleID := "9fe2ff9ee4384b1894a90878d3e92bab"

	err := osinherit.Assign(context.TODO(), identityClient, roleID, osinherit.AssignOpts{
		UserID:   userID,
		domainID: domainID,
	}).ExtractErr()

	if err != nil {
		panic(err)
	}

Example to Assign a Inherited Role to a User to a Project's subtree

	projectID := "a99e9b4e620e4db09a2dfb6e42a01e66"
	userID := "9df1a02f5eb2416a9781e8b0c022d3ae"
	roleID := "9fe2ff9ee4384b1894a90878d3e92bab"

	err := osinherit.Assign(context.TODO(), identityClient, roleID, osinherit.AssignOpts{
		UserID:    userID,
		ProjectID: projectID,
	}).ExtractErr()

	if err != nil {
		panic(err)
	}

Example to validate a Inherited Role to a User to a Project's subtree

	projectID := "a99e9b4e620e4db09a2dfb6e42a01e66"
	userID := "9df1a02f5eb2416a9781e8b0c022d3ae"
	roleID := "9fe2ff9ee4384b1894a90878d3e92bab"

	err := osinherit.Validate(context.TODO(), identityClient, roleID, osinherit.validateOpts{
		UserID:    userID,
		ProjectID: projectID,
	}).ExtractErr()

	if err != nil {
		panic(err)
	}

Example to unassign a Inherited Role to a User to a Project's subtree

	projectID := "a99e9b4e620e4db09a2dfb6e42a01e66"
	userID := "9df1a02f5eb2416a9781e8b0c022d3ae"
	roleID := "9fe2ff9ee4384b1894a90878d3e92bab"

	err := osinherit.Unassign(context.TODO(), identityClient, roleID, osinherit.UnassignOpts{
		UserID:    userID,
		ProjectID: projectID,
	}).ExtractErr()

	if err != nil {
		panic(err)
	}
*/
package osinherit
