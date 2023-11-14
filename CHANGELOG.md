# **unreleased**

## v1.1.0

* fix: add ReadHeaderTimeout to stats server
* fix(lint): typos, unused args, struct alignment
* fix: update to current syntax (release)
* chore: remove deprecated linters
* build(deps): bump github.com/nxadm/tail from 1.4.8 to 1.4.11
* build(deps): bump github.com/circonus-labs/circonus-gometrics/v3 from 3.4.6 to 3.4.7
* build(deps): bump golang.org/x/sys from 0.12.0 to 0.14.0
* build(deps): bump github.com/rs/zerolog from 1.29.1 to 1.31.0
* build(deps): bump golang.org/x/sync from 0.3.0 to 0.5.0
* build(deps): bump github.com/spf13/viper from 1.10.1 to 1.17.0
* build(deps): bump github.com/spf13/cobra from 1.3.0 to 1.8.0
* build(deps): bump github.com/pelletier/go-toml from 1.9.4 to 1.9.5
* fix: update to latest in lint workflow

## v1.0.3

* fix: typo in systemd service file

## v1.0.2

* fix: config unit tests
* fix: test tcp port for agent config
* fix: lowercase config err msgs

## v1.0.1

* fix: lint issues
* upd: remove pkg/errors dependency
* add: lint config
* upd: all dependencies

## v1.0.0

* upd: clarify processing metric vs sending log message
* upd: go1.15
* upd: cgm and viper dependencies
* fix: exec path in service config, add `sbin`

## v1.0.0-beta.3

* fix: let cgm handle retrying check/agent destinations rather than testing the port
* upd: version information in release package
* upd: skip README.md in `log.d`

## v1.0.0-beta.2

* fix: do not skip config if log file not currently available let tail wait for the log
* fix: do not exit if only one config with a missing log file
* upd: switch to nxadm/tail updated tail (hpcloud/tail is not currently being updated/maintained)
* upd: switch to tail poll for better handling of deleted log files

## v1.0.0-beta.1

* UPDATE - stream tag only release
* add: stream tag support
* upd: log id in tag `log_id`
* add: support `tags` attribute to metric stanza in configuration
* add: support interpolation of named regex patterns in `tags` attribute
* upd: dependencies (cgm)

## v0.6.0

* Merge pull request #12 from yargevad/agent-interval
* upd: dependencies, tidy
* upd: deprecated syntax
* add: agent interval to config, command line
* fix: lint issues
* add: build linting
* doc: agent interval

## v0.5.4

* fix: clarify setting for check, it is a check id not a check bundle id

## v0.5.3

* upd: dependencies (cgm - allow/deny checks default, etc.)

## v0.5.2

* upd: dependencies (yaml, x/sync)

## v0.5.1

* upd: dependencies
* upd: disable dir/file permissions tests
* add: ctx cancel function for Stop to use
* fix: calls to New in tests to include a parent context

## v0.5.0

* upd: include service and example configs in release
* add: systemd service configuration in `service/`
* add: example config `etc/example-circonus-logwatch.yaml`
* upd: change found metric message priority from Info to Debug
* upd: tomb->context+errgroup
* add: api ca file load/config for circonus metric destination
* upd: switch to circonus-gometrics v3
* upd: condense/consolidate code
* upd: switch to go mod

## v0.4.0

* upd: release file names use x86_64, facilitate automated builds and testing
* upd: goreleaser, turn off draft

## v0.3.1

* add: freebsd to release

## v0.3.0

* add: example of apache access log  latency histograms
* doc: include `etc/log.d` in release by adding a stub readme
* fix: correctly error when only one config with an error in it
* doc: formatting

## v0.2.0

* adjust handling of tail behavior for rotated logs and nil lines
* fix: statsd base metric name, add separator
* check destination, add api config

## v0.1.1

* update readme, clarify metric type descriptions
* common appstats package

## v0.1.0

* Initial development release
