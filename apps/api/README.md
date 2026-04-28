# Limity API Backend

Simple, minimal rate limiting backend using Go and Upstash Redis.

## Setup

1. Create `.env` from `.env.example`
2. Set your Upstash Redis credentials
3. Run: `go run main.go`

## Testing

```bash
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{
    "key": "user:123",
    "limit": 10,
    "window": 60
  }'
```

Response:
```json
{
  "allowed": true,
  "remaining": 9,
  "reset": 1714396860
}
```
