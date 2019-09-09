# Sensu Go GoAlert Handler

Sensu Go event handler that will create and close alerts in [GoAlert](https://github.com/target/goalert).

## Installation

Download the latest version of the sensu-go-goalert from releases, or create an executable script from this source.

Run `go get github.com/target/sensu-go-goalert/...`

From the local path of the sensu-go-goalert repository:
`go build ./cmd/sensu-goalert-handler`

## Configuration

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
