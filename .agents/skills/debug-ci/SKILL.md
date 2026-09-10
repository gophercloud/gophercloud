---
name: debug-ci
description: "Gophercloud CI: Debug acceptance test failures in GitHub Actions. Download artifacts, analyze Go test output, trace through OpenStack devstack service logs to find root cause."
---

# Gophercloud CI Failure Debugging

Debug failures in the Gophercloud acceptance test suite run via GitHub Actions.

## Prerequisites

- `gh` CLI authenticated with access to `gophercloud/gophercloud`
- Working directory: the gophercloud repository root

## Overview

Gophercloud acceptance tests run Go tests (via `gotestsum --format testname`) against real OpenStack environments deployed with [devstack](https://github.com/gophercloud/devstack-action) in GitHub Actions. Each functional workflow runs against a matrix of OpenStack versions:

| Name | Branch |
|---|---|
| master | `master` |
| gazpacho | `stable/2026.1` |
| epoxy | `stable/2025.1` |

> This table reflects the current CI matrix. Verify against the workflow files under `.github/workflows/functional-*.yaml` if it may be stale. Not all workflows run on every version. For example, `functional-metric` currently only runs on `master`.

When a test fails, the CI uploads devstack logs as artifacts (only on failure).

## Step 1: Identify the Failure

Given a CI run URL, PR number, or run ID, identify which job(s) failed and which test(s) failed:

```bash
# List jobs and their status
gh run view <RUN_ID> --repo gophercloud/gophercloud \
  --json jobs --jq '.jobs[] | select(.conclusion == "failure") | {name: .name, conclusion: .conclusion, id: .databaseId}'

# Get logs for a specific failed job and find the failing test(s)
gh run view <RUN_ID> --repo gophercloud/gophercloud \
  --log --job <JOB_ID> 2>&1 | grep -E "FAIL\b"

# Get the full test output around a failure (error message follows the FAIL line)
gh run view <RUN_ID> --repo gophercloud/gophercloud \
  --log --job <JOB_ID> 2>&1 | grep -A5 "FAIL.*Test"

# If the run URL is from a PR, find runs for that PR
gh run list --repo gophercloud/gophercloud --branch <BRANCH> --limit 5 \
  --json databaseId,status,conclusion,name
```

### Reading gotestsum output

Gophercloud uses `gotestsum --format testname`. In a local terminal, failed tests are marked with `✖`, but in `gh run view --log` output the markers are replaced by ANSI-escaped `FAIL` wrapped in GitHub Actions `##[group]` annotations. Key patterns in the CI log output:

- `##[group] FAIL internal/acceptance/openstack/db/v1.TestFooBar (0.07s)` — test failed (start of group)
- `    foo_test.go:43: <error message>` — the actual error, printed on the line(s) **after** the FAIL line (inside the group block)
- `##[endgroup]` — end of the failed test's output block
- `FAIL Package internal/acceptance/openstack/...` — package-level failure summary

### Finding the test source

Acceptance tests live under `internal/acceptance/openstack/<service>/`. Given a failing test name like `TestServersCreateDestroy`, locate it:

```bash
grep -r "func TestServersCreateDestroy" internal/acceptance/
```

Read the test to understand what OpenStack operations it performs and what assertions it makes.

## Step 2: Download Artifacts

Artifacts are uploaded only when tests fail. The naming pattern is:

```
functional-<service>-<openstack_name>-<run_id>
```

For example: `functional-compute-epoxy-12345678901`

```bash
# List available artifacts
gh api repos/gophercloud/gophercloud/actions/runs/<RUN_ID>/artifacts \
  --jq '.artifacts[] | {name: .name, size_in_bytes: .size_in_bytes}'

# Download artifacts for a specific job
gh run download <RUN_ID> --repo gophercloud/gophercloud \
  --name functional-<service>-<openstack_name>-<RUN_ID> --dir /tmp/ci-artifacts
```

### Artifact Contents

The `script/collectlogs` script collects the following on failure:

| File | Description |
|---|---|
| `journal.log` | Full system journal (`journalctl`) from the runner |
| `devstack-services.txt` | Status of all `devstack@*` systemd services |
| `<service>.log` | Individual devstack service log (one per running service, see mapping below) |
| `free.txt` | Memory usage at time of log collection |
| `dpkg-l.txt` | Installed Debian packages |
| `pip-freeze.txt` | Python package versions (useful for identifying OpenStack service versions) |
| `local.conf` | Devstack configuration used for deployment |

## Step 3: Analyze the Failure

### 3a. Read the Go Test Output

1. Find the failing test function name from the `FAIL` line
2. Read the error message — Gophercloud tests typically print the HTTP response body on failure, which includes the OpenStack error message and request ID
3. Locate the test source under `internal/acceptance/openstack/<service>/`
4. Read the test to understand what it does and what went wrong

### 3b. Trace OpenStack API Errors

Gophercloud's test errors usually include:

- The HTTP status code (e.g., `404 Not Found`, `409 Conflict`)
- The response body with the OpenStack error message
- The OpenStack request ID (in the `X-Openstack-Request-Id` header or response body)

To trace server-side:

```bash
# Find the OpenStack-side processing for a specific request
grep "req-<REQUEST_ID>" /tmp/ci-artifacts/<service>.log

# Find all errors in a service log
grep -i "error\|traceback\|exception" /tmp/ci-artifacts/<service>.log

# Find HTTP requests by status code
grep "HTTP/1.1" /tmp/ci-artifacts/<service>.log | grep " 4[0-9][0-9]\| 5[0-9][0-9]"
```

### 3c. Check for Devstack Infrastructure Issues

```bash
# Check if any devstack services crashed
grep -i "failed\|inactive\|dead" /tmp/ci-artifacts/devstack-services.txt

# Check for OOM conditions
cat /tmp/ci-artifacts/free.txt

# Check for system-level errors (OOM killer, disk full)
grep -i "oom\|out of memory\|no space left" /tmp/ci-artifacts/journal.log

# Check the devstack configuration
cat /tmp/ci-artifacts/local.conf
```

### 3d. Determine Test Timing

Go test output includes timestamps. To correlate with OpenStack logs:

1. Find the test start/end time from the Go test output
2. Search OpenStack service logs within that time window

```bash
# Find log entries in a time window (OpenStack logs use HH:MM:SS format)
grep "14:38\|14:39\|14:40" /tmp/ci-artifacts/<service>.log | grep -i "error"
```

## Devstack Service Log Mapping

### Core services (always present in default devstack)

| OpenStack Service | Log File | Description |
|---|---|---|
| Keystone | `keystone.log` | Identity / authentication |
| Nova API | `n-api.log` | Compute API |
| Nova Compute | `n-cpu.log` | Compute hypervisor |
| Nova Scheduler | `n-sch.log` | Compute scheduling |
| Nova Conductor | `n-cond.log` | Compute conductor |
| Glance API | `g-api.log` | Image service |
| Cinder API | `c-api.log` | Block storage API |
| Cinder Volume | `c-vol.log` | Block storage volume manager |
| Cinder Scheduler | `c-sch.log` | Block storage scheduling |
| Cinder Backup | `c-bak.log` | Block storage backup |
| Neutron | `q-svc.log` | Network API (present in all networking configurations) |
| Placement | `placement-api.log` | Placement API |

> **Note:** The `functional-networking` workflow uses OVN (the devstack default) and does **not** run Neutron agents (`q-l3`, `q-dhcp`, `q-meta`, `q-agt`). These agents are only present in workflows that explicitly set `Q_AGENT=openvswitch` (currently `functional-baremetal` and `functional-fwaas_v2`). Per-workflow log files are listed in the table below.

## Workflow Reference

| Workflow | Test Path | Key Logs |
|---|---|---|
| functional-basic | `internal/acceptance/openstack` | `keystone.log`, `g-api.log`, `s-proxy.log` |
| functional-baremetal | `internal/acceptance/openstack/baremetal/...` | `ir-api.log`, `ir-cond.log`, `n-api.log` |
| functional-blockstorage | `internal/acceptance/openstack/blockstorage/...` | `c-api.log`, `c-vol.log`, `c-sch.log` |
| functional-compute | `internal/acceptance/openstack/compute/...` | `n-api.log`, `n-cpu.log`, `n-sch.log` |
| functional-container | `internal/acceptance/openstack/container/...` | `zun-api.log`, `zun-compute.log` |
| functional-containerinfra | `internal/acceptance/openstack/containerinfra/...` | `magnum-*.log`, `h-api.log` (epoxy only), `h-eng.log` (epoxy only) |
| functional-db | `internal/acceptance/openstack/db/...` | `tr-api.log`, `tr-cond.log`, `tr-tmgr.log` |
| functional-dns | `internal/acceptance/openstack/dns/...` | `designate-api.log`, `designate-central.log` |
| functional-fwaas_v2 | `internal/acceptance/openstack/networking/v2/extensions/fwaas_v2/...` | `q-svc.log`, `q-l3.log`, `q-agt.log` |
| functional-identity | `internal/acceptance/openstack/identity/...` | `keystone.log` |
| functional-image | `internal/acceptance/openstack/image/...` | `g-api.log` |
| functional-keymanager | `internal/acceptance/openstack/keymanager/...` | `barbican-svc.log` |
| functional-loadbalancer | `internal/acceptance/openstack/loadbalancer/...` | `o-api.log`, `o-cw.log`, `o-hm.log` |
| functional-messaging | `internal/acceptance/openstack/messaging/...` | `zaqar-*.log` |
| functional-metric | `internal/acceptance/openstack/metric/...` | `aetos.log`, `ceilometer-*.log` |
| functional-networking | `internal/acceptance/openstack/networking/...` | `q-svc.log` |
| functional-objectstorage | `internal/acceptance/openstack/objectstorage/...` | `s-proxy.log`, `s-account.log` |
| functional-orchestration | `internal/acceptance/openstack/orchestration/...` | `h-api.log`, `h-eng.log` |
| functional-placement | `internal/acceptance/openstack/placement/...` | `placement-api.log` |
| functional-sharedfilesystems | `internal/acceptance/openstack/sharedfilesystems/...` | `manila-*.log` |
| functional-workflow | `internal/acceptance/openstack/workflow/...` | `mistral-api.log`, `mistral-engine.log` |

## Step 4: Common Failure Patterns

| Category | Symptoms | Likely Cause |
|---|---|---|
| Devstack deploy failure | No test output, "Deploy devstack" step failed | OpenStack upstream / devstack bug; check `journal.log` |
| Service crash | Tests fail mid-run with `503` or connection refused | Devstack service crashed; check `devstack-services.txt` and `journal.log` |
| Timeout | `context deadline exceeded` or 60m test timeout | Resource contention, slow service, or leaked resources from prior test |
| 404 Not Found | Resource created then immediately not found | Race condition, wrong endpoint, or version-specific API change |
| 409 Conflict | Resource already exists or state conflict | Test cleanup issue or concurrent test interference |
| 401 Unauthorized | Authentication or token error | Token expiration (long-running tests), service misconfiguration |
| 400 Bad Request | Invalid request body | Gophercloud sending wrong fields, or microversion mismatch |
| Microversion not supported | `Version X.Y is not supported by the API` | Test requires a microversion not available on that OpenStack release |
| Version-specific failure | Fails only on one OpenStack version (e.g., master but not epoxy) | API behavior change between releases |
| OOM | Tests killed, truncated output | Check `free.txt` and `journal.log` for OOM killer |
| Disk full | Various I/O errors | Check `journal.log` for "No space left on device" |

## Step 5: Classify the Failure

| Category | Description | Action |
|---|---|---|
| **OpenStack bug** | Devstack service returned unexpected or incorrect behavior | Document; may need upstream report to OpenStack |
| **Gophercloud bug** | SDK sends wrong request, mishandles response, or incorrect URL construction | File or fix in gophercloud |
| **Test bug** | Wrong assertions, missing setup/cleanup, incorrect environment assumptions | Fix the test |
| **Infrastructure flake** | Timeout, OOM, network issue, devstack instability | Document; consider retry or resource tuning |

### Classification Guardrails

Before classifying a failure as **Infrastructure flake**, verify:

1. You have **positive evidence** of an infrastructure problem (OOM in `free.txt`, service crash in `devstack-services.txt`, network timeout in `journal.log`) — not just "I can't find another explanation."
2. You have explicitly ruled out **Gophercloud bug** (reviewed the SDK code path), **Test bug** (verified assertions and test setup are correct), and **OpenStack bug** (checked service logs for API-level errors) with specific reasoning for why each doesn't apply.
3. State the specific evidence for your classification — if you cannot point to a concrete log line or metric, reconsider your conclusion.

"Unknown" is a valid classification. Do not default to "Infrastructure flake" when uncertain.

## Caveats

### Multiple Failures May Be Independent

When multiple tests fail in the same run, they may not be related — each acceptance test is designed to be self-contained. Analyze each failure independently before assuming a common cause.

### File-Filter May Skip Tests

Each workflow uses `.github/actions/file-filter` to skip tests when only unrelated files changed. If a workflow shows "No relevant files changed - skipping tests", this is not a failure — the job exits successfully. The filter patterns are defined per-workflow in the `patterns` field.

### Artifacts Only on Failure

Log artifacts are uploaded only when the test step fails (`if: failure()`). Successful runs have no artifacts. If you need logs from a successful run, they are only available in the job's stdout via `gh run view --log`.

### Scheduled Runs

Functional workflows also run on a schedule (`cron: '0 0 */3 * *'` — every 3 days). Scheduled run failures that don't reproduce on PRs often indicate flaky tests or transient devstack issues. Check the run's `event` field:

```bash
gh run view <RUN_ID> --repo gophercloud/gophercloud --json event --jq '.event'
```

### Falling Back to journal.log for Missing Service Logs

If a service's individual log file is unexpectedly absent from the artifacts, its output can still be found in `journal.log`:

```bash
grep "devstack@<service>" /tmp/ci-artifacts/journal.log
```

### OpenStack Version Differences

Failures on only one OpenStack version (e.g., `master` but not `epoxy`) may indicate:
- A new API behavior or deprecation in that release
- A devstack configuration difference
- A version-specific bug

Check `pip-freeze.txt` in the artifacts to identify exact service versions.
