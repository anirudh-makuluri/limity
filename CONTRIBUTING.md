# Contributing to Limity

Welcome! Limity is a developer-first rate limiting tool built with simplicity and performance in mind.

## Project Structure

```
limity/
├── apps/
│   └── api/                  # Go backend with Redis
├── packages/
│   ├── core/                 # Core rate limiting logic
│   ├── node/                 # Express middleware
│   └── edge/                 # Fetch/Edge helper
├── examples/
│   └── express-app/          # Example Express app
├── pnpm-workspace.yaml       # Monorepo config
└── package.json
```

## Setup

1. Clone the repo:
   ```bash
   git clone https://github.com/limity/limity.git
   cd limity
   ```

2. Install dependencies:
   ```bash
   pnpm install
   ```

3. Set up environment:
   ```bash
   # Backend
   cp apps/api/.env.example apps/api/.env
   
   # Example app
   cp examples/express-app/.env.example examples/express-app/.env
   ```

## Development

### Build all packages
```bash
pnpm build
```

### Watch mode
```bash
pnpm dev
```

### Run tests
```bash
pnpm test
```

### Run example app
```bash
cd examples/express-app
pnpm dev
```

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for design decisions and system overview.

## Code Style

- **TypeScript:** Use strict mode, explicit types
- **Go:** Follow Go conventions, use `gofmt`
- **Minimal:** Less code is better. No overengineering.
- **Clean:** Readable, well-commented code

### TypeScript Guidelines

```typescript
// Good: explicit, clear
async function checkRateLimit(
  key: string,
  limit: number,
  window: number
): Promise<RateLimitResult>

// Bad: ambiguous, implicit
async function check(k, l, w)
```

### Go Guidelines

```go
// Good: clear error handling
count, err := redisIncr(key)
if err != nil {
  return 0, err
}

// Bad: silent failures
count, _ := redisIncr(key)
```

## Making Changes

1. **Create a branch:**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes** following code style above

3. **Test locally:**
   ```bash
   pnpm build
   pnpm test
   ```

4. **Commit with clear message:**
   ```bash
   git commit -m "feat: add new feature"
   ```

5. **Push and create PR:**
   ```bash
   git push origin feature/my-feature
   ```

## PR Guidelines

- Small, focused PRs are better
- Write clear commit messages
- Update documentation if needed
- Test your changes
- Reference issues if applicable

## Areas to Contribute

- **Documentation** - Improve READMEs, examples, guides
- **Tests** - Add unit and integration tests
- **Performance** - Optimize hot paths
- **Features** - New rate limiting algorithms
- **Integrations** - Support for more frameworks

## Reporting Issues

Please include:
- Description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment (OS, Node version, Go version, etc)
- Error logs

## Questions?

- Check existing docs in README.md and ARCHITECTURE.md
- Look at examples in `examples/`
- Review package READMEs

---

Thanks for contributing to Limity! 🚀
