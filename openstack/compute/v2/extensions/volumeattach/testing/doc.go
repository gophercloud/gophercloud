/*
Package testing holds fixtures (which imports testing),
so that importing volumeattach package does not inadvertently import testing into production code
More information here:
https://github.com/gophercloud/gophercloud/issues/473
*/
package testing
