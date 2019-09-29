# Sensu Go Stale Silence Plugin
[![TravisCI Build Status](https://api.travis-ci.com/rgeniesse/sensu-stale-silence-check.svg?branch=master)](https://travis-ci.com/rgeniesse/sensu-stale-silence-check)
[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/rgeniesse/sensu-stale-silence-check)

A Sensu Go check plugin for finding stale silence entries

## Installation

Download the latest version of the sensu-stale-silence-check from [releases][1],
or create an executable script from this source.

From the local path of the sensu-stale-silence-check repository:

```
go build -o /usr/local/bin/sensu-stale-silence-check main.go
```

## Configuration

Example Sensu Go definition:

```json
{
  "type": "CheckConfig",
  "api_version": "core/v2",
  "metadata": {
    "namespace": "default",
    "name": "check_for_stale_silence_entries"
  },
  "spec": {
    "command": "sensu-stale-silence-check -H 127.0.0.1 -u admin -p secret -t 604800",
    "subscriptions": [
      "system"
    ],
    "handlers": [
      "slack"
    ],
    "interval": 3600,
    "publish": true
  }
}
```

## Usage

Help:

```
A Sensu Go check plugin for finding stale silence entries

Usage:
  sensu-stale-silence-check [flags]

Flags:
  -h, --help              help for sensu-stale-silence-check
  -H, --host string       The Sensu API host.
  -p, --password string   A Sensu Go user's password.
  -P, --port string       The port the Sensu API is listening on. (default "8080")
  -t, --threshold int     Threshold in seconds to consider a silenced entry stale (default 604800)
  -T, --timeout int       Time in seconds to consider the API unresponsive (default 10)
  -u, --username string   A Sensu Go user with API access.
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/rgeniesse/sensu-stale-silence-check/releases
