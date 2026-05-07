# Benchmarks

Benchmarks should be treated as directional and run under CI plus local reproducible settings.

## Core Memory Limiter Baseline

Run:

```bash
pnpm --filter @limity/core build
node packages/core/bench/simple-bench.mjs
```

Track and publish:
- Requests/sec
- Median and p95 latency
- Memory usage under sustained load

## Planned Comparisons
- `@limity/core` vs popular Node limiters on identical workloads.
- Burst-at-window-boundary scenarios.
- Shared-store multi-instance workloads.
