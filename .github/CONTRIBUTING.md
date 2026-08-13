# Contributing to Gophercloud

- [New Contributor Tutorial](#new-contributor-tutorial)
- [Ways to get involved](#ways-to-get-involved)
- [Getting started](#getting-started)
- [Tests](#tests)
- [Style guide](#style-guide)

## New Contributor Tutorial

For new contributors, we've put together a detailed tutorial
[here](https://github.com/gophercloud/gophercloud/tree/main/docs/contributor-tutorial)!

## Ways to get involved

There are several ways you can get involved in our open-source project, each
described briefly below.

### 1. Fixing bugs

If you want to start fixing open bugs, we'd really appreciate that! Bug fixing
is central to any project. The best way to get started is by heading to our
[bug tracker](https://github.com/gophercloud/gophercloud/issues) and finding open
bugs that you think nobody is working on. It might be useful to comment on the
thread to see the current state of the issue and if anybody has made any
breakthroughs on it so far.

New to the project? Look for issues labelled [`good first issue`](https://github.com/gophercloud/gophercloud/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
to find a good starting point.

### 2. Improving documentation

Gophercloud's documentation is automatically generated from the source code
and can be read online at [godoc.org](https://godoc.org/github.com/gophercloud/gophercloud).

If you feel that a certain section could be improved - whether it's to clarify
ambiguity, correct a technical mistake, or to fix a grammatical error - please
feel entitled to do so! We welcome doc pull requests with the same childlike
enthusiasm as any other contribution!

### 3. Working on a new feature

If you've found something we've left out, we'd love for you to add it! Please
first open an issue to indicate your interest to a core contributor - this
enables quick/early feedback and can help steer you in the right direction by
avoiding known issues. It might also help you avoid losing time implementing
something that might not ever work or is outside the scope of the project.

While you're implementing the feature, one tip is to prefix your Pull Request
title with `[wip]` - then people know it's a work in progress. Once the PR is
ready for review, you can remove the `[wip]` tag and request a review.

We ask that you do not submit a feature that you have not spent time researching
and testing first-hand in an actual OpenStack environment. While we appreciate
the contribution, submitting code which you are unfamiliar with is a risk to the
users who will ultimately use it. See our [acceptance tests readme](/internal/acceptance/README.md)
for information about how you can create a local development environment to
better understand the feature you're working on.

### 4. Reviewing pull requests

Reviewing [open pull requests](https://github.com/gophercloud/gophercloud/pulls)
is a valuable contribution in its own right. Fresh eyes catch bugs, improve
code quality, and help contributors learn from each other. Feel free to leave
comments, ask clarifying questions, or approve changes once you're satisfied
with them.

Please do not hesitate to ask questions or request clarification. Your
contribution is very much appreciated and we are happy to work with you to get
it merged.

## Getting Started

As a contributor you will need to setup your workspace. Here are the basic
instructions:

1. Fork the `gophercloud/gophercloud` repository on GitHub, then clone your
fork:

   ```bash
   git clone git@github.com:<my_username>/gophercloud.git
   cd gophercloud
   ```

2. Add the upstream repository as a remote so you can keep your fork up to
date:

   ```bash
   git remote add upstream https://github.com/gophercloud/gophercloud.git
   ```

3. Checkout the latest development branch:

   ```bash
   git checkout main
   git pull upstream main
   ```

4. If you're working on something (discussed more in detail above), you will
need to checkout a new feature branch:

   ```bash
   git checkout -b my-new-feature
   ```

5. Use a standard text editor or IDE of your choice to make your changes to the code or documentation. Once finished, commit them.

   ```bash
   git status
   git add path/to/changed/file.go
   git commit
   ```

6. Run checks locally before opening or updating your Pull Request:

   ```bash
   make unit     # run unit tests
   make lint     # run golangci-lint (requires docker or podman)
   make format   # run gofmt -s
   ```

   If your change should also be validated against a real OpenStack cloud,
   see our [acceptance tests readme](/internal/acceptance/README.md).

7. Before opening the PR, make sure a [GitHub issue](https://github.com/gophercloud/gophercloud/issues)
exists for non-trivial changes and has been acknowledged by a core
contributor. Submit your branch as a [Pull Request](https://help.github.com/articles/creating-a-pull-request/),
following our [Style Guide](/docs/STYLEGUIDE.md). Your PR must also carry one
of the `semver:patch`, `semver:minor`, or `semver:major` labels described in
the [release documentation](/RELEASE.md#the-semver-label). CI will block the
merge until that label is set.

> Further information about using Git can be found [here](https://git-scm.com/book/en/v2).

Happy Hacking!

## Tests

When working on a new or existing feature, testing will be the backbone of your
work since it helps uncover and prevent regressions in the codebase. There are
two types of test we use in Gophercloud: unit tests and acceptance tests, which
are both described below.

### Unit tests

Unit tests are the fine-grained tests that establish and ensure the behavior
of individual units of functionality. We usually test on an
operation-by-operation basis (an operation typically being an API action) with
the use of mocking to set up explicit expectations. Each operation will set up
its HTTP response expectation, and then test how the system responds when fed
this controlled, pre-determined input.

To make life easier, we've introduced a bunch of test helpers to simplify the
process of testing expectations with assertions:

```go
import (
  "testing"

  "github.com/gophercloud/gophercloud/testhelper"
)

func TestSomething(t *testing.T) {
  result, err := Operation()

  testhelper.AssertEquals(t, "foo", result.Bar)
  testhelper.AssertNoErr(t, err)
}

func TestSomethingElse(t *testing.T) {
  testhelper.CheckEquals(t, "expected", "actual")
}
```

`AssertEquals` and `AssertNoErr` will throw a fatal error if a value does not
match an expected value or if an error has been declared, respectively. You can
also use `CheckEquals` and `CheckNoErr` for the same purpose; the only difference
being that `t.Errorf` is raised rather than `t.Fatalf`.

Here is a truncated example of mocked HTTP responses:

```go
import (
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

func TestGet(t *testing.T) {
	// Setup the HTTP request multiplexer and server
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/networks/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		// Test we're using the correct HTTP method
		th.TestMethod(t, r, "GET")

		// Test we're setting the auth token
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		// Set the appropriate headers for our mocked response
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Set the HTTP body
		fmt.Fprintf(w, `
{
    "network": {
        "status": "ACTIVE",
        "name": "private-network",
        "admin_state_up": true,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "shared": true,
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22"
    }
}
			`)
	})

	// Call our API operation
	network, err := networks.Get(fake.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()

	// Assert no errors and equality
	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Status, "ACTIVE")
}
```

### Acceptance tests

As we've already mentioned, unit tests have a very narrow and confined focus -
they test small units of behavior. Acceptance tests on the other hand have a
far larger scope: they are fully functional tests that test the entire API of a
service in one fell swoop. They don't care about unit isolation or mocking
expectations, they instead do a full run-through and consequently test how the
entire system _integrates_ together. When an API satisfies expectations, it
proves by default that the requirements for a contract have been met.

Please be aware that acceptance tests will hit a live API - and may incur
service charges from your provider. Although most tests handle their own
teardown procedures, it is always worth manually checking that resources are
deleted after the test suite finishes.

We provide detailed information about how to set up local acceptance test
environments in our [acceptance tests readme](/internal/acceptance/README.md).

### Running tests

To run all tests:

  ```bash
  go test -tags fixtures ./...
  ```

To run all tests with verbose output:

  ```bash
  go test -v -tags fixtures ./...
  ```

To run tests that match certain [build tags]():

  ```bash
  go test -tags "fixtures foo bar" ./...
  ```

To run tests for a particular sub-package:

  ```bash
  cd ./path/to/package && go test -tags fixtures ./...
  ```

## Style guide

See [here](/docs/STYLEGUIDE.md)
