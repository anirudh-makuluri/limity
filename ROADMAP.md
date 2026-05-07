# Roadmap

## Reliability First (Current)
- Stable release process with changelog and upgrade notes.
- CI required on every PR.
- Concurrency and edge-case tests for core algorithm.
- Benchmark harness and published baseline numbers.

## Next
- Shared-store reference for multi-instance deployments.
- Strategy selection (fixed window, sliding window, token bucket).
- Production-ready examples with retries and `Retry-After` handling.

## Later
- Cross-region reliability tests (clock skew and network jitter).
- Hosted mode SLO targets and incident transparency.
