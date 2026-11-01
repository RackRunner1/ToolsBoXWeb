# TBXM static web server

A small Go HTTP server for serving a static website. It supports extensionless
URLs for `.html` files, redirects `.html` URLs to clean URLs, hides `.git` and
`.env` paths, and does not list directory contents.

## Run

```sh
TBXM_HOST=127.0.0.1 TBXM_PORT=8080 TBXM_ROOT=/var/www/tbxm go run ./cmd/tbxm
```

Configuration:

| Variable | Default | Description |
| --- | --- | --- |
| `TBXM_HOST` | `0.0.0.0` | Address to bind to |
| `TBXM_PORT` | `8080` | TCP port |
| `TBXM_ROOT` | `/var/www/tbxm` | Existing directory containing the site |

Place an `index.html` in each directory intended to be served. Directories
without one return 404 rather than exposing a file listing.

## Caching and HTTPS

HTML responses use `Cache-Control: no-cache`; successful non-HTML files are
cached for one hour; error responses use `no-store`. Security response headers
are added by the server. TLS and HSTS are intentionally left to a TLS-terminating
reverse proxy, which should also redirect HTTP to HTTPS and set HSTS once HTTPS
is permanently enabled for the domain.

## Logging and privacy

The server does **not** log HTTP requests, URLs, client IP addresses, user
agents, referrers, or other visitor activity. Logs are limited to startup,
shutdown, and operational errors. A reverse proxy or hosting platform may have
its own access logs; disable those separately if visitor-level logging is not
desired.

## Development

```sh
go test ./...
go vet ./...
go build ./cmd/tbxm
```

The CI workflow runs tests and `go vet` for pushes and pull requests. Releases
are built manually through the GitHub Actions workflow dispatch interface.

## Source layout

- `cmd/tbxm`: executable entry point.
- `internal/config`: environment configuration and validation.
- `internal/server`: HTTP listener, timeouts, logging, and graceful shutdown.
- `internal/web`: composition of the website handler.
- `internal/web/urls`: clean URLs and `.html` redirects.
- `internal/web/security`: `.git`/`.env` filtering and security headers.
- `internal/web/files`: static files without directory listings.
- `internal/web/cache`: HTTP cache policy.
- `internal/web/*_test.go`: integration tests grouped by web feature.
