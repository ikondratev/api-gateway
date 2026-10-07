# Changelog

## [Unreleased]

## [0.1.0] - 2026-10-07

### Added

- HTTP API gateway service: settings loaded from `settings/{env}_settings.json`, a JSON logger backed by `slog`, and graceful shutdown on `SIGINT` and `SIGTERM`.
- `GET /ping` and request logging middleware that records method, path, status, response size, and duration.
- Container image in `deploy/Dockerfile`.
- Docker Hub publish on a `v*` tag push (`.github/workflows/publish.yml`). The tag must match the `VERSION` file; for this release that is `v0.1.0`.
- Pull request checks: `golangci-lint` and `go test` (`.github/workflows/ci.yml`, `.golangci.yml`).

[Unreleased]: https://github.com/ikondratev/api-gateway/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ikondratev/api-gateway/releases/tag/v0.1.0
