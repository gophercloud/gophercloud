# clouds.yaml → auth package bridge

Status: approved for planning
Branch: keystone-v3-auth
Date: 2026-09-04

## Problem

`keystone-v3-auth` replaced the old v2 gophercloud auth path with a new
`auth` package (`auth.Authenticator`, `auth.AuthOptionsV2`/`AuthOptionsV3`,
per-mechanism opts like `V3PasswordOpts`, `V3ApplicationCredentialOpts`,
etc.). `openstack.AuthenticatedClient` and `config.NewProviderClient` now
require `auth.Authenticator`.

`openstack/config/clouds.Parse()` (unmodified by this branch) still parses
`clouds.yaml`/`secure.yaml`/`clouds-public.yaml` and flattens the result into
the legacy `gophercloud.AuthOptions` struct, which has no `Authenticate`/
`GetAuthURL` methods and does not satisfy `auth.Authenticator`. There is no
code anywhere that bridges `clouds.Parse()`'s output into the new `auth`
package. The package doc comment in `clouds.go` still shows the old,
now-nonfunctional composition. `gophercloud.AuthOptions` itself is being
retired, so the fix must not route through it.

## Goal

Make `clouds.yaml`-based authentication work against the new `auth`
package, with authentication-selection logic living entirely in `auth`,
and `openstack/config/clouds` doing file I/O/parsing only.

## Package boundary

- `openstack/config/clouds` is responsible only for locating and parsing
  `clouds.yaml`/`secure.yaml`/`clouds-public.yaml`, merging cloud entries
  (profile/public-cloud merge, secure.yaml merge), resolving TLS config,
  and resolving `gophercloud.EndpointOpts` (region/interface). It has no
  opinion on which auth mechanism to use or how to build a token request.
- `auth` owns everything about turning parsed cloud data into an
  `auth.Authenticator`: mechanism selection, scope resolution, and any
  runtime overrides of credential fields.

## `clouds.Parse()` signature change

```go
func Parse(opts ...ParseOption) (Cloud, gophercloud.EndpointOpts, *tls.Config, error)
```

`Parse()` returns the merged `Cloud` struct (as defined in
`openstack/config/clouds/types.go`) directly, instead of flattening it into
`gophercloud.AuthOptions`. Internal merge logic (`mergeClouds`,
`mergeWithPublicClouds`, `computeTLSConfig`, `computeAvailability`) is
unchanged.

### `ParseOption`s that stay in `clouds`

`WithCloudName`, `WithLocations`, `WithPublicLocations`, `WithCloudsYAML`,
`WithSecureYAML`, `WithCloudsPublicYAML`, `WithRegion`, `WithEndpointType`,
`WithIdentityEndpoint`, `WithCACertPath`, `WithClientCertPath`,
`WithClientKeyPath`, `WithInsecure`.

`WithIdentityEndpoint` stays here even though it writes into
`Cloud.AuthInfo.AuthURL` — it's an endpoint locator (like region), not
credential material.

### `ParseOption`s removed from `clouds` (moved to `auth`)

`WithUsername`, `WithUserID`, `WithPassword`, `WithToken`, `WithDomainID`,
`WithDomainName`, `WithProjectID`, `WithProjectName`,
`WithApplicationCredentialID`, `WithApplicationCredentialName`,
`WithApplicationCredentialSecret`, `WithScope`. No backwards-compatibility
aliases are kept in `clouds` for these — this is pre-1.0, breaking-change
territory the branch already embraces.

## New `auth/clouds.go`

```go
func AuthOptionsFromCloud(cloud clouds.Cloud, opts ...CloudOption) (Authenticator, error)
func AuthOptionsFromCloudV2(cloud clouds.Cloud, opts ...CloudOption) (*AuthOptionsV2, error)
func AuthOptionsFromCloudV3(cloud clouds.Cloud, opts ...CloudOption) (*AuthOptionsV3, error)
```

This mirrors the existing `AuthOptionsFromEnv`/`AuthOptionsFromEnvV2`/
`AuthOptionsFromEnvV3` triad in `auth/env.go`. `AuthOptionsFromCloud` picks
V2 vs V3 based on `cloud.IdentityAPIVersion == "2.0"`, defaulting to V3 —
the same convention `AuthOptionsFromEnv` already uses for
`OS_IDENTITY_API_VERSION`.

`auth` importing `openstack/config/clouds` does not create an import
cycle: `clouds` only imports the root `gophercloud` package and the YAML
library.

### `CloudOption`

A new functional-option type in `auth`, applied on top of the values read
from `cloud.AuthInfo` before mechanism/scope selection runs:

- `WithUsername(string)`
- `WithUserID(string)`
- `WithPassword(string)`
- `WithToken(string)`
- `WithPasscode(string)` — new; `clouds.AuthInfo` has no passcode field
  (a one-time TOTP code has no business living in a static file), so this
  is the only way to add TOTP on top of an otherwise file-sourced config.
- `WithDomainID(string)` / `WithDomainName(string)`
- `WithProjectID(string)` / `WithProjectName(string)`
- `WithApplicationCredentialID(string)` / `WithApplicationCredentialName(string)` / `WithApplicationCredentialSecret(string)`
- `WithScope(*Scope)` — short-circuits scope resolution entirely, same as
  today's `clouds.WithScope`.

## Mechanism and version selection

1. If `cloud.AuthType` is one of the version-specific values
   (`v2password`, `v2token`, `v3password`, `v3token`,
   `v3applicationcredential`), use it directly — it overrides
   `IdentityAPIVersion`.
2. If `cloud.AuthType` is a version-agnostic value (`password`, `token`),
   resolve V2 vs V3 via `IdentityAPIVersion`, then pick the mechanism.
3. If `cloud.AuthType` is empty, resolve V2 vs V3 via `IdentityAPIVersion`
   (default V3), then infer the mechanism from populated fields using the
   same relative precedence as `auth/env.go`'s `AuthOptionsFromEnvV3`:
   password > token > application credential > passcode (passcode can
   only be populated via `WithPasscode`, so it is checked last and only
   matters when explicitly supplied).
4. If nothing resolves, return an error (`gophercloud.ErrMissingInput{Argument: "Auth"}`).

## Scope resolution

`clouds.AuthInfo` already keeps `UserDomainID`/`UserDomainName`,
`ProjectDomainID`/`ProjectDomainName`, generic `DomainID`/`DomainName`, and
`DefaultDomain` as distinct fields, and `auth.Scope` has separate
`ProjectDomainID`/`ProjectDomainName` fields. Unlike today's
`gophercloud.AuthOptions`-based flattening (which conflates user and
project domain into one field), the new mapping keeps them distinct:

```
userDomainID, userDomainName       := coalesce(UserDomainID, DomainID), coalesce(UserDomainName, DomainName)
    if both empty: userDomainID = DefaultDomain
projectDomainID, projectDomainName := coalesce(ProjectDomainID, DomainID), coalesce(ProjectDomainName, DomainName)
    if both empty: projectDomainID = DefaultDomain
```

`SystemScope` (currently silently ignored by `clouds.Parse()`), `TrustID`,
`ProjectID`, and `ProjectName` map straight into `auth.Scope`. `WithScope`
bypasses all of this and is used verbatim. Application-credential auth
never builds a scope (matches existing `V3ApplicationCredentialOpts`
behavior — app creds carry an implicit scope in Keystone).

## Known, called-out behavior changes

- `AuthInfo.AllowReauth` becomes a no-op. The new `auth` opts structs
  hardcode reauth-ability per mechanism (password/app-cred: reauthable;
  token/TOTP: not), so there is nothing to plumb through anymore.
- TOTP cannot be sourced from a `clouds.yaml` file itself — only via
  `WithPasscode`.

## Doc/example updates

`openstack/config/clouds/clouds.go`'s package-level doc comment currently
shows:

```go
ao, eo, tlsConfig, err := clouds.Parse()
providerClient, err := config.NewProviderClient(ctx, ao, config.WithTLSConfig(tlsConfig))
```

which no longer compiles (`ao` doesn't satisfy `auth.Authenticator`).
Update to:

```go
cloud, eo, tlsConfig, err := clouds.Parse()
authOpts, err := auth.AuthOptionsFromCloud(cloud)
providerClient, err := config.NewProviderClient(ctx, authOpts, config.WithTLSConfig(tlsConfig))
```

## Test plan

- `openstack/config/clouds/clouds_test.go`: update assertions from
  `ao.Field` to `cloud.AuthInfo.Field` (mechanical — no behavior change to
  the parsing/merging logic itself). Remove/relocate the subtests that
  exercised the now-removed auth-field `With*` options
  (`ExampleWithUserID` moves to `auth/testing`).
- New `auth/testing/clouds_test.go`: mechanism selection (explicit
  `AuthType` for every value, inferred from fields, V2 vs V3 via
  `IdentityAPIVersion`), scope resolution (user vs. project domain
  separation, `DefaultDomain` fallback, system scope, trust ID,
  application-credential no-scope), every `CloudOption` override, and the
  two error cases (missing `AuthURL`, no usable credentials).

## Non-goals

- No change to `computeTLSConfig`, `mergeClouds`, `mergeWithPublicClouds`,
  or any other parsing/merge logic in `clouds.go`.
- No combined one-shot "parse and authenticate" convenience function —
  callers compose `clouds.Parse()` and `auth.AuthOptionsFromCloud()`
  themselves, same as the doc example.
- No V2 `clouds.yaml` acceptance/integration test beyond unit coverage —
  v2 Keystone is effectively dead in practice; V2 support here exists for
  API symmetry with `auth.AuthOptionsFromEnv`.
