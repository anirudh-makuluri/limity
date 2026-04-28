# Publishing Guide

## npm Publishing

### Prerequisites
1. **npm account** - Create one at https://npmjs.com
2. **Node.js** - v18+ installed
3. **Build packages** - Compile TypeScript to JavaScript

### Build All Packages
```bash
pnpm build
```

This will compile TypeScript to `dist/` in each package.

### Version Management
Update versions in:
- `packages/core/package.json`
- `packages/node/package.json`
- `packages/edge/package.json`

```bash
# For patch version (0.1.0 → 0.1.1)
# Edit package.json or use npm version

pnpm -r --filter="@limity/*" exec npm version patch
```

### Login to npm
```bash
npm login
# Enter username, password, email, OTP
```

### Publish Packages
```bash
# Publish all packages at once
pnpm -r --filter="@limity/*" publish --access public

# Or publish individual packages
cd packages/core && npm publish --access public
cd packages/node && npm publish --access public
cd packages/edge && npm publish --access public
```

**Note:** `--access public` is required for scoped packages (@limity/*)

### Verify Published Packages
```bash
npm search @limity
npm view @limity/core
npm view @limity/node
npm view @limity/edge
```

---

## Python Publishing (PyPI)

### Prerequisites
1. **PyPI account** - Create one at https://pypi.org
2. **Build tools** - `pip install build twine`

### Build Package
```bash
cd packages/python
python -m build
```

Generates:
- `dist/limity-0.1.0.tar.gz` (source)
- `dist/limity-0.1.0-py3-none-any.whl` (wheel)

### Test Package Locally
```bash
cd packages/python
pip install -e .

# Test import
python -c "from limity import rate_limit; print(rate_limit.__doc__)"
```

### Configure twine (One-time)
Create `~/.pypirc`:
```ini
[pypi]
username = __token__
password = pypi-your-token-here
```

Or use environment variable:
```bash
export TWINE_USERNAME=__token__
export TWINE_PASSWORD=pypi-your-token-here
```

### Publish to PyPI
```bash
cd packages/python
python -m twine upload dist/*
```

### Publish to Test PyPI (Optional)
```bash
cd packages/python
python -m twine upload -r testpypi dist/*
# https://test.pypi.org/project/limity/
```

### Verify Published Package
```bash
pip install limity

# Test
python -c "from limity import rate_limit; print('Success!')"
```

---

## Using in Your Projects

### TypeScript/JavaScript Projects

#### Option 1: From npm (Published)
```bash
npm install @limity/core
# or just the middleware you need
npm install @limity/node    # Express
npm install @limity/edge    # Serverless
```

Usage:
```typescript
import { rateLimit } from '@limity/core';
// or
import { rateLimit } from '@limity/node';
```

#### Option 2: From GitHub (Monorepo)
```bash
npm install github:anirudh-makuluri/limity#main
# or specific package
npm install github:anirudh-makuluri/limity/packages/core#main
```

#### Option 3: Locally (Development)
```bash
# In another project
npm install ../path/to/limity/packages/core
```

### Python Projects

#### Option 1: From PyPI (Published)
```bash
pip install limity
```

Usage:
```python
from limity import rate_limit
```

#### Option 2: From GitHub
```bash
pip install git+https://github.com/anirudh-makuluri/limity.git#egg=limity
```

#### Option 3: Locally (Development)
```bash
pip install -e ../path/to/limity/packages/python
```

---

## Versioning Strategy

### Semantic Versioning
Format: `MAJOR.MINOR.PATCH`

Examples:
- `0.1.0` → `0.2.0` (backward-compatible features)
- `0.1.0` → `0.1.1` (bug fixes)
- `0.1.0` → `1.0.0` (breaking changes)

### Publishing Checklist
- [ ] Update version in all package.json files
- [ ] Update CHANGELOG (if you have one)
- [ ] Run `pnpm build` and verify no errors
- [ ] Run examples to test
- [ ] Commit: `git commit -m "chore: bump version to 0.2.0"`
- [ ] Tag: `git tag v0.2.0`
- [ ] Push: `git push && git push --tags`
- [ ] Publish to npm: `pnpm -r --filter="@limity/*" publish`
- [ ] Publish to PyPI: `cd packages/python && twine upload dist/*`

---

## CI/CD (Optional GitHub Actions)

Create `.github/workflows/publish.yml`:

```yaml
name: Publish Packages

on:
  push:
    tags:
      - v*

jobs:
  publish-npm:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v2
      - uses: actions/setup-node@v4
        with:
          node-version: 18
          registry-url: https://registry.npmjs.org
      - run: pnpm install
      - run: pnpm build
      - run: pnpm -r --filter="@limity/*" publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}

  publish-pypi:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v4
      - run: pip install build twine
      - run: cd packages/python && python -m build
      - run: twine upload packages/python/dist/*
        env:
          TWINE_USERNAME: __token__
          TWINE_PASSWORD: ${{ secrets.PYPI_TOKEN }}
```

Then:
1. Generate npm token at https://npmjs.com/settings/tokens
2. Generate PyPI token at https://pypi.org/manage/account/tokens/
3. Add as GitHub secrets: `NPM_TOKEN` and `PYPI_TOKEN`
4. Tag and push: `git tag v0.2.0 && git push --tags`

---

## Troubleshooting

### npm publish errors
- **"You must be logged in"** → `npm login`
- **"Package already published"** → Update version number
- **"Package not found in dist"** → Run `pnpm build` first
- **"Forbidden (403)"** → Check npm account permissions

### PyPI upload errors
- **"Invalid token"** → Check `~/.pypirc` or env variable
- **"Filename already exists"** → Update version in pyproject.toml
- **"Upload rejected"** → Run `twine check dist/*` to validate

### Dependency issues
```bash
# Update @limity/core version in node/edge
cd packages/node
npm install @limity/core@latest
```

---

## Summary

**To publish v0.2.0:**

```bash
# 1. Update versions
nano packages/*/package.json

# 2. Build
pnpm build

# 3. Commit and tag
git add . && git commit -m "chore: release v0.2.0"
git tag v0.2.0
git push && git push --tags

# 4. Publish
pnpm -r --filter="@limity/*" publish --access public
cd packages/python && python -m build && twine upload dist/*

# 5. Verify
npm view @limity/core@0.2.0
pip install limity==0.2.0
```

Done! 🚀
