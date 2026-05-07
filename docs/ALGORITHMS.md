# Algorithm Notes

## Current Implementation
Limity currently uses a **fixed-window counter** algorithm.

- Key format: `"{key}:{windowStart}"`
- Window unit: seconds
- Allow rule: `count <= limit`
- Reset timestamp: start of current window + window size

## Why Fixed Window
- Very low overhead and simple memory model.
- Easy to reason about request budgeting per time block.
- Works consistently in memory and hosted modes.

## Tradeoffs
- Allows burstiness at window boundaries (end of one window + start of next).
- Accuracy is coarser than sliding-window or token-bucket models.
- In distributed systems, correctness depends on shared state and clock consistency.

## Production Behavior Expectations
- Single process memory mode: suitable for local development or one instance.
- Multi-instance deployments: use hosted/shared backing store for consistent limiting.
- Clock drift: can affect reset timing when multiple machines evaluate windows independently.

## Roadmap Direction
- Add optional sliding-window and token-bucket strategies.
- Add distributed stress tests with injected clock skew.
- Publish benchmark comparisons per algorithm.
