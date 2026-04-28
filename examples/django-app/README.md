# Limity Django Example

Django server with Limity rate limiting.

## Setup

```bash
pip install -r requirements.txt
python manage.py runserver 0.0.0.0:8000
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
- Accepts: `{"userId": "user123"}` (optional)
- Returns: created resource

## Testing

### Single request
```bash
curl http://localhost:8000/api/data
```

### Create resource
```bash
curl -X POST http://localhost:8000/api/create \
  -H "Content-Type: application/json" \
  -d '{"userId":"user123"}'
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

Applied via middleware:
- Limit: 1000 requests per 60 seconds per IP
- Returns 429 if exceeded

## How It Works

1. **Middleware** checks global rate limit
2. **View function** applies endpoint-specific limits
3. **Response headers** added with rate limit info
4. **Returns 429** if any limit exceeded

## Customization

Edit `manage.py`:

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
python manage.py runserver
```

Without it, uses fast in-memory limiter.

## Files

```
manage.py         - Django app and views
requirements.txt  - Dependencies
```

Note: This is a simplified Django setup without a traditional project structure.
For production, create a proper Django project with apps.
