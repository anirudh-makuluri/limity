export interface RateLimitOptions {
  key: string;
  limit?: number;
  window?: number;
}

export interface RateLimitResult {
  allowed: boolean;
  remaining: number;
  reset: number;
}

export interface WindowEntry {
  count: number;
  reset: number;
}
