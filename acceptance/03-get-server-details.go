package main

import (
	"fmt"
	"flag"
	"github.com/rackspace/gophercloud"
)

var id = flag.String("i", "", "Server ID to get info on.  Defaults to first server in your account if unspecified.")
var rgn = flag.String("r", "DFW", "Datacenter region")
var quiet = flag.Bool("quiet", false, "Run quietly, for acceptance testing.  $? non-zero if issue.")

func main() {
	flag.Parse()

	withIdentity(func(auth gophercloud.AccessProvider) {
		withServerApi(auth, func(servers gophercloud.CloudServersProvider) {
			serverId := ""
			if *id == "" {
				ss, err := servers.ListServers()
				if err != nil {
					panic(err)
				}
				// We could just cheat and dump the server details from ss[0].
				// But, that tests ListServers(), and not ServerById().  So, we
				// elect not to cheat.
				serverId = ss[0].Id
			} else {
				serverId = *id
			}

			s, err := servers.ServerById(serverId)
			if err != nil {
				panic(err)
			}

			configs := []string {
				"Access IPv4: %s\n",
				"Access IPv6: %s\n",
				"    Created: %s\n",
				"     Flavor: %s\n",
				"    Host ID: %s\n",
				"         ID: %s\n",
				"      Image: %s\n",
				"       Name: %s\n",
				"   Progress: %s\n",
				"     Status: %s\n",
				"  Tenant ID: %s\n",
				"    Updated: %s\n",
				"    User ID: %s\n",
			}

			values := []string {
				s.AccessIPv4,
				s.AccessIPv6,
				s.Created,
				s.Flavor.Id,
				s.HostId,
				s.Id,
				s.Image.Id,
				s.Name,
				fmt.Sprintf("%d", s.Progress),
				s.Status,
				s.TenantId,
				s.Updated,
				s.UserId,
			}

			if !*quiet {
				fmt.Println("Server info:")
				for i, _ := range configs {
					fmt.Printf(configs[i], values[i])
				}
			}
		})
	})
}
