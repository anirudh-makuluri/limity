import { performance } from 'node:perf_hooks';
import { rateLimit } from '../dist/index.js';

const REQUESTS = 100000;
const WINDOW = 60;
const LIMIT = REQUESTS;

const start = performance.now();
for (let i = 0; i < REQUESTS; i += 1) {
  await rateLimit({ key: `bench:${i % 1000}`, limit: LIMIT, window: WINDOW });
}
const elapsedMs = performance.now() - start;
const rps = Math.round((REQUESTS / elapsedMs) * 1000);

console.log(JSON.stringify({ requests: REQUESTS, elapsedMs: Number(elapsedMs.toFixed(2)), rps }, null, 2));
