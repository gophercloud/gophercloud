/*
Package extendedavailabilityzone provides the ability the ability to extend a
server result with the extended availability zone information.

Example to Get an extended information:

  type serverAvailabilityZoneExt struct {
    servers.Server
    extendedavailabilityzone.AvailabilityZoneExt
  }
  var serverWithAvailabilityZoneExt serverAvailabilityZoneExt

  err := servers.Get(computeClient, "d650a0ce-17c3-497d-961a-43c4af80998a").ExtractInto(&serverWithAvailabilityZoneExt)
  if err != nil {
    panic(err)
  }

  fmt.Printf("%+v\n", serverWithAvailabilityZoneExt)
*/
package extendedavailabilityzone
