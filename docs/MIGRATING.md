# Migration guide

Gophercloud follows [semver](https://semver.org/) and each major release brings
a number of changes breaking backward compatibility. This guide details those
changes and explains how to migrate from one major version of Gophercloud to
another.

## From v1 to v2

### Change import path

The module is now named `github.com/gophercloud/gophercloud/v2`. Consequently,
you need to update all your imports:

```diff
import (
-	"github.com/gophercloud/gophercloud"
-	"github.com/gophercloud/gophercloud/pagination"
+	"github.com/gophercloud/gophercloud/v2"
+	"github.com/gophercloud/gophercloud/v2/pagination"
)
```

### Go version

The minimum go version for Gophercloud v2 is now v1.21.6.

### Context-awareness

Gophercloud is now context aware, for tracing and cancellation. All function
signatures triggering an HTTP call now take a `context.Context` as their first
argument. 

While you previously called:
```go
myServer, err := servers.Get(client, server.ID)
```

You now need to pass it a context, for example:
```go
ctx := context.TODO()
myServer, err := servers.Get(ctx, client, server.ID)
```

As a consequence of this change, now that every method accept a context, it is no
longer possible to attach a context to the Provider client. Use per-call context instead.

Another consequence is that all the `WaitFor` functions now take a context as
well, and we've dropped the timeout argument. This means that the following
code:

```go
err = attachments.WaitForStatus(client, attachment.ID, "attached", 60)
```

Must be changed to use a context with timeout, such as:
```go
ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
defer cancel()

err = attachments.WaitForStatus(ctx, client, attachment.ID, "attached")
```

### Type changes

`loadbalancer/v2/pools/CreateOpts.Members` is now a slice of `CreateMemberOpts` rather than a slice of `BatchUpdateMemberOpts`.

`blockstorage/v3/volumes/CreateOpts.Multiattach` was removed. Use a volume type with `multiattach` capability instead.

The following structs are no longer comparable due to the addition of a non comparable field:
- `compute/v2/flavors/Flavor`
- `loadbalancer/v2/l7policies/CreateRuleOpts`
- `loadbalancer/v2/l7policies/UpdateOpts`
- `loadbalancer/v2/l7policies/UpdateRuleOpts`
- `loadbalancer/v2/listeners/ListOpts`
- `loadbalancer/v2/monitors/ListOpts`
- `loadbalancer/v2/monitors/CreateOpts`
- `loadbalancer/v2/monitors/UpdateOpts`
- `loadbalancer/v2/pools/ListOpts`

This means that you were previously able to use `==` to compare these objects,
this is no longer the case with Gophercloud v2.

### Image service

The `imageservice` service was renamed to simply `image` to conform with the other services.

If you previously imported from
`github.com/gophercloud/gophercloud/v2/openstack/imageservice/`, you now need
to import from `github.com/gophercloud/gophercloud/v2/openstack/image/`.

Additionally, `NewImageServiceV2()` was renamed `NewImageV2()`.

### Baremetal inventory

The Baremetal inventory types moved from
`baremetalintrospection/v1/introspection` to `baremetal/inventory`. This
includes `BootInfoType`, `CPUType`, `LLDPTLVType`, `InterfaceType`,
`InventoryType`, `MemoryType`, `RootDiskType`, `SystemFirmwareType`,
`SystemVendorType`, `ExtraHardwareDataType`, `ExtraHardwareData`,
`ExtraHardwareDataSection`, `NUMATopology`, `NUMACPU`, `NUMANIC`, and
`NUMARAM`.

Additionally, a few of these types were renamed in the process:
- `ExtraHardwareDataType` became `ExtraDataType`
- `ExtraHardwareData` became `ExtraDataItem`
- `ExtraHardwareDataSection` became `ExtraHardwareDataSection`

### Object storage

Gophercloud now escapes container and object names in all `objects` and
`containers` functions. If you were previously escaping names (with, for
example: `url.PathEscape` or `url.QueryEscape`), then you should REMOVE that
and pass the intended names to Gophercloud directly.

The `containers.ListOpts#Full` and `objects.ListOpts#Full` properties are
REMOVED from the Gophercloud API. The reason for that is: plaintext listing is
unfixably wrong and won't handle special characters reliably (i.e. `\n`).
Object listing and container listing now always behave like “Full” did.

Empty container names, container names containing a slash (`/`), and empty
object names are now rejected in Gophercloud before any call to Swift.

In `objectstorage`: `containers.ErrInvalidContainerName` is now `v1.ErrInvalidContainerName`.
New name validation errors in `objectstorage`:
  * `v1.ErrEmptyContainerName`
  * `v1.ErrEmptyObjectName`

In `objects.Copy`: the `destination` field (e.g.
`objects.CopyOpts#Destination`) must be in the form `/container/object`: the
function will reject a destination path if it doesn't start with a slash (`/`).

### Removed services and extensions

Users that still rely on theses old services should continue using Gophercloud v1.

- Cinder (Blockstorage) v1
- Neutron LBaaS and LBaaS v2 extensions. They have been replaced by Octavia.
- Neutron FWaaS extension.
- Poppy (CDNaaS).
