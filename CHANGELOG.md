# Changelog

All notable changes to this project are documented in this file.

The format is based on Keep a Changelog and this project follows Semantic Versioning.

## [0.3.0] - 2026-05-04

### Added
- CI pipeline for JavaScript/TypeScript and Python.
- Core test suite with behavioral, boundary, and concurrency coverage.
- Reliability-focused docs: algorithm notes, roadmap, contribution guide.

### Changed
- Added `Retry-After` response header in `@limity/node` when requests are blocked.
- Clarified package boundaries and reliability status in top-level docs.
- Aligned package versions to `0.3.0` to reflect active pre-1.0 progress.

### Notes
- Current algorithm is fixed window (documented in `docs/ALGORITHMS.md`).

## [0.1.3] - 2026-04-27

### Added
- Initial published packages for `@limity/core`, `@limity/node`, `@limity/edge`.
