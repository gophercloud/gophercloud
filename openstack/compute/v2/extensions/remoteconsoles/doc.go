/*
Package remoteconsoles provides the ability to create server remote consoles
through the Compute API.

Example of Creating a new RemoteConsole

  createOpts := remoteconsoles.CreateOpts{
    Protocol: string(remoteconsoles.RemoteConsoleVNCProtocol),
    Type:     string(remoteconsoles.RemoteConsoleNoVNCType),
	}
	serverID := "b16ba811-199d-4ffd-8839-ba96c1185a67"
  remtoteConsole, err := remoteconsoles.Create(computeClient, serverID, createOpts).Extract()
  if err != nil {
    panic(err)
  }

  fmt.Printf("%+v\n", task)
*/
package remoteconsoles
