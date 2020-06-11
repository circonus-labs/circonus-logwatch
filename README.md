# Circonus Log Watch

A small utility for extracting metrics from log files and forwarding to Circonus.

> Note: In v1 all metrics use stream tags, for non-stream tagged metrics use the v0 branch/releases.

## Install

1. `mkdir -p /opt/circonus/{sbin,etc,etc/log.d}`
1. Download [latest release](../../releases/latest) from repository (or [build manually](#manual-build))
1. Extract archive into `/opt/circonus`
1. Create config in `/opt/circonus/etc` and log configs in `/opt/circonus/etc/log.d`

## Options

```sh
/opt/circonus/sbin/circonus-logwatchd -h

Flags:
      --api-app string              [ENV: CLW_API_APP] Circonus API Token app (default "circonus-logwatch")
      --api-ca-file string          [ENV: CLW_API_CA_FILE] Circonus API CA certificate file (optional, Inside API server not using public certs)
      --api-key string              [ENV: CLW_API_KEY] Circonus API Token key or 'cosi' to use COSI config
      --api-url string              [ENV: CLW_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
  -c, --config string               config file (default is /opt/circonus/etc/circonus-logwatch.(json|toml|yaml)
  -d, --debug                       [ENV: CLW_DEBUG] Enable debug messages
      --debug-cgm                   [ENV: CLW_DEBUG_CGM] Enable CGM & API debug messages
      --debug-metric                [ENV: CLW_DEBUG_METRIC] Enable metric rule evaluation tracing debug messages
      --debug-tail                  [ENV: CLW_DEBUG_TAIL] Enable log tailing messages
      --dest string                 [ENV: CLW_DESTINATION] Destination[agent|check|log|statsd] type for metrics (default "log")
      --dest-agent-interval string  [ENV: CLW_DEST_AGENT_INTERVAL] Destination[agent] Interval for metric submission to agent (default "60s")
      --dest-cid string             [ENV: CLW_DEST_CID] Destination[check] Check ID (not check bundle)
      --dest-id string              [ENV: CLW_DEST_ID] Destination[statsd|agent] metric group ID (default "circonus-logwatch")
      --dest-instance-id string     [ENV: CLW_DEST_INSTANCE_ID] Destination[check] Check Instance ID
      --dest-port string            [ENV: CLW_DEST_PORT] Destination[agent|statsd] port (agent=2609, statsd=8125)
      --dest-statsd-prefix string   [ENV: CLW_DEST_STATSD_PREFIX] Destination[statsd] Prefix prepended to every metric sent to StatsD (default "host.")
      --dest-tag string             [ENV: CLW_DEST_TAG] Destination[check] Check search tag
      --dest-target string          [ENV: CLW_DEST_TARGET] Destination[check] Check target (default hostname)
      --dest-url string             [ENV: CLW_DEST_URL] Destination[check] Check Submission URL
  -h, --help                        help for circonus-logwatch
  -l, --log-conf-dir string         [ENV: CLW_PLUGIN_DIR] Log configuration directory (default "/opt/circonus/etc/log.d")
      --log-level string            [ENV: CLW_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty                  [ENV: CLW_LOG_PRETTY] Output formatted/colored log lines
      --show-config                 Show config (json|toml|yaml) and exit
      --stat-port string            [ENV: CLW_STAT_PORT] Exposes app stats while running (default "33284")
  -V, --version                     Show version and exit

```

## Destinations

* `--dest check` metrics are sent directly to the circonus broker (will create a check if `--dest-cid` not provided). `--dest-instance-id`, `--dest-target`, and `--dest-tag` can be used to customize the check created.
* `--dest agent` metrics are sent to `/write` endpoint of local circonus-agent (`http://localhost:2609/write/id`) uses `--dest-id` to categorize the metrics. `--dest-port` controls the agent port (default 2609)
* `--dest statsd` metrics sent to statsd listener of local circonus-agent (`localhost:8125`) uses `--statsd-prefix` for each metric name, followed by `--dest-id` (`--dest-statsd-prefix` should match circonus-agent `--statsd-host-prefix` to ensure metrics are routed to correct destination by the agent). `--dest-port` controls the agent statsd port (default 8125)

## Config

Create a JSON, YAML, or TOML config in `/opt/circonus/etc/circonus-logwatch.(json|yaml|toml)`. Or, use environment variables and/or command line parameters.

YAML with a StatsD destination (send metrics to local circonus-agent statsd listener)

```yaml
---
app_stat_port: "33284"
debug: false
debug_cgm: false
debug_metric: false
debug_tail: false
log_conf_dir: /opt/circonus/etc/log.d
api:
  key: circonus api token key
  app: circonus-logwatch
  url: 'https://api.circonus.com/v2/'
  ca_file: path to api ca file, if needed
destination:
  type: statsd
  config:
    id: circonus-logwatch
    port: "8125"
    statsd_prefix: "host."
```

JSON with a Check destination (send metrics directly to a Circonus check, a check will be created if not supplied via cid or url config settings)

> NOTE: use `cid` or `url` to identify an existing check. Or, to search for a check a combination of `target`, `search_tag` and/or `instance_id`.

```json
{
    "api": {
        "app": "circonus-logwatch",
        "ca_file": "path to api ca file, if needed",
        "key": "circonus api token key",
        "url": "https://api.circonus.com/v2/"
    },
    "app_stat_port": "33284",
    "debug": false,
    "debug_cgm": false,
    "debug_metric": false,
    "debug_tail": false,
    "destination": {
        "type": "check",
        "config": {
            "cid": "check id (of existing check)",
            "url": "submission url of existing check",
            "target": "to find|create a check",
            "search_tag": "to find|create a check",
            "instance_id": "to find|create a check"
        }
    },
    "log": {
        "level": "info",
        "pretty": false
    },
    "log_conf_dir": "/opt/circonus/etc/log.d"
}
```

TOML with an Agent destination (send metrics to local circonus-agent)

```toml
app_stat_port = "33284"
debug = false
debug_cgm = false
debug_metric = false
debug_tail = false
log_conf_dir = "/opt/circonus/etc/log.d"

[api]
  key = "circonus api token key"
  app = "circonus-logwatch"
  url = "https://api.circonus.com/v2/"
  ca_file = "path to api ca file, if needed"

[destination]
  type = "agent"

  [destination.config]
    id = "circonus-logwatch"
    port = "2609"

[log]
  level = "info"
  pretty = false
```

## Log configs

Create one config (JSON, YAML, or TOML) in `--log-conf-dir` for each distinct log. Examples [`etc/log.d`](etc/log.d/) in this repository

1. `id` of the log, short identifier - optional, the base file name will be used if omitted
1. `log_file` path to the log
1. `metrics` a list of:
    1. `match` regular expression to identify lines and optionally extract named subexpressions for value and metric name
    1. `name` a static string to use as the metric name or a template for naming the metric if named subexpressions were used in match regex
    1. `tags` comma separated list of k:v pairs, templating can be used accessing named subexpressions (e.g. `foo:bar,yabba:dabba` or `foo:{{.id}},bar:baz`)
    1. `type` what type of metric (all numbers are 64bit)
        * `c` counter int
        * `g` gauge int or float
        * `ms` timing (treated as histogram) value can be a float (3.12) or duration (25ms, 3.2s, 1m10s, etc.) durations are converted to milliseconds, floats are assumed to already represent milliseconds
        * `h` histogram float
        * `s` set (ala statsd set metrics) unique string to count
        * `t` text string

### Log configuration notes

* any metric which does not have a `type` will be treated as a counter.
* any metric which does not have a subexpression named '*Value*' (case insensitive) will be treated as a counter.
* named subexpressions can be used in the name template and tag list with the following syntax `{{.id}}` where `id` is the name given to a named subexpression in the match regex
* metrics will have a stream tag added for the log `id` (e.g. for a log with an id of "foo" the tag would be `log_id:foo`)

## Manual build

1. Clone repo (outside if `GOPATH`)`git clone https://github.com/circonus-labs/circonus-logwatch && cd circonus-logwatch`
1. Build `go build -o circonus-logwatchd`
1. Ensure target directories exist `mkdir -p /opt/circonus/logwatch/{sbin,etc,etc/log.d}`
1. Install `cp circonus-logwatchd /opt/circonus/logwatch/sbin`
