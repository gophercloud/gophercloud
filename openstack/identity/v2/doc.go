/*
Package identity provides convenient OpenStack Identity V2 API client access.
This package currently doesn't support the administrative access endpoints, but may appear in the future based on demand.

Authentication

Established convention in the OpenStack community suggests the use of environment variables to hold authentication parameters.
For example, the following settings would be sufficient to authenticate against Rackspace:

	# assumes Bash shell on a POSIX environment; use SET command for Windows.
	export OS_AUTH_URL=https://identity.api.rackspacecloud.com/v2.0
	export OS_USERNAME=xxxx
	export OS_PASSWORD=yyyy

while you'd need these additional settings to authenticate against, e.g., Nebula One:

	export OS_TENANT_ID=zzzz
	export OS_TENANT_NAME=wwww

Be sure to consult with your provider to see which settings you'll need to authenticate with.

A skeletal client gets started with Gophercloud by authenticating against his/her provider, like so:

	package main

	import (
		"fmt"
		"github.com/rackspace/gophercloud/openstack/identity"
		"github.com/rackspace/gophercloud/openstack/utils"
	)

	func main() {
		// Create an initialized set of authentication options based on available OS_*
		// environment variables.
		ao, err := utils.AuthOptions()
		if err != nil {
			panic(err)
		}

		// Attempt to authenticate with them.
		r, err := identity.Authenticate(ao)
		if err != nil {
			panic(err)
		}

		// With each authentication, you receive a master directory of all the services
		// your account can access.  This "service catalog", as OpenStack calls it,
		// provides you the means to exploit other OpenStack services.
		sc, err := identity.GetServiceCatalog(r)
		if err != nil {
			panic(err)
		}

		// Find the desired service(s) for our application.
		computeService, err := findService(sc, "compute", ...)
		if err != nil {
			panic(err)
		}

		blockStorage, err := findService(sc, "block-storage", ...)
		if err != nil {
			panic(err)
		}

		// ... etc ...
	}

NOTE!
Unlike versions 0.1.x of the Gophercloud API,
0.2.0 and later will not provide a service look-up mechanism as a built-in feature of the Identity SDK binding.
The 0.1.x behavior potentially opened its non-US users to legal liability by potentially selecting endpoints in undesirable regions
in a non-obvious manner if a specific region was not explicitly specified.
Starting with 0.2.0 and beyond, you'll need to use either your own service catalog query function or one in a separate package.
This makes it plainly visible to a code auditor that if you indeed desired automatic selection of an arbitrary region,
you made the conscious choice to use that feature.

Extensions

Some OpenStack deployments may support features that other deployments do not.
Anything beyond the scope of standard OpenStack must be scoped by an "extension," a named, yet well-known, change to the API.
Users may invoke IsExtensionAvailable() after grabbing a list of extensions from the server with GetExtensions().
This of course assumes you know the name of the extension ahead of time.

Here's a simple example of listing all the aliases for supported extensions.
Once you have an alias to an extension, everything else about it may be queried through accessors.

	package main

	import (
		"fmt"
		"github.com/rackspace/gophercloud/openstack/identity"
		"github.com/rackspace/gophercloud/openstack/utils"
	)

	func main() {
		// Create an initialized set of authentication options based on available OS_*
		// environment variables.
		ao, err := utils.AuthOptions()
		if err != nil {
			panic(err)
		}

		// Attempt to query extensions.
		exts, err := identity.GetExtensions(ao)
		if err != nil {
			panic(err)
		}

		// Print out a summary of supported extensions
		aliases, err := exts.Aliases()
		if err != nil {
			panic(err)
		}
		fmt.Println("Extension Aliases:")
		for _, alias := range aliases {
			fmt.Printf("  %s\n", alias)
		}
	}
*/
package identity
