# Speculator
---

Speculator is an automated scaling service for AWS EC2 AutoScale Groups based
on Kubernetes Pod CPU/Memory Request Percentages of the cluster.

Default Ideal range is lower than typical solutions (40%-57%) because of the way
Deis deploys on my cluster.  Range is modifyable via environment variables.

Twirp is a framework for service-to-service communication emphasizing simplicity
and minimalism. It generates routing and serialization from API definition files
and lets you focus on your application's logic instead of thinking about
folderol like HTTP methods and paths and JSON.

### Documentation

Thorough documentation is [on the wiki](https://github.com/undeadops/speculator/wiki).

### Releases
Twirp follows semantic versioning through git tags, and uses Github Releases for
release notes and upgrade guides:
[Twirp Releases](https://github.com/twitchtv/twirp/releases)

### Contributing
Check out [CONTRIBUTING.md](./CONTRIBUTING.md) for notes on making contributions.

### License

This library is licensed under the Apache 2.0 License. 
