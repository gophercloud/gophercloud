# apicheck

`apicheck` measures how much of the OpenStack API surface Gophercloud covers. It
diffs the authoritative [OpenStack OpenAPI schemas][openstack-openapi] against
what the Gophercloud source actually implements and reports the gaps: missing
operations, missing request/response fields, and missing query parameters — each
annotated with the microversion that introduced it.

It is a standalone Go module (its own `go.mod`) so the OpenAPI-parser
dependencies never leak into the dependency-light main SDK module.

[openstack-openapi]: https://github.com/gtema/openstack

## What it does

```
OpenAPI YAML ──▶ spec model ──┐
                              ├──▶ diff ──▶ report (Markdown or JSON)
Go source (AST) ─▶ impl model ┘
```

- **spec/** parses each service's OpenAPI 3.1 document, flattening the
  microversion `oneOf` unions into a single field set and tagging every field and
  operation with its `min-ver`/`max-ver`.
- **impl/** statically analyses the Gophercloud packages (via
  `golang.org/x/tools/go/packages` + `go/types`) to discover operations and their
  fields. Static analysis — not reflection — is required because Gophercloud
  builds many request bodies by hand (`map[string]any{...}`, `b["key"] = ...`),
  composes URLs through helper functions, and uses `json:"-"` fields, none of
  which reflection can see.
- **diff/** matches spec operations to impl operations by a normalised
  `METHOD /path[:action]` key (path parameters collapsed to `{}`) and reports the
  set-difference of their fields.
- **report/** renders Markdown (human) or JSON (machine/CI).

## Usage

The OpenAPI schemas are **not** vendored into this repository. Clone the
[gtema/openstack][openstack-openapi] repository and point the tool at its
`openstack_types/data` directory:

```bash
git clone https://github.com/gtema/openstack
make apicheck APICHECK_SPEC_ROOT=<clone>/openstack_types/data
```

Forward extra flags with `APICHECK_ARGS`:

```bash
# One service, JSON output
make apicheck APICHECK_SPEC_ROOT=<clone>/openstack_types/data APICHECK_ARGS="-service compute -format json"
```

Or run the tool directly from `tools/apicheck`:

```bash
cd tools/apicheck
go run . -spec-root <clone>/openstack_types/data
```

There is no default for the schema location, so `-spec-root` (or
`APICHECK_SPEC_ROOT`) is required; the tool exits with a reminder if it is
missing.

### Flags

| Flag | Description |
|---|---|
| `-config` | Path to the config file (default `apicheck.yaml`). |
| `-spec-root` | Override `spec_root` from the config. |
| `-gc-root` | Override `gophercloud_root` from the config. |
| `-service` | Limit the run to a single service (e.g. `compute`). |
| `-format` | `markdown` (default) or `json`. |
| `-dump-spec` | Dump the parsed spec model as JSON and exit. |
| `-dump-impl` | Dump the parsed impl model as JSON and exit. |

### Configuration

`apicheck.yaml` maps each Gophercloud service package to its OpenAPI file and the
version prefix to strip from spec paths (e.g. `/v2.1` for compute).

`spec_root` and `gophercloud_root` are intentionally left unset (`null`) in the
committed config so that no machine-specific paths are checked in:

- `spec_root` has no portable default and must be supplied via `-spec-root` /
  `APICHECK_SPEC_ROOT` (see [Usage](#usage)).
- `gophercloud_root` defaults to this repository's root, resolved relative to the
  config file; override it with `-gc-root` only if you run the tool from an
  unusual location.

## Microversions

Following the project decision, coverage is measured against the **latest**
`vX.yaml` (the union of all microversions), and every gap is **annotated** with
the `min-ver` that introduced the operation or field. A microversion is derived
from either an `x-openstack.min-ver` extension or the `_<mv>` suffix on a
component schema name (e.g. `ServersCreate_294` → `2.94`).

## Baseline

`baseline.json` is a committed snapshot of the full JSON report — a point-in-time
audit of coverage across all services. Regenerate it with:

```bash
cd tools/apicheck
go run . -spec-root <data> -format json > baseline.json
```

To use apicheck as a **non-blocking CI regression gate**, run the tool in CI
(checking out openstack-openapi alongside this repo) and diff the fresh JSON
against `baseline.json`; fail only on *newly introduced* gaps. Refresh the
baseline as coverage improves.

## Known limitations

The report is deliberately honest about what it cannot see:

- **Field coverage is measured at the top (envelope) level.** The spec flattener
  descends the single request/response envelope (`{"server": {...}}`,
  `{"servers": [...]}`) but treats nested domain objects and arrays-of-objects
  (e.g. `block_device_mapping_v2` items) as single leaf fields, because
  Gophercloud models those as separate Go structs referenced by the parent field.
  Sub-fields of nested objects are therefore not individually diffed.
- **Response fields are matched against a package-wide pool.** Impl response
  structs cannot be reliably tied to a specific operation, so response-field
  coverage is approximated by the union of json-tagged fields across the package's
  non-request structs. This over-approximates (favouring false negatives over
  false positives).
- **`json:"-"` fields** are surfaced as `ManualHandled` and counted as present
  (they exist but are serialised by hand, so their wire name cannot be verified).
- **Unresolved URLs** (dynamically-built URLs the analyser cannot fold to a static
  path) are counted and reported on stderr, never silently dropped, so coverage
  numbers stay honest.
- **Unmatched impl operations** (impl ops with no matching spec op) are listed
  under a collapsible section; these are usually path-mapping mismatches or
  Gophercloud extensions absent from the schema.
- **block-storage** (Cinder) lists every path twice in its schema — once
  project-scoped (`/v3/{project_id}/volumes`) and once project-less
  (`/v3/volumes`). Gophercloud only builds the project-less form (the project ID
  is carried in the catalog endpoint), so the config sets `strip_project_id:
  true`, which collapses the leading `{project_id}` and dedupes the pairs.
  Without it the denominator doubles and coverage is halved (~20% instead of the
  real ~40%).
- **object-store** requires special handling because Swift's service-catalog
  endpoint already embeds `/v1/{account}`: the config strips that whole prefix,
  and the extractor understands `client.Endpoint` (the account root) and the
  generic `client.Request(ctx, "VERB", ...)` calls the Swift packages use. The
  one remaining gap (`DELETE /`, delete-account) is genuinely unimplemented, and
  `COPY` (a Swift extension verb) shows up as an unmatched impl operation because
  it is not part of the OpenAPI schema.

## Testing

```bash
make apicheck-test        # from the repo root
# or
cd tools/apicheck && go test ./...
```

Tests live in `testing/` and run the extractor and spec parser against
hand-written fixtures in `testing/fixtures/` (an isolated Go module mirroring the
Gophercloud request/URL/result patterns, plus a small OpenAPI document), covering:
`BuildRequestBody` struct bodies, `q:` query opts, `json:"-"` fields, manual
`map[string]any` action bodies, nested url-func resolution, the microversion
`oneOf` flattening, and the envelope-only field flattening.
