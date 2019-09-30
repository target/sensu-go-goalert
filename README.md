# Sensu Go GoAlert Handler
[![Bonsai Asset Badge](https://img.shields.io/badge/Sensu%2Go%2Goalert%2Handler-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/target/sensu-go-goalert) [![TravisCI Build Status](https://travis-ci.org/target/sensu-go-goalert.svg?branch=master)](https://travis-ci.org/target/sensu-go-goalert)

Sensu Go event handler that will create and close alerts in [GoAlert](https://github.com/target/goalert).

## Configuration

### Asset Registration

Assets are the best way to make use of this plugin. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 or later, you can use the following command to add the asset: 

`sensuctl asset add target/sensu-go-goalert:VERSION`

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index](https://bonsai.sensu.io/assets/target/sensu-go-goalert).

### Examples

Example handler config:

```json
{
  "type": "Handler",
  "api_version": "core/v2",
  "metadata": {
    "name": "goalert",
    "namespace": "default"
  },
  "spec": {
    "type": "pipe",
    "command": "sensu-goalert-handler",
    "env_vars": ["GOALERT_URL=ENTER_GENERIC_INTEGRATION_KEY_URL_HERE"],
    "filters": [
      "is_incident",
      "not_silenced"
      ],
    "handlers": [],
    "runtime_assets": [],
    "timeout": 15
  }
}
```

Example check config:

```json
{
  "api_version": "core/v2",
  "type": "CheckConfig",
  "metadata": {
    "namespace": "default",
    "name": "example-health"
  },
  "spec": {
    "command": "curl -s http://localhost:3030/health",
    "subscriptions": ["example"],
    "publish": true,
    "interval": 10,
    "handlers": ["goalert"]
  }
}
```

## Installation

Download the latest version of the sensu-goalert-handler from releases, or create an executable script from this source.

Run  
`go get github.com/target/sensu-go-goalert/...`  
to install the `sensu-goalert-handler` binary into your GOPATH automatically.

Or, run  
`go build ./cmd/sensu-goalert-handler`  
to build the `sensu-goalert-handler` binary into the current directory.
