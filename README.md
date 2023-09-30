# loggy

## Motivation

This is a simple webserver that persists your log calls.
It  can be run for example behind a caddy reverse proxy or next to your deployments.
Since it is written in go and uses sqlite3 as a database, it is very easily deployable.

## Todos

- [] Enable json logging

## Installation

Should be as a easy as

```
go install github.com/4thel00z/loggy/...@latest
```

## License

This project is licensed under the GPL-3 license.
