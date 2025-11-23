# go-multi-logger-slog

env-based multi logger that:

- if you have grafana loki endpoint, it logs to it and to stdout with json
- if you don't specify grafana loki endpoint, it nicely logs with [lmittmann/tint](github.com/lmittmann/tint)

## Setup

```bash
go get github.com/tomek7667/go-multi-logger-slog@latest
```

In order for loki x json logging you need to specify the following env variables:

- `LOKI_ENDPOINT` - probably something like `https://logs-prod-___.grafana.net/loki/api/v1/push`
- `LOKI_USERNAME` - `10......`
- `LOKI_PASSWORD` - `glc_eyJ....`
- `APP_NAME` - I'd suggest something like `<env>/<kebab-case-app-name>` e.g.: `prod/kup-mi-backend`

You can also specify:

- `LOG_LEVEL` _(default: `DEBUG`)_
- `ENVIRONMENT` _(default: `dev`)_

## Example usage

```go
package main

import (
    "log/slog"

    "github.com/tomek7667/go-multi-logger-slog"
    "github.com/joho/godotenv"
)

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		slog.Info("godotenv load failed", "err", err)
	}

    logger.SetLogLevel()
}

func main() {
    slog.Info("Hello world!")
}
```

## Note

Please note that this package is not optimized, however gets the job done by logging to grafana dashboard. For every log called, it calls a POST to your loki endpoint.
