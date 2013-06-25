# gophercloud &mdash; V0.0.0
The Go ecosystem seems to lack a comprehensive cloud services API (at the time this README was first written). As both Go and cloud services are trending in many businesses, and with Go used increasingly in infrastructure, it seems like an odd omission. To fill this gap, gorax provides a Go binding to the Rackspace cloud APIs. Rackspace offers many APIs that are compatible with OpenStack, and thus provides an ideal springboard for wider OpenStack technology adoption in the Go community.

**This library is still in the very early stages of development. Unless you want to contribute, it probably isn't what you want.**

## How to Contribute
### Familiarize Yourself with Go
To contribute to gophercloud, you'll need some passing familiarity with Go, and how it uses certain concepts.  If you've never worked with Go before, I *strongly* encourage the interested reader to [install the latest version of Go](http://golang.org/), and follow through the excellent online book titled [Effective Go](http://golang.org/doc/effective_go.html).[^1]

[^1]: A common problem found often in newcomers to the Go community is preconceptions instilled through years of working with other languages or environments.  If you find something about Go or its workflow that irks you to the point of kvetching, I humbly ask that you be patient with yourself, keep any discussion of the matter informed and civil, and please do not use gophercloud's issue-tracking system as a soapbox unless it *really* impacts the usability of the gophercloud package by other Go developers.  Be respectful of other Go users who might disagree with you.

### Installing gophercloud in a Workspace

<span style="font-size: 300%; font-weight: bold; text-align:center; color:red">STOP!</span>

**Please** do not just clone this repository expecting it to work like any other Python, Ruby, Java, or C/C++ repo.  Go packages don't work that way!  (You *did* read Effective Go, right?)

#### Temporary: gophercloud is a Private Repository!  How do I install?

***Please note that these instructions are here temporarily until we make this repository public.***

1. Navigate to the scripts/create-environment.sh file in the GitHub user interface
2. Click on 'Raw'
3. Copy and paste into a file locally
4. Supposing that file is named `/tmp/gcsetup.sh`, you can invoke it as follows:

	`source /tmp/gcsetup.sh`

#### Installing for POSIX (Linux, BSD, MacOS X, et. al.)

***Please note that this set of installation instructions will work as soon the repoistory is public***

You can execute the following command to create a brand new Go workspace that is minimally configured for use with gophercloud.  This should work for any reasonable POSIX-compatible environment.

	source <(curl "https://raw.github.com/rackspace/gophercloud/master/scripts/create-env.sh")


### Contributing Features or Bug-Fixes
After installing gophercloud and after running `env.sh` (only needed once per shell session), you will find the source files in the `$GOPHERCLOUD` directory.  Feel free to open your favorite editor inside that directory and poke around.

Features and bug-fixes **must** appear on their own *feature branches*.  The name of the branch should be fairly descriptive, but try to avoid verbosity for verbosity's sake.  Examples of good feature branch names include:

* script-environment-setup
* server-creation
* issue-43-memory-leak-fix

Some examples of not-so-good branch names include:

* cloud-server-api-server-creation-endpoint-support
* tk
* anything/with/slashes

The former is lengthy without delivering much value, the second is too short to be useful to anyone other than the submitter, and the last tries to work around Git usability issues, where some commands separate origins from branch names using slashes, and thus could be considered ambiguous to a human operator.

For example, if you're looking to fix a memory leak that is documented in, just to pick a number, issue 42, you might follow a sequence of commands such as the following:

	cd $GOPHERCLOUD
	git checkout working
	git checkout -b issue-42-fix-memory-leak
	# edit edit edit ...
	# commits happen here ...
	git push -u origin issue-42-fix-memory-leak

At this point, you may now switch to the GitHub user interface, and open a pull-request for the feature branch.  This pull request should be descriptive.  Basically, you want to give a small code walkthrough in the pull request summary.  You should be able to answer, at a minimum, four basic questions, as appropriate for the nature of the patch:

1.  What is the problem?
2.  Why is it a problem?
3.  What is your solution?
4.  How does your solution actually work?

Here's a made-up example:

	Fix memory leak detailed in issue #42.
	
	The Rackspace provider interface tended to leak memory every fifth
	Saturday of February.  Over the course of several decades, we find
	we run out of memory.  Killing and restarting the process periodically
	restores service, but is a burden on the ops team.  This PR fixes this
	bug permanently.
	
	The barProvider structure found in
	provider/barisp.go defines a FooSet as a slice, as seen on line 314.
	Per services/auth/keystone2.go line 628, Keystone authentication
	only ever uses the first three	elements of this FooSet.  Line 42 shows
	where FooSet is initialized to an empty slice, but on line 512, we see
	a function that appends to this slice unconditionally.
	
	I'm not sure where the logic exists to determine where this function is
	called; so, I've adjusted the provider/barisp.go file to truncate this
	FooSet to only three items, maximum on behalf of the caller.  This seems
	to solve the problem in my test cases.  See included tests.

Obviously, please use common sense!  In situations where these questions do not apply, please don't make up filler information.

|NOTES|
|-----|
|All bug-fix PRs **MUST** reference at least one open issue.  New feature PRs **SHOULD** reference at least one open issue.  This convention helps track *why* certain code is written the way it is, and maintains historical context.|
|You may elide answers to the questions above if the answers already appear in the referenced PR(s).  We care that the answers exist and may be easily found, not so much about *where* the answers may be found.|

### Master versus Working Branch

Many projects will happily let you create a feature branch off the master branch.  However, Go environments place special significance on master branches of packages.  Because the `go get` command *is not* intended to perform complete package management tasks, but merely serve as a convenience for establishing your Go work environment, it will always fetch from the master branch of any repository you specify.  **Therefore, the master branch MUST always represent a customer-installable package.**  Not only that, but interface changes **must** be backward compatible at all times.

To facilitate development efforts, then, we maintain a *working* branch.  New features and bug fixes merge into the `working` branch, where it remains staged for some future release date.  Ideally, every push to github and every merge to `working` should kick off a batch of tests to validate the product still works.  Assuming that `working` tests all pass, *and* your features or bug-fixes are both code- and feature-complete, then and only then should `working` be merged into `master`.
