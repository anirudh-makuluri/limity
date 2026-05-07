import limity


def setup_function() -> None:
    limity._memory_limiter.store.clear()
    limity._hosted_limiter = None


def test_allows_within_limit(monkeypatch):
    monkeypatch.delenv("LIMITY_API_KEY", raising=False)
    monkeypatch.delenv("RATE_LIMIT_API_KEY", raising=False)
    monkeypatch.setattr(limity.time, "time", lambda: 1_700_000_000)

    first = limity.rate_limit("user:1", limit=2, window=60)
    second = limity.rate_limit("user:1", limit=2, window=60)
    third = limity.rate_limit("user:1", limit=2, window=60)

    assert first.allowed is True
    assert first.remaining == 1
    assert second.allowed is True
    assert second.remaining == 0
    assert third.allowed is False
    assert third.remaining == 0


def test_resets_after_window(monkeypatch):
    monkeypatch.delenv("LIMITY_API_KEY", raising=False)
    monkeypatch.delenv("RATE_LIMIT_API_KEY", raising=False)

    monkeypatch.setattr(limity.time, "time", lambda: 1_700_000_000)
    limity.rate_limit("user:2", limit=1, window=60)

    monkeypatch.setattr(limity.time, "time", lambda: 1_700_000_061)
    result = limity.rate_limit("user:2", limit=1, window=60)

    assert result.allowed is True
    assert result.remaining == 0


def test_rate_limit_result_dict():
    result = limity.RateLimitResult(allowed=True, remaining=42, reset=123456)
    assert result.to_dict() == {
        "allowed": True,
        "remaining": 42,
        "reset": 123456,
    }
