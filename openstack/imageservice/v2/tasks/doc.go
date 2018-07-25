/*
Package tasks enables management and retrieval of tasks from the OpenStack
Imageservice.

Example to List Tasks

  listOpts := tasks.ListOpts{
    Owner: "424e7cf0243c468ca61732ba45973b3e",
  }

  allPages, err := tasks.List(imagesClient, listOpts).AllPages()
  if err != nil {
    panic(err)
  }

  allTasks, err := tasks.ExtractTasks(allPages)
  if err != nil {
    panic(err)
  }

  for _, task := range allTasks {
    fmt.Printf("%+v\n", task)
  }
*/
package tasks
