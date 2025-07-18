# Microversions

## Table of Contents

* [Introduction](#introduction)
* [Client Configuration](#client-configuration)
* [Gophercloud Developer Information](#gophercloud-developer-information)
* [Application Developer Information](#application-developer-information)

## Introduction

Microversions are an OpenStack API ability which allows developers to add and
remove features while still retaining backwards compatibility for all prior
versions of the API.

More information can be found here:

> Note: these links are not an exhaustive reference for microversions.

* http://specs.openstack.org/openstack/api-wg/guidelines/microversion_specification.html
* https://developer.openstack.org/api-guide/compute/microversions.html
* https://github.com/openstack/keystoneauth/blob/master/doc/source/using-sessions.rst

## Client Configuration

You can set a specific microversion on a Service Client by doing the following:

```go
client, err := openstack.NewComputeV2(context.TODO(), providerClient, nil)
client.Microversion = "2.52"
```

## Gophercloud Developer Information

Microversions change several aspects about API interaction.

### Existing Fields, New Values

This is when an existing field behaves like an "enum" and a new valid value
is possible by setting the client's microversion to a specific version.

An example of this can be seen with Nova/Compute's Server Group `policy` field
and the introduction of the [`soft-affinity`](https://developer.openstack.org/api-ref/compute/?expanded=create-server-group-detail#create-server-group)
value.

Unless Gophercloud is limiting the valid values that are passed to the
Nova/Compute service, no changes are required in Gophercloud.

### New Request Fields

This is when a microversion enables a new field to be used in an API request.
When implementing this kind of change, it is imperative that the field has
the `omitempty` attribute set. If `omitempty` is not set, then the field will
be used for _all_ microversions and possibly cause an error from the API
service. You may need to use a pointer field in order for this to work.

When adding a new field, please make sure to include a GoDoc comment about
what microversions the field is valid for.

Please see [here](https://github.com/gophercloud/gophercloud/blob/917735ee91e24fe1493e57869c3b42ee89bc95d8/openstack/compute/v2/servers/requests.go#L215-L217) for an example.

### New Response Fields

This is when a microversion includes new fields in the API response. The
correct way of implementing this in Gophercloud is to add the field to the
resource's "result" struct (in the `results.go` file) as a *pointer*. This
way, the developer can check for a `nil` value to see if the field was set
from a microversioned result.

When adding a new field, please make sure to include a GoDoc comment about
what microversions the field is valid for.

Please see [here](https://github.com/gophercloud/gophercloud/blob/ed4deec00ff1d4d4c8a762af0c6360d4184a4bf4/openstack/compute/v2/servers/results.go#L221-L223) for an example.

### Modified Response Fields

This is when the new type of the returned field is incompatible with the
original type. When this happens, an entire new result struct must be
created with new Extract methods to account for both the original result
struct and new result struct.

These new structs and methods need to be defined in a new `microversions.go`
file.

Please see [here](https://github.com/gophercloud/gophercloud/blob/917735ee91e24fe1493e57869c3b42ee89bc95d8/openstack/container/v1/capsules/microversions.go) for an example.

## Application Developer Information

Gophercloud does not perform any validation checks on the API request to make
sure it is valid for a specific microversion. It is up to you to ensure that
the API request is using the correct fields and functions for the microversion.
