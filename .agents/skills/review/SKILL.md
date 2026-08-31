---
name: review-gophercloud
description: Review Gophercloud Go SDK code for API coverage gaps, convention violations, microversion mistakes, missing tests, and bugs. Use when the user invokes this skill or asks to review Gophercloud packages, PRs, diffs, or OpenStack API implementations.
---

# Review Gophercloud

Review Gophercloud (`github.com/gophercloud/gophercloud/v2`) as a core contributor would. Do not implement fixes unless the user asks.

## Patterns

Do not restate repository conventions here. Read [AGENTS.md](../../../AGENTS.md) and cite the matching section when a finding is a pattern violation, for example `AGENTS.md § 2.4` or `see § 5.3`.

Consult AGENTS.md for: running tests and lint (§ 1); imports, file layout, naming, and OptsBuilder (§ 2.1–2.3); pointers, `omitempty`, errors, and GoDoc (§ 2.4–2.6); unit and acceptance tests (§ 2.7); microversions (§ 3, § 6); pull requests (§ 4); `context.Context`, pagination, and the KeyedPage interface (§ 5).

For extra detail only: `docs/STYLEGUIDE.md`, `docs/MICROVERSIONS.md`, `.github/CONTRIBUTING.md`.

## Scope

Determine the target, then review only that:

- PR, branch, or diff: changed files vs `main` (committed + uncommitted)
- A package or service: that package tree, including `testing/` and matching acceptance tests
- New API / resource: completeness of `requests.go`, `results.go`, `urls.go`, `microversions.go`, unit tests, acceptance tests
- Unspecified: ask once whether to review the current diff or a named package

When the user gives a GitHub PR URL, fetch it with `gh pr view NNNN --json title,body,files,state,headRefName`, check it out with `gh pr checkout NNNN` (stash local changes first if needed), and read any linked issue with `gh issue view NNNN`. Note the OpenStack source links the PR provides as proof; if it claims a field/path/method and provides none, that is a suggestion to add one.

Always cover: correctness/bugs, Gophercloud conventions (AGENTS.md), microversions, tests, and API coverage gaps vs OpenStack.

## Workflow

1. Identify service, version, and resource (`openstack/<service>/<version>/<resource>/`).
2. Read `requests.go`, `results.go`, `urls.go`, `doc.go`, and `microversions.go` if present. A missing file is not automatically wrong; check sibling packages before flagging it (§ 2.2).
3. Read unit tests under `testing/` and acceptance tests under `internal/acceptance/openstack/<service>/`. Confirm tests exercise **every** request and response field (§ 2.7), not a happy-path subset.
4. Compare the code to AGENTS.md. For a convention miss, name the § in the finding. Before calling a pattern wrong, grep for how sibling packages do it.
5. For API completeness, compare Go types to **OpenStack source on a non-`master` branch** (e.g. `stable/2025.1`), not API docs alone. Link the proving source in findings.
6. Run the package unit tests when reviewing local changes (§ 1.1):

   ```bash
   go test ./openstack/<service>/<version>/<resource>/...
   ```

7. Report findings. Do not commit, push, or open a PR.

Call out missing `omitempty`, missing pointer response fields, or missing GoDoc version notes (§ 3). Do NOT suggest client-side validation of microversions — Gophercloud does not validate microversion compatibility (§ 6).

When you lack the context to judge something, say so instead of guessing: phrase it as a question or a "consider" suggestion rather than a blocker. A false-positive blocker costs more reviewer trust than an honest question.

## Output

Use this structure. Omit a section only when it has no items (do not invent filler). Put **file and line** on every finding that maps to code (`path:line`). Package-level or PR-description issues can omit a line. When the issue is a documented pattern, add `AGENTS.md § N`.

```markdown
## Summary
What was reviewed (diff, package, or PR) and the overall verdict in 2–4 sentences.

## Blockers
Must change before merge or before treating the API as complete.

- `openstack/.../requests.go:215` — one-sentence problem (`AGENTS.md § 2.4`).
  Include an OpenStack source link when claiming a field/path/method mismatch.

## Warnings
Should change; likely bugs, missing tests that would miss the bug, or incomplete coverage of claimed work.

- `openstack/.../results.go:88` — ... (`AGENTS.md § 5.3`)

## Suggestions
Optional improvements: style, GoDoc, extra tests, API gaps the PR did not claim to add, PR hygiene.

- `openstack/.../testing/requests_test.go:40` — ... (`AGENTS.md § 2.7`)

## Positives
What is already in good shape. Be specific; cite file and line, and § when it matches a pattern.
```

Classification:

- **Blockers**: wrong HTTP method/path, broken Extract, data loss, panic, auth/token mishandling, missing `omitempty` on a microversioned request field (§ 2.4, § 3), required field omitted, incorrect JSON names, wrong time type, list `Page` with multiple top-level JSON keys and no `ResourceKey()` when `AllPages` would pick the wrong collection (§ 5.3)
- **Warnings**: tests that cannot catch the bug (§ 2.7), incomplete field coverage for new/claimed APIs, missing acceptance test for a new operation, naming/OptsBuilder drift (§ 2.3), pointer/`omitempty` uncertainty (§ 2.4), list body with extra top-level keys but no `ResourceKey()` yet (§ 5.3)
- **Suggestions**: style, import order (§ 2.1), comment nits (§ 2.6), extra OpenStack operations/fields not in scope, PR description gaps (§ 4), missing links to OpenStack source code
- **Positives**: always include at least one concrete item when anything is correct

When you lack the context to classify something, raise it as a question in the closest-fitting section rather than forcing it into Blockers.

If there are no blockers, say the code is not blocked and keep warnings/suggestions/positives as applicable. Residual risk (no live API, untested microversions) goes under Warnings.

## Rules of the review

- Discover before deciding: grep for how sibling packages do it before flagging a pattern as wrong.
- Follow existing patterns (§ 2, § 5); do not invent new ones.
- Cite code precisely: `path:line` or `path:startLine-endLine`.
- Prefer OpenStack source over docs; link to service source on `stable/*` branches, not `master` or a commit SHA (§ 4).
- Test before claiming: run `go test ./...` for affected packages.
- When uncertain, ask: flag as a question or "consider" rather than a blocker.
- Respect the author: they may know something you do not. Ask when the rationale is not clear rather than asserting a fix.
- Do not squash-commit advice unless the user is writing the PR description (§ 4).
