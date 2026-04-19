# Contributing to BucketDesk

Thanks for helping improve BucketDesk.

## Good First Contributions

- Test BucketDesk on Windows, macOS, or Linux and report packaging issues.
- Improve Persian or English copy.
- Add screenshots or short usage videos.
- Improve MinIO/S3 policy examples.
- Add small UI fixes.
- Add automated tests for storage/profile logic.

## Development Setup

```bash
git clone git@github.com:PouryaMansouri/BucketDesk.git
cd BucketDesk
go test ./...
go run ./cmd/bucketdesk
```

Frontend development is optional:

```bash
npm install
npm run dev:web
```

## Pull Requests

Before opening a PR:

1. Run `go test ./...`.
2. Keep changes focused.
3. Update docs when user behavior changes.
4. Mention the OS you tested on.

## Code Style

- Prefer simple Go code and standard library behavior where possible.
- Keep the app local-first.
- Do not add telemetry or external services.
- Avoid storing secrets in logs.

## Security Issues

Please do not open public issues for vulnerabilities. See [SECURITY.md](./SECURITY.md).
