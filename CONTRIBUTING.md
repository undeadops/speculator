# Contributing #

First, if you have run into a bug, please file an issue. I will try to get back to
issue reporters within a day or two.

If you'd rather not publicly discuss the issue, please email mitch@undeadops.xyz.

Issues are also a good place to present experience reports or requests for new
features.

## Setup Requirements ##

You will need git and I haven't tested on anything lower than Go 1.8

I work off feature branches, and merge to master.  Tag commits in master for
releases. 

## Developer Loop ##

Application Tests are needed, if you wanted to help that would be helpful :)

## Contributing Code ##

Github pull requests will be accpeted. Fork from `master`, and submit PR.

## Releasing Versions ##

Releasing versions is the responsibility of the core maintainers. Most people
don't need to know this stuff.

Speculator uses [Semantic versioning](http://semver.org/): `v<major>.<minor>.<patch>`.

 * Increment major if you're making a backwards-incompatible change.
 * Increment minor if you're adding a feature that's backwards-compatible.
 * Increment patch if you're making a bugfix.

Speculator uses Github releases. To make a new release:
 1. Merge all changes that should be included in the release into the master
    branch.
 2. Add a new commit to master with a message like "Version vX.X.X release".
 3. Tag the commit you just made: `git tag <version number>` and `git push
    origin --tags`

