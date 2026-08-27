---
name: review-gophercloud
description: Review Gophercloud Go SDK code for API coverage gaps, convention violations, microversion mistakes, missing tests, and bugs. Use when the user invokes this skill or asks to review Gophercloud packages, PRs, diffs, or OpenStack API implementations.
disable-model-invocation: true
---

# Review Gophercloud

Review Gophercloud (`github.com/gophercloud/gophercloud/v2`) as a core contributor would. Do not implement fixes unless the user asks.

Read [checklist.md](checklist.md) before reporting findings. Read repo docs only as needed: `docs/STYLEGUIDE.md`, `docs/MICROVERSIONS.md`, `.github/CONTRIBUTING.md`.

## Scope

Determine the target, then review only that:

| User asks for | Target |
|---|---|
| PR, branch, or diff | Changed files vs `main` (committed + uncommitted) |
| A package or service | That package tree, including `testing/` and matching acceptance tests |
| New API / resource | Completeness of `requests.go`, `results.go`, `urls.go`, `microversions.go`, unit tests, acceptance tests |
| Unspecified | Ask once whether to review the current diff or a named package |

Always cover: correctness/bugs, Gophercloud conventions, microversions, tests, and API coverage gaps vs OpenStack.

## Workflow

1. Identify service, version, and resource (`openstack/<service>/<version>/<resource>/`).
2. Read `requests.go`, `results.go`, `urls.go`, `doc.go`, and `microversions.go` if present.
3. Read unit tests under `testing/` and acceptance tests under `internal/acceptance/openstack/<service>/`.
4. For API completeness, compare Go types to **OpenStack source on a non-`master` branch** (e.g. `stable/2025.1`), not API docs alone. Link the proving source in findings.
5. Check tests actually exercise **every** request and response field, not a happy-path subset.
6. Run the package unit tests when reviewing local changes:
   ```bash
   go test ./openstack/<service>/<version>/<resource>/...
   ```
7. Report findings. Do not commit, push, or open a PR.

Gophercloud does **not** validate microversion compatibility. Call out missing `omitempty`, missing pointer response fields, or missing GoDoc version notes — not "the client should reject the field."

## Output

Use this structure. Omit a section only when it has no items (do not invent filler). Put **file and line** on every finding that maps to code (`path:line` in the heading, plus a code citation when quoting). Package-level or PR-description issues can omit a line.

```markdown
## Summary
What was reviewed (diff, package, or PR) and the overall verdict in 2–4 sentences.

## Blockers
Must change before merge or before treating the API as complete.

- `openstack/.../requests.go:215` — one-sentence problem and why it breaks callers or the API.
  Include an OpenStack source link when claiming a field/path/method mismatch.

## Warnings
Should change; likely bugs, missing tests that would miss the bug, or incomplete coverage of claimed work.

- `openstack/.../results.go:88` — ...

## Suggestions
Optional improvements: style, GoDoc, extra tests, API gaps the PR did not claim to add, PR hygiene (`For #N`, stable-branch proof links, `semver:*`).

- `openstack/.../testing/requests_test.go:40` — ...

## Positives
What is already in good shape (patterns followed, tests that do cover fields, correct microversion handling). Be specific; cite file and line when pointing at a good example.
```

Classification:

- **Blockers**: wrong HTTP method/path, broken Extract, data loss, panic, auth/token mishandling, missing `omitempty` on a microversioned request field, required field omitted, incorrect JSON names, wrong time type, list `Page` with multiple top-level JSON keys and no `ResourceKey()` when `AllPages` would pick the wrong collection
- **Warnings**: tests that cannot catch the bug, incomplete field coverage for new/claimed APIs, missing acceptance test for a new operation, naming/OptsBuilder drift, pointer/`omitempty` uncertainty, list body with extra top-level keys but no `ResourceKey()` yet (heuristic may still work)
- **Suggestions**: style, import order, comment nits, extra OpenStack operations/fields not in scope, PR description gaps
- **Positives**: always include at least one concrete item when anything is correct; if the review is all blockers, still note any solid parts

If there are no blockers, say the code is not blocked and keep warnings/suggestions/positives as applicable. Residual risk (no live API, untested microversions) goes under Warnings.

## Rules of the review

- Follow existing package patterns over inventing new ones.
- List `Page` types whose JSON body has more than one top-level key (besides `*links`) must implement `pagination.KeyedPage` (`ResourceKey()`), so `AllPages` does not guess the collection key. See `pagination/pager.go` and [gophercloud#3947](https://github.com/gophercloud/gophercloud/pull/3947).
- Assertion order is **expected, actual**: `th.AssertEquals(t, expected, actual)`.
- Cite code with `startLine:endLine:path`.
- Do not demand Gophercloud validate microversions.
- Do not treat API documentation as sufficient proof; prefer service source on `stable/*`.
- Do not squash-commit advice unless the user is writing the PR description.
