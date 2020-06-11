# v1.0.0-beta.1

* UPDATE - stream tag only release
* add: stream tag support
* upd: log id in tag `log_id`
* add: support `tags` attribute to metric stanza in configuration
* add: support interpolation of named regex patterns in `tags` attribute
* upd: depdencies (cgm)

# v0.6.0

* Merge pull request #12 from yargevad/agent-interval
* upd: depedencies, tidy
* upd: deprecated syntax
* add: agent interval to config, command line
* fix: lint issues
* add: build linting
* doc: agent interval

# v0.5.4

* fix: clarify setting for check, it is a check id not a check bundle id

# v0.5.3

* upd: dependencies (cgm - allow/deny checks default, etc.)

# v0.5.2

* upd: dependencies (yaml, x/sync)

# v0.5.1

* upd: dependencies
* upd: disable dir/file permissions tests
* add: ctx cancel function for Stop to use
* fix: calls to New in tests to include a parent context

# v0.5.0

* upd: include service and example configs in release
* add: systemd service configuration in `service/`
* add: example config `etc/example-circonus-logwatch.yaml`
* upd: change found metric message priority from Info to Debug
* upd: tomb->context+errgroup
* add: api ca file load/config for circonus metric destination
* upd: switch to circonus-gometrics v3
* upd: condense/consolidate code
* upd: switch to go mod

# v0.4.0

* upd: release file names use x86_64, facilitate automated builds and testing
* upd: goreleaser, turn off draft

# v0.3.1

* add: freebsd to release

# v0.3.0

* add: example of apache access log  latency histograms
* doc: include `etc/log.d` in release by adding a stub readme
* fix: correctly error when only one config with an error in it
* doc: formatting

# v0.2.0

* adjust handling of tail behavior for rotated logs and nil lines
* fix: statsd base metric name, add separator
* check destination, add api config

# v0.1.1

* update readme, clarify metric type descriptions
* common appstats package

# v0.1.0

* Initial development release
