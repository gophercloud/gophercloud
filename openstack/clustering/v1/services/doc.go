/*
Package services enables listing services for senlin engine from the OpenStack
Clustering Service.

Example to list services

  allServices, err := services.ListDetail(serviceClient).AllPages()
  if err != nil {
    panic(err)
  }

	allServices, err := services.ExtractServices(service)
	if err != nil {
		panic(err)
	}

  for _, image := range allServices {
		fmt.Printf("%+v\n", service)
	}
*/
package services
