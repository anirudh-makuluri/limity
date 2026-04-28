# Limity FastAPI Example

FastAPI server with Limity rate limiting.

## Setup

```bash
# Using poetry
poetry install
poetry run python main.py

# Or with pip
pip install -r requirements.txt
python main.py
```

Runs on `http://localhost:8000`

## Endpoints

### GET /health
- Health check (not rate limited)

### GET /api/data
- Limit: 100 requests per 60 seconds per IP
- Returns: data with rate limit info

### GET /api/limited
- Limit: 10 requests per 60 seconds per IP (strict!)
- Returns: limited endpoint warning

### POST /api/create
- Limit: 20 requests per 60 seconds per user/IP
- Returns: created resource

## Testing

### Single request
```bash
curl http://localhost:8000/api/data
```

### Flood requests
```bash
for i in {1..15}; do curl http://localhost:8000/api/data; done
```

### Check headers
```bash
curl -i http://localhost:8000/api/data
```

## Global Rate Limiting

The middleware applies to all routes (except health/docs):
- Limit: 1000 requests per 60 seconds per IP
- Returns 429 if exceeded

## How It Works

1. **Global middleware** checks rate limit
2. **Route handlers** apply endpoint-specific limits
3. **Headers** added to response
4. **Returns 429** if any limit exceeded

## Customization

Edit `main.py`:

```python
# Change global limit
result = rate_limit(
    key=f"global:{client_ip}",
    limit=1000,   # <-- Change this
    window=60,
)
```

## Hosted Mode

Enable by setting environment variable:

```bash
export RATE_LIMIT_API_KEY=your_api_key
python main.py
```

Without it, uses fast in-memory limiter.

## Files

```
main.py           - FastAPI server
pyproject.toml    - Poetry config
requirements.txt  - Dependencies
```
