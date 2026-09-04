# Gophercloud review checklist

Copy this list into the review notes and check items that apply. Skip items that are N/A and say so.

## Layout and naming

- [ ] Package lives at `openstack/<service>/<version>/<resource>/`.
- [ ] Files: `requests.go` (HTTP + opts), `results.go` (response types + Extract), `urls.go` (paths), `microversions.go` only when response types change incompatibly.
- [ ] Unit tests in `testing/`; acceptance tests in `internal/acceptance/openstack/<service>/`.
- [ ] Result receiver `r`; unmarshalled value `s`; request return value `r` except pagers.
- [ ] Body opts: `<Action>OptsBuilder` with `To<Resource><Action>Map` (e.g. `ToServerCreateMap`).
- [ ] Query opts: `<Action>OptsBuilder` with `To<Resource><Action>Query`.
- [ ] `BuildRequestBody(opts, "<resource>")` uses the correct JSON envelope key.
- [ ] Imports: stdlib, then external, then `github.com/gophercloud/gophercloud/v2/...`, groups separated by blank lines.
- [ ] Module path is `github.com/gophercloud/gophercloud/v2` (the `/v2` suffix).

## Requests

- [ ] Every exported API function takes `context.Context` as the first argument after the client (match neighboring functions).
- [ ] Required fields have **no** `omitempty`. Optional fields **always** have `omitempty`.
- [ ] New microversion request fields: `omitempty` **and** a GoDoc line stating the minimum microversion. Prefer pointers when `false`/`0`/`""` must be distinguishable from omit.
- [ ] Do not invent client-side enum validation unless the package already does that for the same resource.
- [ ] List operations use `pagination.Pager` / `EachPage` consistently with similar resources.
- [ ] URL helpers in `urls.go` match the service's actual paths (pluralization, nested IDs).
- [ ] Delete/Get/Update IDs are passed in the path, not duplicated incorrectly in the body.

## Results

- [ ] JSON tags match OpenStack field names exactly (no invented aliases).
- [ ] New microversion **response** fields are **pointers** so callers can nil-check.
- [ ] Incompatible response type changes go in `microversions.go` with new structs and Extract methods — do not silently change the existing struct's type.
- [ ] Every struct field has a GoDoc comment. Microversion-dependent fields document the required version.
- [ ] `Extract` / `ExtractErr` / `Extract...` match sibling resources (`gophercloud.Result` vs `ErrResult`).
- [ ] Time fields use the same type as nearby resources (`time.Time` vs `gophercloud.JSONRFC3339Milli` etc.). Wrong time type is a **blocker**.
- [ ] List `Page` types implement `NextPageURL` / `IsEmpty` like sibling resources.
- [ ] If the list JSON object has more than one top-level key besides `*links`, the `Page` type implements `pagination.KeyedPage` with `ResourceKey()` set to the collection envelope (e.g. `"networks"`, `"inventories"`). `AllPages` otherwise scans for the first non-links array key, which is wrong or unstable when extra keys exist ([gophercloud#3947](https://github.com/gophercloud/gophercloud/pull/3947)). Missing `ResourceKey()` is a **blocker** when more than one top-level value is an array, otherwise a **warning**.
- [ ] `ResourceKey()` matches the JSON array field that Extract methods read.

## Microversions

See `docs/MICROVERSIONS.md`.

- [ ] New request field → `omitempty` + GoDoc version.
- [ ] New response field → pointer + GoDoc version.
- [ ] Changed response type → new types in `microversions.go`.
- [ ] New enum **value** on an existing string field usually needs **no** Gophercloud change unless the SDK restricts allowed values.
- [ ] Tests that need a microversion must set `client.Microversion` on the fake/service client.

## Unit tests (`testing/`)

- [ ] `fakeServer := th.SetupHTTP()` and `defer fakeServer.Teardown()`.
- [ ] Mock asserts method, path, headers (`X-Auth-Token`, `Content-Type`, microversion header when relevant), and body.
- [ ] **Every** field on Create/Update/List opts appears in at least one request test.
- [ ] **Every** field on result structs is asserted in a get/list extract test (including pointer/microversion fields).
- [ ] Use `th.Assert*` (fatal) / `th.Check*` (non-fatal) from `testhelper`; expected first.
- [ ] Fixtures live in `fixtures_test.go` (or the package's existing fixture pattern). Do not leave response JSON that omits new fields.

## Acceptance tests

- [ ] Path: `internal/acceptance/openstack/<service>/...`
- [ ] Cover the new operation and its option variants, not only create/delete smoke.
- [ ] Tests clean up resources; note any admin-only or charge-incurring calls.
- [ ] Skip/guard patterns must match the rest of that service's acceptance package.

## Bugs and correctness

- [ ] Nil pointer dereference on optional nested structs.
- [ ] `json:",omitempty"` on slices/maps: empty vs omitted semantics vs the API.
- [ ] `*bool` / `*int` needed when `false`/`0` is a valid distinct payload.
- [ ] Error results: return `err` directly; wrap only when adding useful context (project convention is usually not to wrap).
- [ ] Pagination: next-page URL / marker handling vs the service; `AllPages` uses `ResourceKey()` when the page implements `KeyedPage`.
- [ ] Identity vs project-scoped endpoints if the service is special-cased in `openstack/` client constructors.

## API coverage gaps

When reviewing a package (not only a tiny diff):

- [ ] List operations in OpenStack (index, create, get, update, delete, actions).
- [ ] List query parameters and body fields from service source on `stable/*`.
- [ ] Report anything implemented in OpenStack but missing in Gophercloud under **Suggestions** (or **Warnings** if the PR/package claimed to add it).
- [ ] Report extra Gophercloud fields that the service never accepted/returned.
- [ ] Proof links are GitHub (or opendev) **source** on a **non-master** branch, e.g. `https://github.com/openstack/nova/blob/stable/2025.1/...`. Docs-only links are insufficient.

## PR hygiene (PRs only)

- [ ] GitHub issue exists; description has `For #<number>`.
- [ ] Description links OpenStack source (stable branch) proving request/response fields.
- [ ] Scope is focused; unrelated refactors called out.
- [ ] Title `[wip]` / `[Pending #PRNUM]` used correctly if applicable.
- [ ] Reminder: CI expects a `semver:patch`, `semver:minor`, or `semver:major` label.
- [ ] Do not tell the author to squash during review unless they asked.
