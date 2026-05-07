# Contributing

Thanks for helping improve Limity.

## Principles
- Reliability over feature count.
- Backward compatibility by default.
- Explicit docs for behavior changes.

## Local Setup
```bash
pnpm install
pnpm build
pnpm test
```

Python package checks:
```bash
cd packages/python
pip install -e .[dev]
pytest --cov=limity
```

## Pull Request Checklist
- Add or update tests for behavior changes.
- Update docs (`README`, package docs, or algorithm notes).
- Add an entry to `CHANGELOG.md` under the next release.
- Keep changes scoped and explain migration impact.

## Release Policy
- Semantic versioning (`MAJOR.MINOR.PATCH`).
- Pre-1.0: breaking changes still bump minor.
- Every release must include upgrade notes.
