# Changelog

## [0.3.0] - 2026-11-06

### Added
- `POST /api/v1/auth/register`. The gateway calls auth-service `Register` and returns the token pair.
- gRPC contract in `api/auth/v1/auth.proto`. The client reads `upstream.auth_service.url`.
- `X-Trace-Id`. A missing or invalid header is replaced with a UUID. The id is logged as `trace_id` and sent to event-service.

### Changed
- Request logs apply to `/api` only. `GET /ping` stays a health check and is no longer logged.
- Both gRPC clients close when the process exits, including when the HTTP server fails to start.
- Settings are validated on load. `auth_service.url` is required. The write-timeout key is `write_timeout`.

### Notes
- Invalid email and invalid password respond with HTTP 400 (`invalid_email`, `invalid_password`). An existing email responds with HTTP 409 (`email_taken`).
- `Login`, `Refresh`, and `Logout` are not routed yet.

[0.3.0]: https://github.com/ikondratev/api-gateway/releases/tag/v0.3.0

## [0.2.0] - 2026-10-09

### Added

- `POST /api/v1/event` for web and mobile clients.
- gRPC contract in `api/event/v1/event.proto` 
- The gRPC connection is closed after HTTP shutdown on `SIGINT` and `SIGTERM`.
- Client errors are split from internal failures.
### Notes

- The gateway forwards the bearer token and does not verify it. Signature and `events:write` stay on event-service.
- event-service does not serve this gRPC method yet. Until it listens on the configured address, `POST /api/v1/event` responds with HTTP 500.

## [0.1.0] - 2026-10-07

### Added

- HTTP API gateway service: settings loaded from `settings/{env}_settings.json`, a JSON logger backed by `slog`, and graceful shutdown on `SIGINT` and `SIGTERM`.
- `GET /ping` and request logging middleware that records method, path, status, response size, and duration.
- Container image in `deploy/Dockerfile`.
- Docker Hub publish on a `v*` tag push (`.github/workflows/publish.yml`). The tag must match the `VERSION` file; for this release that is `v0.1.0`.
- Pull request checks: `golangci-lint` and `go test` (`.github/workflows/ci.yml`, `.golangci.yml`).

[Unreleased]: https://github.com/ikondratev/api-gateway/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ikondratev/api-gateway/releases/tag/v0.1.0
