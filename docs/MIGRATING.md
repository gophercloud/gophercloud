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

If using [gophercloud/utils](https://github.com/gophercloud/utils), you will
also need to update those imports:

```diff
import (
-	"github.com/gophercloud/gophercloud"
-	serverutils "github.com/gophercloud/utils/openstack/compute/v2/servers"
+	"github.com/gophercloud/gophercloud/v2"
+	serverutils "github.com/gophercloud/utils/v2/openstack/compute/v2/servers"
)
```

### Go version

The minimum go version for Gophercloud v2 is now v1.22.

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

Now that every method accept a context, it is no longer possible to attach
a context to the Provider client. Use per-call context instead.

The `WaitFor` functions now take a context as well, and we've dropped the
timeout argument. This means that the following code:

```go
err = attachments.WaitForStatus(client, attachment.ID, "attached", 60)
```

Must be changed to use a context with timeout. For example:

```go
ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
defer cancel()

err = attachments.WaitForStatus(ctx, client, attachment.ID, "attached")
```

### Error handling

The error types for specific response codes (`ErrDefault400`, `ErrDefault401`, etc.) have been removed.
All unexpected response codes will now return `ErrUnexpectedResponseCode` instead.
For quickly checking whether a request resulted in a specific response code, use the new `ResponseCodeIs` function:

```go
server, err := servers.Get(ctx, client, serverID).Extract()

// before
if _, ok := err.(gophercloud.ErrDefault404); ok {
  handleServerNotFound()
}

// after
if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
  handleServerNotFound()
}
```

Furthermore, the error messages returned by ErrUnexpectedResponseCode now include less newlines than before.
If you match on error messages using regexes, please double-check your regexes.

### Removal of `extensions` modules

A number of services previously supported API extensions but have long since
switched to using microversions to allow API changes. This is now reflected
in Gophercloud v2 and the contents of the follow modules have been largely
migrated:

- `openstack/blockstorage/extensions`
- `openstack/compute/v2/extensions`
- `openstack/identity/v2/extensions`
- `openstack/identity/v3/extensions`

The replacement for these depends on the type of the former extension. For
extensions that added wholly new APIs, these APIs have been moved into the
main module for the corresponding service. These are:

- `openstack/blockstorage/extensions/availabilityzones`

  Moved to `openstack/blockstorage/v2/availabilityzones` and
  `openstack/blockstorage/v3/availabilityzones`.

- `openstack/blockstorage/extensions/backups`

  Moved to `openstack/blockstorage/v2/backups` and
  `openstack/blockstorage/v3/backups`.

- `openstack/blockstorage/extensions/limits`

  Moved to `openstack/blockstorage/v2/limits` and
  `openstack/blockstorage/v3/limits`.

- `openstack/blockstorage/extensions/quotasets`

  Moved to `openstack/blockstorage/v2/quotasets` and
  `openstack/blockstorage/v3/quotasets`.

- `openstack/blockstorage/extensions/schedulerstats`

  Moved to `openstack/blockstorage/v2/schedulerstats` and
  `openstack/blockstorage/v3/schedulerstats`.

- `openstack/blockstorage/extensions/services`

  Moved to `openstack/blockstorage/v2/services` and
  `openstack/blockstorage/v3/services`.

- `openstack/blockstorage/extensions/volumetransfers`

  Moved to `openstack/blockstorage/v2/transfers` and
  `openstack/blockstorage/v3/transfers`.

- `openstack/compute/v2/extensions/aggregates`

  Moved to `openstack/compute/v2/aggregates`.

- `openstack/compute/v2/extensions/attachinterfaces`

  Moved to `openstack/compute/v2/attachinterfaces`.

- `openstack/compute/v2/extensions/diagnostics`

  Moved to `openstack/compute/v2/diagnostics`.

- `openstack/compute/v2/extensions/hypervisors`

  Moved to `openstack/compute/v2/hypervisors`.

- `openstack/compute/v2/extensions/instanceactions`

  Moved to `openstack/compute/v2/instanceactions`.

- `openstack/compute/v2/extensions/keypairs`

  Moved to `openstack/compute/v2/keypairs`.

- `openstack/compute/v2/extensions/quotasets`

  Moved to `openstack/compute/v2/quotasets`.

- `openstack/compute/v2/extensions/remoteconsoles`

  Moved to `openstack/compute/v2/remoteconsoles`.

- `openstack/compute/v2/extensions/secgroups`

  Moved to `openstack/compute/v2/secgroups`.

- `openstack/compute/v2/extensions/servergroups`

  Moved to `openstack/compute/v2/servergroups`.

- `openstack/compute/v2/extensions/services`

  Moved to `openstack/compute/v2/services`.

- `openstack/compute/v2/extensions/tags`

  Moved to `openstack/compute/v2/tags`.

- `openstack/compute/v2/extensions/usage`

  Moved to `openstack/compute/v2/usage`.

- `openstack/compute/v2/extensions/volumeattach`

  Moved to `openstack/compute/v2/volumeattach`.

- `openstack/identity/v2/extensions/admin/roles`

  Moved to `openstack/identity/v2/roles`.

- `openstack/identity/v3/extensions/ec2credentials`

  Moved to `openstack/identity/v3/ec2credentials`.

- `openstack/identity/v3/extensions/ec2tokens`

  Moved to `openstack/identity/v3/ec2tokens`.

- `openstack/identity/v3/extensions/federation`

  Moved to `openstack/identity/v3/federation`

- `openstack/identity/v3/extensions/oauth1`.

  Moved to `openstack/identity/v3/oauth1`

- `openstack/identity/v3/extensions/projectendpoints`

  Moved to `openstack/identity/v3/projectendpoints`.

For extensions that modified existing APIs, these modifications have been
folded into the modified APIs. These are:

- `openstack/blockstorage/extensions/schedulerhints`

  `SchedulerHintOpts` has been renamed to `SchedulerHints` and moved to
  `openstack/blockstorage/v2/volumes` and `openstack/blockstorage/v3/volumes`.
  This is now a required argument of `volumes.Create` for both modules.

- `openstack/blockstorage/extensions/volumeactions`

  All functions and supporting structs and interfaces have been moved to
  `openstack/blockstorage/v2/volumes` and `openstack/blockstorage/v3/volumes`.

- `openstack/blockstorage/extensions/volumehost`

  The `VolumeHostExt` struct has been removed and a `Host` field added to the
  `Volume` struct in `openstack/blockstorage/v2/volumes` and
  `openstack/blockstorage/v3/volumes`.

- `openstack/blockstorage/extensions/volumetenants`

  The `VolumeTenantExt` struct has been removed and a `TenantID` field added to
  the `Volume` struct in `openstack/blockstorage/v2/volumes` and
  `openstack/blockstorage/v3/volumes`.

- `openstack/compute/v2/extensions/bootfromvolume`

  The `CreateOptsExt` struct has been removed and a `BlockDevice` field added
  to the `CreateOpts` struct in `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/diskconfig`

  The `CreateOptsExt` struct has been removed and a `DiskConfig` field added to
  the `CreateOpts` struct in `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/evacuate`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/extendedserverattributes`

  The `ServerAttributesExt` struct has been removed and all fields added to the
  `Server` struct in `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/extendedstatus`

  The `ServerExtendedStatusExt` struct has been removed and all fields added to
  the `Server` struct in `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/injectnetworkinfo`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/lockunlock`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/migrate`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/pauseunpause`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/rescueunrescue`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/resetnetwork`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/resetstate`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/schedulerhints`

  `SchedulerHintOpts` has been moved to `openstack/compute/v2/servers` and
  renamed to `SchedulerHints`. This is now a required argument of
  `servers.Create`.

- `openstack/compute/v2/extensions/serverusage`

  The `serverusage` struct has been removed and all fields added to the
  `Server` struct in `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/shelveunshelve`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/startstop`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

- `openstack/compute/v2/extensions/suspendresume`

  All functions and supporting structs and interfaces have been moved to
  `openstack/compute/v2/servers`.

Finally, for extensions that added new APIs *and* modified existing APIs, the
new APIs are moved into the main module of the corresponding service while the
modifications are folded into the modified APIs. These are:

- `openstack/compute/v2/extensions/availabilityzones`

  The `ServerAvailabilityZoneExt` struct has been removed and a
  `AvailabilityZone` field added to the `Server` struct in
  `openstack/compute/v2/servers`. Everything else is moved moved to
  `openstack/compute/v2/availabilityzones`.

- `openstack/identity/v3/extensions/trusts`

  The `AuthOptsExt` struct has been removed and a `TrustID` field added to the
  `Scope` struct in `openstack/identity/v3/tokens`. Everything else is moved
  moved to `openstack/identity/v3/trusts`.

### Type changes

`loadbalancer/v2/pools/CreateOpts.Members` is now a slice of `CreateMemberOpts`
rather than a slice of `BatchUpdateMemberOpts`.

`blockstorage/v3/volumes/CreateOpts.Multiattach` is removed. Use a volume type
with `multiattach` capability instead.

The following structs are no longer comparable due to the addition of a non-comparable field:

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

### Image

The `imageservice` service is renamed to simply `image` to conform with the other services.

If you previously imported from
`github.com/gophercloud/gophercloud/v2/openstack/imageservice/`, you now need
to import from `github.com/gophercloud/gophercloud/v2/openstack/image/`.

Additionally, `NewImageServiceV2()` is renamed `NewImageV2()`.

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
example, `url.PathEscape` or `url.QueryEscape`), then you should REMOVE that
and pass the intended names to Gophercloud directly.

The `objectstorage/v1/containers.ListOpts#Full` and
`objectstorage/v1/objects.ListOpts#Full` properties are removed from the
Gophercloud API. Plaintext listing is unfixably wrong and won't handle special
characters reliably (i.e. `\n`). Object listing and container listing now
always behave like “Full” did.

Empty container names, container names containing a slash (`/`), and empty
object names are now rejected in Gophercloud before any call to Swift.

The `ErrInvalidContainerName` error has been moved from
`objectstorage/v1/containers` to `objectstorage/v1`. In addition, two new name
validation errors have been added: `objectstorage.v1.ErrEmptyContainerName` and
`objectstorage.v1.ErrEmptyObjectName`.

The `objectstorage/v1/objects.Copy#Destination` field must be in the form
`/container/object`. The function will reject a destination path if it doesn't
start with a slash (`/`).

### Removed services and extensions

Support for services that are no longer supported upstream has been removed.
Users that still rely on theses old services should continue using Gophercloud v1.

- Cinder (Blockstorage) v1
- Neutron (Networking) LBaaS and LBaaS v2 extensions
- Neutron (Networking) FWaaS extension
- Poppy (CDNaaS) service
- Senlin (Clustering) service
