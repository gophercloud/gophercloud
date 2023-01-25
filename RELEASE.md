# Gophercloud release

## Contributions

### The semver label

Gophercloud follows [semver](https://semver.org/).

Each Pull request must have a label indicating its impact on the API:
* `semver:patch` for changes that don't impact the API
* `semver:minor` for changes that impact the API in a backwards-compatible fashion
* `semver:major` for changes that introduce a breaking change in the API

Automation prevents merges if the label is not present.

### Metadata

The release notes for a given release are generated based on the PR title and its milestone:
* make sure that the PR title is descriptive
* add a milestone based on the semver label: x++ if major, y++ if minor, z++ if patch.

## Release of a new version

### Step 1: Check the metadata

Check that all pull requests merged since the last release have the right milestone.

### Step 2: Release notes and version string

Once all PRs have a sensible title and are added to the right milestone, generate the release notes with the [`gh`](https://github.com/cli/cli) tool:
```shell
gh pr list \
        --state merged \
        --search 'milestone:vx.y.z' \
        --json number,title \
        --template \
        '{{range .}}* {{ printf "[GH-%v](https://github.com/gophercloud/gophercloud/pull/%v)" .number .number }} {{ .title }}
{{end}}'
```

Replace `x.y.z` with the current milestone.

Add that to the top of `CHANGELOG.md`. Also add any information that could be useful to consumers willing to upgrade.

**Set the new version string in the `DefaultUserAgent` constant in `provider_client.go`.**

Create a PR with these two changes. The new PR should be labeled with the semver label corresponding to the type of bump, and the milestone corresponding to its version.

### Step 3: Git tag and Github release

The Go mod system relies on Git tags. In order to simulate a review mechanism, we rely on Github to create the tag through the Release mechanism.

* [Prepare a new release](https://github.com/gophercloud/gophercloud/releases/new)
* Let Github generate the  release notes by clicking on Generate release notes
* Click on **Save draft**
* Ask another Gophercloud maintainer to review and publish the release

_Note: never change a release or force-push a tag. Tags are almost immediately picked up by the Go proxy and changing the commit it points to will be detected as tampering._
