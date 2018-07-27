/*
Package extendedserverattributes provides the ability to extend a
server result with the extended usage information.
You need to specify at least "2.3" microversion for the ComputeClient to
retrieve most of the OS-EXT-SRV-ATTR fields.

Example to Get an extended information:

  computeClient, err := openstack.NewComputeV2(providerClient, endpointOptions)
  computeClient.Microversion = "2.3"

  type serverAttributesExt struct {
    servers.Server
    extendedserverattributes.ServerAttributesExt
  }
  var serverWithAttributesExt serverAttributesExt

  err := servers.Get(computeClient, "d650a0ce-17c3-497d-961a-43c4af80998a").ExtractInto(&serverWithAttributesExt)
  if err != nil {
    panic(err)
  }

  fmt.Printf("%+v\n", serverWithAttributesExt)
*/
package extendedserverattributes
