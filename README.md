# gophercloud

gophercloud is a flexible SDK that allows you to consume and work with OpenStack
clouds in a simple and idiomatic way using golang. Many services are supported,
including Compute, Block Storage, Object Storage, Networking, and Identity.
Each service API is backed with documentation, code samples, unit tests and
acceptance tests.

## How to install

Before installing, you need to ensure that your [GOPATH environment variable](https://golang.org/doc/code.html#GOPATH)
is pointing to an appropriate directory where you want to install gophercloud:

```bash
mkdir $HOME/go
export GOPATH=$HOME/go
```

Once this is set up, you can install the gophercloud package like so:

```bash
go get github.com/rackspace/gophercloud
```

This will install all the source files you need into a `pkg` directory, which is
referenceable from your own source files.

## Getting started

### Credentials

Because you'll be hitting an API, you will need to retrieve your OpenStack
credentials and either store them as environment variables or in your local Go
files. The first method is recommended because it decouples credential
information from source code, allowing you to push the latter to your version
control system without any security risk.

You will need to retrieve the following:

* username
* password
* tenant name or tenant ID
* a valid Keystone identity URL

For users that have the OpenStack dashboard installed, there's a shortcut. If
you visit the `project/access_and_security` path in Horizon and click on the
"Download OpenStack RC File" button at the top right hand corner, you will
download a bash file that exports all of your access details to environment
variables. To execute the file, run `source admin-openrc.sh` and you will be
prompted for your password.

### Authentication

Once you have access to your credentials, you can begin plugging them into
gophercloud. The next step is authentication, and this is handled by a base
"Provider" struct. To get one, you can either pass in your credentials
explicitly, or tell gophercloud to use environment variables:

```go
// Option 1: Pass in the values yourself
opts := gophercloud.AuthOptions{
  IdentityEndpoint: "https://my-openstack.com:5000/v2.0",
  Username: "{username}",
  Password: "{password}",
  TenantID: "{tenant_id}",
}

// Option 2: Use a utility function to retrieve all your environment variables
import "github.com/rackspace/gophercloud/openstack/utils"
opts, err := utils.AuthOptions()
```

Once you have the `opts` variable, you can pass it in and get back a
`ProviderClient` struct:

```go
import "github.com/rackspace/gophercloud/openstack"

provider, err := openstack.AuthenticatedClient(opts)
```

The `ProviderClient` is the top-level client that all of your OpenStack services
derive from. The provider contains all of the authentication details that allow
your Go code to access the API - such as the base URL and token ID.

### Next steps

Cool! You've handled authentication and got your `ProviderClient`. You're now
ready to use an OpenStack service.

* [Getting started with Compute](http://gophercloud.io/docs/compute)
* [Getting started with Object Storage](http://gophercloud.io/docs/object-storage)
* [Getting started with Networking](http://gophercloud.io/docs/networking)
* [Getting started with Block Storage](http://gophercloud.io/docs/block-storage)
* [Getting started with Identity](http://gophercloud.io/docs/identity)

## Contributing

Engaging the community and lowering barriers for contributors is something we
care a lot about. For this reason, we've taken the time to write a [contributing
guide](./CONTRIBUTING.md) for folks interested in getting involved in the project.
If you're not sure what you can do to get involved, feel free to submit an issue
or ping us an e-mail! You don't need to be a Go expert - all members of the
community are welcome.

## Help and feedback

If you're struggling with something or have spotted a potential bug, feel free
to submit an issue to our [bug tracker](/issues) or e-mail us directly at
sdk-support@rackspace.com.
