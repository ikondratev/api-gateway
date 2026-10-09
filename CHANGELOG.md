# Changelog

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
