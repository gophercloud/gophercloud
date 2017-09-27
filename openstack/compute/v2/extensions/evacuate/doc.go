/*
Package evacuate provides functionality to evacuates servers that have been
provisioned by the OpenStack Compute service from a failed host to a new host.

Example to Evacuate a Server

	serverID := "b16ba811-199d-4ffd-8839-ba96c1185a67"
	err := evacuate.Evacuate(computeClient, serverID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package evacuate
