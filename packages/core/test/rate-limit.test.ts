import { beforeEach, describe, expect, it, vi } from 'vitest';
import { rateLimit } from '../src/index.js';
import { __resetMemoryStoreForTests } from '../src/memory.js';

describe('rateLimit (fixed window)', () => {
  beforeEach(() => {
    __resetMemoryStoreForTests();
    vi.restoreAllMocks();
  });

  it('allows requests up to the limit and blocks the next one', async () => {
    vi.spyOn(Date, 'now').mockReturnValue(1_700_000_000_000);

    const first = await rateLimit({ key: 'user:1', limit: 2, window: 60 });
    const second = await rateLimit({ key: 'user:1', limit: 2, window: 60 });
    const third = await rateLimit({ key: 'user:1', limit: 2, window: 60 });

    expect(first.allowed).toBe(true);
    expect(first.remaining).toBe(1);

    expect(second.allowed).toBe(true);
    expect(second.remaining).toBe(0);

    expect(third.allowed).toBe(false);
    expect(third.remaining).toBe(0);
  });

  it('resets after the time window advances', async () => {
    vi.spyOn(Date, 'now').mockReturnValue(1_700_000_000_000);
    await rateLimit({ key: 'user:2', limit: 1, window: 60 });

    vi.spyOn(Date, 'now').mockReturnValue(1_700_000_061_000);
    const result = await rateLimit({ key: 'user:2', limit: 1, window: 60 });

    expect(result.allowed).toBe(true);
    expect(result.remaining).toBe(0);
  });

  it('handles bursty concurrency consistently within one process', async () => {
    vi.spyOn(Date, 'now').mockReturnValue(1_700_000_000_000);

    const attempts = 200;
    const limit = 75;

    const results = await Promise.all(
      Array.from({ length: attempts }, () => rateLimit({ key: 'user:3', limit, window: 60 }))
    );

    const allowed = results.filter((r) => r.allowed).length;
    const blocked = results.filter((r) => !r.allowed).length;

    expect(allowed).toBe(limit);
    expect(blocked).toBe(attempts - limit);
  });
});
