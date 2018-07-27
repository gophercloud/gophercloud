/*
Package networkipavailabilities provides the ability to retrieve and manage
networkipavailabilities through the Neutron API.

Example of Listing NetworkIPAvailabilities

  allPages, err := networkipavailabilities.List(networkClient, networkipavailabilities.ListOpts{}).AllPages()
  if err != nil {
    panic(err)
  }

  allAvailabilities, err := subnetpools.ExtractSubnetPools(allPages)
  if err != nil {
    panic(err)
  }

  for _, availability := range allAvailabilities {
    fmt.Printf("%+v\n", availability)
  }
*/
package networkipavailabilities
