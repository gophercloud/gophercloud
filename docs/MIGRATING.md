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

#### With gophercloud/utils

If using the `utils` library, note that the `IDFromName` functions return
`ErrResourceNotFound` rather than `ErrUnexpectedResponseCode`. In that
scenario, type assertions for a "not found" error are still necessary:

```Go
func IsNotFound(err error) bool {
	if _, ok := err.(gophercloud.ErrResourceNotFound); ok { // <-- this
		return true
	}

	return gophercloud.ResponseCodeIs(err, http.StatusNotFound)
}
```

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

- `openstack/compute/v2/extensions/limits`

  Moved to `openstack/compute/v2/limits`.

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

  `SchedulerHints` has been renamed to `SchedulerHintOpts` and moved to
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

  `SchedulerHints` has been moved to `openstack/compute/v2/servers` and
  renamed to `SchedulerHintOpts`. This is now a required argument of
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

For extensions that added new APIs *and* modified existing APIs, the new APIs
are moved into the main module of the corresponding service while the
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

Finally, for extensions that are deprecated and have been removed in a
microversion, the APIs were removed entirely. These are:

- `openstack/compute/v2/extensions/defsecrules`

  This was a proxy for the Networking service, Neutron. Use
  `openstack/networking/v2/extensions/security/groups` instead.

- `openstack/compute/v2/extensions/floatingips`

  This was a proxy for the Networking service, Neutron. Use
  `openstack/networking/v2/extensions/layer3/floatingips` instead.

- `openstack/compute/v2/extensions/images`

  This was a proxy for the Image service, Glance. Use
  `openstack/image/v2/images` instead.

- `openstack/compute/v2/extensions/networks`

  This was a proxy for the Networking service, Neutron. Use
  `openstack/networking/v2/networks` instead.

- `openstack/compute/v2/extensions/tenantnetworks`

  This was a proxy for the Networking service, Neutron. Use
  `openstack/networking/v2/networks` instead.

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

- Cinder (Blockstorage) v1 (`openstack/blockstorage/v1`)
- Neutron (Networking) LBaaS and LBaaS v2 extensions
  (`openstack/networking/v2/extensions/lbaas`,
  `openstack/networking/v2/extensions/lbaas_v2`)
- Neutron (Networking) FWaaS extension
  (`openstack/networking/v2/extensions/fwaas`)
- Poppy (CDNaaS) service (`openstack/cdn`)
- Senlin (Clustering) service (`openstack/clustering`)

### Script-assisted migration

#### Expected outcome

After running the script, your code may not compile. The idea is that at this point, you're only left with a few changes that can't reasonably be automated.

#### What it does

* Add `/v2` to all Gophercloud imports, except to the packages that have been removed without replacement
* Adjust the import path of moved packages
* Adjust the package identifier in the code where possible
* Add `context.TODO()` where required

#### Limitations

* it doesn't fix the use of removed extensions. For example, if you used `openstack/blockstorage/extensions/availabilityzones`, you will have to manually put that back into e.g. `servers.CreateOpts`
* it will just put `context.TODO()` where a context is required to satisfy the function signature. It's up to you to actually replace that with a variable and provide proper cancellation
* it will add `context.TODO()` to `blockstorage/v1` calls, even though that package only exists in Gophercloud v1

```bash
# Adjust the blockstorage version appropriately
blockstorageversion=v3

openstack='github.com/gophercloud/gophercloud/openstack'
openstack_utils='github.com/gophercloud/utils/openstack'
find . -type f -name '*.go' -not -path "*/vendor/*" -exec sed -i '
    /^import ($/,/^)$/ {

        # 1: These packages have been removed and their functionality moved into the main module for the corresponding service.
        /\(\/openstack\/blockstorage\/v1\|\/openstack\/networking\/v2\/extensions\/lbaas\|\/openstack\/networking\/v2\/extensions\/lbaas_v2\|\/openstack\/networking\/v2\/extensions\/fwaas\|\/openstack\/cdn\|\/openstack\/clustering\)/! {
            /\/openstack\/blockstorage\/extensions\/volumehost/d
            /\/openstack\/blockstorage\/extensions\/volumetenants/d
            /\/openstack\/compute\/v2\/extensions\/bootfromvolume/d
            /\/openstack\/compute\/v2\/extensions\/diskconfig/d
            /\/openstack\/compute\/v2\/extensions\/extendedserverattributes/d
            /\/openstack\/compute\/v2\/extensions\/extendedstatus/d
            /\/openstack\/compute\/v2\/extensions\/schedulerhints/d
            /\/openstack\/compute\/v2\/extensions\/serverusage/d
            /\/openstack\/compute\/v2\/extensions\/availabilityzones/d
            /\/openstack\/identity\/v3\/extensions\/trusts/d
        }

        '"
        # 2: Functions and supporting structs and interfaces of these packages have been moved to an existing package
        s|${openstack}/blockstorage/extensions/schedulerhints|${openstack}/blockstorage/${blockstorageversion}/volumes|g
        s|${openstack}/blockstorage/extensions/volumeactions|${openstack}/blockstorage/${blockstorageversion}/volumes|g
        s|${openstack}/compute/v2/extensions/evacuate|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/injectnetworkinfo|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/lockunlock|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/migrate|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/pauseunpause|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/rescueunrescue|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/resetnetwork|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/resetstate|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/shelveunshelve|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/startstop|${openstack}/compute/v2/servers|g
        s|${openstack}/compute/v2/extensions/suspendresume|${openstack}/compute/v2/servers|g

        # 3: These packages have been renamed
        s|${openstack}/imageservice|${openstack}/image|g
        s|${openstack_utils}/imageservice|${openstack_utils}/image|g
        s|${openstack}/blockstorage/extensions/availabilityzones|${openstack}/blockstorage/${blockstorageversion}/availabilityzones|g
        s|${openstack}/blockstorage/extensions/backups|${openstack}/blockstorage/${blockstorageversion}/backups|g
        s|${openstack}/blockstorage/extensions/limits|${openstack}/blockstorage/${blockstorageversion}/limits|g
        s|${openstack}/blockstorage/extensions/quotasets|${openstack}/blockstorage/${blockstorageversion}/quotasets|g
        s|${openstack}/blockstorage/extensions/schedulerstats|${openstack}/blockstorage/${blockstorageversion}/schedulerstats|g
        s|${openstack}/blockstorage/extensions/services|${openstack}/blockstorage/${blockstorageversion}/services|g
        s|${openstack}/blockstorage/extensions/volumetransfers|${openstack}/blockstorage/${blockstorageversion}/transfers|g
        s|${openstack}/compute/v2/extensions/aggregates|${openstack}/compute/v2/aggregates|g
        s|${openstack}/compute/v2/extensions/attachinterfaces|${openstack}/compute/v2/attachinterfaces|g
        s|${openstack}/compute/v2/extensions/diagnostics|${openstack}/compute/v2/diagnostics|g
        s|${openstack}/compute/v2/extensions/hypervisors|${openstack}/compute/v2/hypervisors|g
        s|${openstack}/compute/v2/extensions/instanceactions|${openstack}/compute/v2/instanceactions|g
        s|${openstack}/compute/v2/extensions/keypairs|${openstack}/compute/v2/keypairs|g
        s|${openstack}/compute/v2/extensions/limits|${openstack}/compute/v2/limits|g
        s|${openstack}/compute/v2/extensions/quotasets|${openstack}/compute/v2/quotasets|g
        s|${openstack}/compute/v2/extensions/remoteconsoles|${openstack}/compute/v2/remoteconsoles|g
        s|${openstack}/compute/v2/extensions/secgroups|${openstack}/compute/v2/secgroups|g
        s|${openstack}/compute/v2/extensions/servergroups|${openstack}/compute/v2/servergroups|g
        s|${openstack}/compute/v2/extensions/services|${openstack}/compute/v2/services|g
        s|${openstack}/compute/v2/extensions/tags|${openstack}/compute/v2/tags|g
        s|${openstack}/compute/v2/extensions/usage|${openstack}/compute/v2/usage|g
        s|${openstack}/compute/v2/extensions/volumeattach|${openstack}/compute/v2/volumeattach|g
        s|${openstack}/identity/v2/extensions/admin/roles|${openstack}/identity/v2/roles|g
        s|${openstack}/identity/v3/extensions/ec2credentials|${openstack}/identity/v3/ec2credentials|g
        s|${openstack}/identity/v3/extensions/ec2tokens|${openstack}/identity/v3/ec2tokens|g
        s|${openstack}/identity/v3/extensions/federation|${openstack}/identity/v3/federation|g
        s|${openstack}/identity/v3/extensions/oauth1|${openstack}/identity/v3/oauth1|g
        s|${openstack}/identity/v3/extensions/projectendpoints|${openstack}/identity/v3/projectendpoints|g

        # 4: These removed packages existed as proxies of others
        s|${openstack}/compute/v2/extensions/defsecrules|${openstack}/networking/v2/extensions/security/groups|g
        s|${openstack}/compute/v2/extensions/floatingips|${openstack}/networking/v2/extensions/layer3/floatingips|g
        s|${openstack}/compute/v2/extensions/images|${openstack}/image/v2/images|g
        s|${openstack}/compute/v2/extensions/networks|${openstack}/networking/v2/networks|g
        s|${openstack}/compute/v2/extensions/tenantnetworks|${openstack}/networking/v2/networks|g
        "'

        # 5: Update to v2, except for packages that were removed without replacement
        s|github.com/gophercloud/utils|github.com/gophercloud/utils/v2|g
        /\(\/openstack\/blockstorage\/v1\|\/openstack\/networking\/v2\/extensions\/lbaas\|\/openstack\/networking\/v2\/extensions\/lbaas_v2\|\/openstack\/networking\/v2\/extensions\/fwaas\|\/openstack\/cdn\|\/openstack\/clustering\)/! s|github.com/gophercloud/gophercloud|github.com/gophercloud/gophercloud/v2|g
    }

    /^)$/,$ {

        # 6: Rename identifiers of items of step 2 above
        s#\(schedulerhints\|volumeactions\)\.\([A-Z][A-Z_a-z_0-9]*\)#volumes.\2#g
        s#\(evacuate\|injectnetworkinfo\|lockunlock\|migrate\|pauseunpause\|rescueunrescue\|resetnetwork\|resetstate\|shelveunshelve\|startstop\|suspendresume\)\.\([A-Z][A-Z_a-z_0-9]*\)#servers.\2#g

        # 7: Add context.TODO()
        s#\(accept\.Create\|accept\.Get\|accounts\.Get\|accounts\.Update\|acls\.DeleteContainerACL\|acls\.DeleteSecretACL\|acls\.GetContainerACL\|acls\.GetSecretACL\|acls\.SetContainerACL\|acls\.SetSecretACL\|acls\.UpdateContainerACL\|acls\.UpdateSecretACL\|addressscopes\.Create\|addressscopes\.Delete\|addressscopes\.Get\|addressscopes\.Update\|agents\.Delete\|agents\.Get\|agents\.ListDHCPNetworks\|agents\.ListL3Routers\|agents\.RemoveBGPSpeaker\|agents\.RemoveDHCPNetwork\|agents\.RemoveL3Router\|agents\.ScheduleBGPSpeaker\|agents\.ScheduleDHCPNetwork\|agents\.ScheduleL3Router\|agents\.Update\|aggregates\.AddHost\|aggregates\.Create\|aggregates\.Delete\|aggregates\.Get\|aggregates\.RemoveHost\|aggregates\.SetMetadata\|aggregates\.Update\|allocations\.Create\|allocations\.Delete\|allocations\.Get\|amphorae\.Failover\|amphorae\.Get\|apiversions\.Get\|apiversions\.List\|applicationcredentials\.Create\|applicationcredentials\.Delete\|applicationcredentials\.DeleteAccessRule\|applicationcredentials\.Get\|applicationcredentials\.GetAccessRule\|attachinterfaces\.Create\|attachinterfaces\.Delete\|attachinterfaces\.Get\|attachments\.Complete\|attachments\.Create\|attachments\.Delete\|attachments\.Get\|attachments\.Update\|attachments\.WaitForStatus\|backups\.Create\|backups\.Delete\|backups\.Export\|backups\.ForceDelete\|backups\.Get\|backups\.Import\|backups\.ResetStatus\|backups\.RestoreFromBackup\|backups\.Update\|bgpvpns\.Create\|bgpvpns\.CreateNetworkAssociation\|bgpvpns\.CreatePortAssociation\|bgpvpns\.CreateRouterAssociation\|bgpvpns\.Delete\|bgpvpns\.DeleteNetworkAssociation\|bgpvpns\.DeletePortAssociation\|bgpvpns\.DeleteRouterAssociation\|bgpvpns\.Get\|bgpvpns\.GetNetworkAssociation\|bgpvpns\.GetPortAssociation\|bgpvpns\.GetRouterAssociation\|bgpvpns\.Update\|bgpvpns\.UpdatePortAssociation\|bgpvpns\.UpdateRouterAssociation\|buildinfo\.Get\|capsules\.Create\|capsules\.Delete\|capsules\.Get\|certificates\.Create\|certificates\.Get\|certificates\.Update\|claims\.Create\|claims\.Delete\|claims\.Get\|claims\.Update\|clusters\.Create\|clusters\.Delete\|clusters\.Get\|clusters\.Resize\|clusters\.Update\|clusters\.Upgrade\|clustertemplates\.Create\|clustertemplates\.Delete\|clustertemplates\.Get\|clustertemplates\.Update\|conductors\.Get\|config\.NewProviderClient\|configurations\.Create\|configurations\.Delete\|configurations\.Get\|configurations\.GetDatastoreParam\|configurations\.GetGlobalParam\|configurations\.Replace\|configurations\.Update\|containers\.BulkDelete\|containers\.Create\|containers\.CreateConsumer\|containers\.CreateSecretRef\|containers\.Delete\|containers\.DeleteConsumer\|containers\.DeleteSecretRef\|containers\.Get\|containers\.Update\|credentials\.Create\|credentials\.Delete\|credentials\.Get\|credentials\.Update\|crontriggers\.Create\|crontriggers\.Delete\|crontriggers\.Get\|databases\.Create\|databases\.Delete\|datastores\.Get\|datastores\.GetVersion\|diagnostics\.Get\|domains\.Create\|domains\.Delete\|domains\.Get\|domains\.Update\|drivers\.GetDriverDetails\|drivers\.GetDriverDiskProperties\|drivers\.GetDriverProperties\|ec2credentials\.Create\|ec2credentials\.Delete\|ec2credentials\.Get\|ec2tokens\.Create\|ec2tokens\.ValidateS3Token\|endpointgroups\.Create\|endpointgroups\.Delete\|endpointgroups\.Get\|endpointgroups\.Update\|endpoints\.Create\|endpoints\.Delete\|endpoints\.Update\|executions\.Create\|executions\.Delete\|executions\.Get\|extensions\.Get\|extraroutes\.Add\|extraroutes\.Remove\|federation\.CreateMapping\|federation\.DeleteMapping\|federation\.GetMapping\|federation\.UpdateMapping\|flavorprofiles\.Create\|flavorprofiles\.Delete\|flavorprofiles\.Get\|flavorprofiles\.Update\|flavors\.AddAccess\|flavors\.Create\|flavors\.CreateExtraSpecs\|flavors\.Delete\|flavors\.DeleteExtraSpec\|flavors\.Get\|flavors\.GetExtraSpec\|flavors\.ListExtraSpecs\|flavors\.RemoveAccess\|flavors\.Update\|flavors\.UpdateExtraSpec\|floatingips\.Create\|floatingips\.Delete\|floatingips\.Get\|floatingips\.Update\|gophercloud\.WaitFor\|groups\.Create\|groups\.Delete\|groups\.Get\|groups\.RemoveEgressPolicy\|groups\.RemoveIngressPolicy\|groups\.Update\|hypervisors\.Get\|hypervisors\.GetStatistics\|hypervisors\.GetUptime\|ikepolicies\.Create\|ikepolicies\.Delete\|ikepolicies\.Get\|ikepolicies\.Update\|imagedata\.Download\|imagedata\.Stage\|imagedata\.Upload\|imageimport\.Create\|imageimport\.Get\|images\.Create\|images\.Delete\|images\.Get\|images\.Update\|instanceactions\.Get\|instances\.AttachConfigurationGroup\|instances\.Create\|instances\.Delete\|instances\.DetachConfigurationGroup\|instances\.EnableRootUser\|instances\.Get\|instances\.IsRootEnabled\|instances\.Resize\|instances\.ResizeVolume\|instances\.Restart\|introspection\.AbortIntrospection\|introspection\.GetIntrospectionData\|introspection\.GetIntrospectionStatus\|introspection\.ReApplyIntrospection\|introspection\.StartIntrospection\|ipsecpolicies\.Create\|ipsecpolicies\.Delete\|ipsecpolicies\.Get\|ipsecpolicies\.Update\|keypairs\.Create\|keypairs\.Delete\|keypairs\.Get\|l7policies\.Create\|l7policies\.CreateRule\|l7policies\.Delete\|l7policies\.DeleteRule\|l7policies\.Get\|l7policies\.GetRule\|l7policies\.Update\|l7policies\.UpdateRule\|limits\.BatchCreate\|limits\.Delete\|limits\.Get\|limits\.GetEnforcementModel\|limits\.Update\|listeners\.Create\|listeners\.Delete\|listeners\.Get\|listeners\.GetStats\|listeners\.Update\|loadbalancers\.Create\|loadbalancers\.Delete\|loadbalancers\.Failover\|loadbalancers\.Get\|loadbalancers\.GetStats\|loadbalancers\.GetStatuses\|loadbalancers\.Update\|members\.Create\|members\.Delete\|members\.Get\|members\.Update\|messages\.Create\|messages\.Delete\|messages\.DeleteMessages\|messages\.Get\|messages\.GetMessages\|messages\.PopMessages\|monitors\.Create\|monitors\.Delete\|monitors\.Get\|monitors\.Update\|networkipavailabilities\.Get\|networks\.Create\|networks\.Delete\|networks\.Get\|networks\.Update\|nodegroups\.Create\|nodegroups\.Delete\|nodegroups\.Get\|nodegroups\.Update\|nodes\.AttachVirtualMedia\|nodes\.ChangePowerState\|nodes\.ChangeProvisionState\|nodes\.Create\|nodes\.CreateSubscription\|nodes\.Delete\|nodes\.DeleteSubscription\|nodes\.DetachVirtualMedia\|nodes\.Get\|nodes\.GetAllSubscriptions\|nodes\.GetBIOSSetting\|nodes\.GetBootDevice\|nodes\.GetInventory\|nodes\.GetSubscription\|nodes\.GetSupportedBootDevices\|nodes\.GetVendorPassthruMethods\|nodes\.InjectNMI\|nodes\.ListBIOSSettings\|nodes\.ListFirmware\|nodes\.SetBootDevice\|nodes\.SetMaintenance\|nodes\.SetRAIDConfig\|nodes\.UnsetMaintenance\|nodes\.Update\|nodes\.Validate\|nodes\.WaitForProvisionState\|oauth1\.AuthorizeToken\|oauth1\.Create\|oauth1\.CreateAccessToken\|oauth1\.CreateConsumer\|oauth1\.DeleteConsumer\|oauth1\.GetAccessToken\|oauth1\.GetAccessTokenRole\|oauth1\.GetConsumer\|oauth1\.RequestToken\|oauth1\.RevokeAccessToken\|oauth1\.UpdateConsumer\|objects\.BulkDelete\|objects\.Copy\|objects\.Create\|objects\.CreateTempURL\|objects\.Delete\|objects\.Download\|objects\.Get\|objects\.Update\|openstack\.Authenticate\|openstack\.AuthenticatedClient\|openstack\.AuthenticateV2\|openstack\.AuthenticateV3\|orders\.Create\|orders\.Delete\|orders\.Get\|osinherit\.Assign\|osinherit\.Unassign\|osinherit\.Validate\|pagination\.Request\|peers\.Create\|peers\.Delete\|peers\.Get\|peers\.Update\|policies\.Create\|policies\.Delete\|policies\.Get\|policies\.InsertRule\|policies\.RemoveRule\|policies\.Update\|pools\.BatchUpdateMembers\|pools\.Create\|pools\.CreateMember\|pools\.Delete\|pools\.DeleteMember\|pools\.Get\|pools\.GetMember\|pools\.Update\|pools\.UpdateMember\|portforwarding\.Create\|portforwarding\.Delete\|portforwarding\.Get\|portforwarding\.Update\|ports\.Create\|ports\.Delete\|ports\.Get\|ports\.Update\|projectendpoints\.Create\|projectendpoints\.Delete\|projects\.Create\|projects\.Delete\|projects\.DeleteTags\|projects\.Get\|projects\.ListTags\|projects\.ModifyTags\|projects\.Update\|qos\.Associate\|qos\.Create\|qos\.Delete\|qos\.DeleteKeys\|qos\.Disassociate\|qos\.DisassociateAll\|qos\.Get\|qos\.Update\|queues\.Create\|queues\.Delete\|queues\.Get\|queues\.GetStats\|queues\.Purge\|queues\.Share\|queues\.Update\|quotas\.Create\|quotasets\.Delete\|quotasets\.Get\|quotasets\.GetDefaults\|quotasets\.GetDetail\|quotasets\.GetUsage\|quotasets\.Update\|quotas\.Get\|quotas\.GetDetail\|quotas\.Update\|rbacpolicies\.Create\|rbacpolicies\.Delete\|rbacpolicies\.Get\|rbacpolicies\.Update\|recordsets\.Create\|recordsets\.Delete\|recordsets\.Get\|recordsets\.Update\|regions\.Create\|regions\.Delete\|regions\.Get\|regions\.Update\|registeredlimits\.BatchCreate\|registeredlimits\.Delete\|registeredlimits\.Get\|registeredlimits\.Update\|remoteconsoles\.Create\|replicas\.Create\|replicas\.Delete\|replicas\.ForceDelete\|replicas\.Get\|replicas\.GetExportLocation\|replicas\.ListExportLocations\|replicas\.Promote\|replicas\.ResetState\|replicas\.ResetStatus\|replicas\.Resync\|request\.Create\|request\.Delete\|request\.Get\|request\.Update\|resourceproviders\.Create\|resourceproviders\.Delete\|resourceproviders\.Get\|resourceproviders\.GetAllocations\|resourceproviders\.GetInventories\|resourceproviders\.GetTraits\|resourceproviders\.GetUsages\|resourceproviders\.Update\|resourcetypes\.GenerateTemplate\|resourcetypes\.GetSchema\|resourcetypes\.List\|roles\.AddUser\|roles\.Assign\|roles\.Create\|roles\.CreateRoleInferenceRule\|roles\.Delete\|roles\.DeleteRoleInferenceRule\|roles\.DeleteUser\|roles\.Get\|roles\.GetRoleInferenceRule\|roles\.ListRoleInferenceRules\|roles\.Unassign\|roles\.Update\|routers\.AddInterface\|routers\.Create\|routers\.Delete\|routers\.Get\|routers\.RemoveInterface\|routers\.Update\|rules\.Create\|rules\.CreateBandwidthLimitRule\|rules\.CreateDSCPMarkingRule\|rules\.CreateMinimumBandwidthRule\|rules\.Delete\|rules\.DeleteBandwidthLimitRule\|rules\.DeleteDSCPMarkingRule\|rules\.DeleteMinimumBandwidthRule\|rules\.Get\|rules\.GetBandwidthLimitRule\|rules\.GetDSCPMarkingRule\|rules\.GetMinimumBandwidthRule\|rules\.Update\|rules\.UpdateBandwidthLimitRule\|rules\.UpdateDSCPMarkingRule\|rules\.UpdateMinimumBandwidthRule\|ruletypes\.GetRuleType\|secgroups\.AddServer\|secgroups\.Create\|secgroups\.CreateRule\|secgroups\.Delete\|secgroups\.DeleteRule\|secgroups\.Get\|secgroups\.RemoveServer\|secgroups\.Update\|secrets\.Create\|secrets\.CreateMetadata\|secrets\.CreateMetadatum\|secrets\.Delete\|secrets\.DeleteMetadatum\|secrets\.Get\|secrets\.GetMetadata\|secrets\.GetMetadatum\|secrets\.GetPayload\|secrets\.Update\|secrets\.UpdateMetadatum\|securityservices\.Create\|securityservices\.Delete\|securityservices\.Get\|securityservices\.Update\|servergroups\.Create\|servergroups\.Delete\|servergroups\.Get\|servers\.ChangeAdminPassword\|servers\.ConfirmResize\|servers\.Create\|servers\.CreateImage\|servers\.CreateMetadatum\|servers\.Delete\|servers\.DeleteMetadatum\|servers\.Evacuate\|servers\.ForceDelete\|servers\.Get\|servers\.GetPassword\|servers\.InjectNetworkInfo\|servers\.LiveMigrate\|servers\.Lock\|servers\.Metadata\|servers\.Metadatum\|servers\.Migrate\|servers\.Pause\|servers\.Reboot\|servers\.Rebuild\|servers\.Rescue\|servers\.ResetMetadata\|servers\.ResetNetwork\|servers\.ResetState\|servers\.Resize\|servers\.Resume\|servers\.RevertResize\|servers\.Shelve\|servers\.ShelveOffload\|servers\.ShowConsoleOutput\|servers\.Start\|servers\.Stop\|servers\.Suspend\|servers\.Unlock\|servers\.Unpause\|servers\.Unrescue\|servers\.Unshelve\|servers\.Update\|servers\.UpdateMetadata\|servers\.WaitForStatus\|services\.Create\|services\.Delete\|services\.Get\|services\.Update\|shareaccessrules\.Get\|shareaccessrules\.List\|sharenetworks\.AddSecurityService\|sharenetworks\.Create\|sharenetworks\.Delete\|sharenetworks\.Get\|sharenetworks\.RemoveSecurityService\|sharenetworks\.Update\|shares\.Create\|shares\.Delete\|shares\.DeleteMetadatum\|shares\.Extend\|shares\.ForceDelete\|shares\.Get\|shares\.GetExportLocation\|shares\.GetMetadata\|shares\.GetMetadatum\|shares\.GrantAccess\|shares\.ListAccessRights\|shares\.ListExportLocations\|shares\.ResetStatus\|shares\.Revert\|shares\.RevokeAccess\|shares\.SetMetadata\|shares\.Shrink\|shares\.Unmanage\|shares\.Update\|shares\.UpdateMetadata\|sharetransfers\.Accept\|sharetransfers\.Create\|sharetransfers\.Delete\|sharetransfers\.Get\|sharetypes\.AddAccess\|sharetypes\.Create\|sharetypes\.Delete\|sharetypes\.GetDefault\|sharetypes\.GetExtraSpecs\|sharetypes\.RemoveAccess\|sharetypes\.SetExtraSpecs\|sharetypes\.ShowAccess\|sharetypes\.UnsetExtraSpecs\|siteconnections\.Create\|siteconnections\.Delete\|siteconnections\.Get\|siteconnections\.Update\|snapshots\.Create\|snapshots\.Delete\|snapshots\.ForceDelete\|snapshots\.Get\|snapshots\.ResetStatus\|snapshots\.Update\|snapshots\.UpdateMetadata\|snapshots\.UpdateStatus\|snapshots\.WaitForStatus\|speakers\.AddBGPPeer\|speakers\.AddGatewayNetwork\|speakers\.Create\|speakers\.Delete\|speakers\.Get\|speakers\.RemoveBGPPeer\|speakers\.RemoveGatewayNetwork\|speakers\.Update\|stackevents\.Find\|stackevents\.Get\|stackresources\.Find\|stackresources\.Get\|stackresources\.MarkUnhealthy\|stackresources\.Metadata\|stackresources\.Schema\|stackresources\.Template\|stacks\.Abandon\|stacks\.Adopt\|stacks\.Create\|stacks\.Delete\|stacks\.Find\|stacks\.Get\|stacks\.Preview\|stacks\.Update\|stacks\.UpdatePatch\|stacktemplates\.Get\|stacktemplates\.Validate\|subnetpools\.Create\|subnetpools\.Delete\|subnetpools\.Get\|subnetpools\.Update\|subnets\.Create\|subnets\.Delete\|subnets\.Get\|subnets\.Update\|swauth\.Auth\|swauth\.NewObjectStorageV1\|tags\.Add\|tags\.Check\|tags\.Delete\|tags\.DeleteAll\|tags\.List\|tags\.ReplaceAll\|tasks\.Create\|tasks\.Get\|tenants\.Create\|tenants\.Delete\|tenants\.Get\|tenants\.Update\|tokens\.Create\|tokens\.Get\|tokens\.Revoke\|tokens\.Validate\|transfers\.Accept\|transfers\.Create\|transfers\.Delete\|transfers\.Get\|trunks\.AddSubports\|trunks\.Create\|trunks\.Delete\|trunks\.Get\|trunks\.GetSubports\|trunks\.RemoveSubports\|trunks\.Update\|trusts\.CheckRole\|trusts\.Create\|trusts\.Delete\|trusts\.Get\|trusts\.GetRole\|users\.AddToGroup\|users\.ChangePassword\|users\.Create\|users\.Delete\|users\.Get\|users\.IsMemberOfGroup\|users\.RemoveFromGroup\|users\.Update\|utils\.ChooseVersion\|utils\.GetSupportedMicroversions\|utils\.RequireMicroversion\|volumeattach\.Create\|volumeattach\.Delete\|volumeattach\.Get\|volumes\.Attach\|volumes\.BeginDetaching\|volumes\.ChangeType\|volumes\.Create\|volumes\.Delete\|volumes\.Detach\|volumes\.ExtendSize\|volumes\.ForceDelete\|volumes\.Get\|volumes\.InitializeConnection\|volumes\.ReImage\|volumes\.Reserve\|volumes\.ResetStatus\|volumes\.SetBootable\|volumes\.SetImageMetadata\|volumes\.TerminateConnection\|volumes\.Unreserve\|volumes\.Update\|volumes\.UploadImage\|volumes\.WaitForStatus\|volumetypes\.AddAccess\|volumetypes\.Create\|volumetypes\.CreateEncryption\|volumetypes\.CreateExtraSpecs\|volumetypes\.Delete\|volumetypes\.DeleteEncryption\|volumetypes\.DeleteExtraSpec\|volumetypes\.Get\|volumetypes\.GetEncryption\|volumetypes\.GetEncryptionSpec\|volumetypes\.GetExtraSpec\|volumetypes\.ListExtraSpecs\|volumetypes\.RemoveAccess\|volumetypes\.Update\|volumetypes\.UpdateEncryption\|volumetypes\.UpdateExtraSpec\|workflows\.Create\|workflows\.Delete\|workflows\.Get\|zones\.Create\|zones\.Delete\|zones\.Get\|zones\.Update\)(#\1(context.TODO(), #g
        s#\(\.AllPages(\)#\1context.TODO(), #g
        s#\(\.EachPage(\)\(func(\)#\1context.TODO(), \2ctx context.Context, #g

        # 8: Rename identifiers that were changed in v2
        s#\(\(volumes\|servers\)\.SchedulerHint\)s#\2.SchedulerHintOpts#g

        # 9: Tentatively replace error handling
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault400); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusBadRequest) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault401); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusUnauthorized) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault403); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusForbidden) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault404); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusNotFound) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault405); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusMethodNotAllowed) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault408); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusRequestTimeout) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault409); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusConflict) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault429); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusTooManyRequests) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault500); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusInternalServerError) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault502); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusBadGateway) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault503); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusServiceUnavailable) {#g
        s#\(\t\+\)if _, ok := err.(gophercloud.ErrDefault504); \(!\?\)ok {#\1if \2gophercloud.ResponseCodeIs(err, http.StatusGatewayTimeout) {#g
    }
    ' {} \;

grep -r -l 'context\.TODO' | xargs -r sed -i '
    /^import ($/ a  "context"
    '

grep -r -l 'http\.Status' | xargs -r sed -i '
    /^import ($/ a  "net/http"
    '

goimports -format-only -w .
go mod tidy
```
