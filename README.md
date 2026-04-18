# BucketDesk

**BucketDesk** is a lightweight bilingual desktop/local web app for managing MinIO and S3-compatible buckets without exposing the full MinIO Console.

BucketDesk is useful when teams need to upload, browse, copy URLs, and delete objects in a controlled bucket, while administrators keep access limited through scoped S3 credentials and bucket policies.

زبان فارسی: [راهنمای فارسی](./docs/README.fa.md)

## Features

- Bilingual UI: English and Persian.
- RTL/LTR layout switching.
- Multiple MinIO/S3 connection profiles.
- Endpoint, bucket, region, CDN URL, and path-style configuration.
- Connection test with bucket access and write-permission validation.
- Folder-like browsing through S3 prefixes.
- Multi-file upload to the current prefix.
- Select and delete objects.
- Copy public object URLs.
- Local-first app: no external service, database, or telemetry.

## Why BucketDesk?

In many companies, the MinIO Console is not shared with regular users because it exposes too much operational power. BucketDesk gives those users a focused object-management interface while admins can restrict access with S3 policies.

## Recommended Security Model

Do not use MinIO root credentials.

Create a dedicated user/access key for each team or workflow, then scope it to the required bucket and prefixes. See [IAM policy examples](./docs/IAM_POLICIES.md).

## Tech Stack

- **Go** backend for S3/MinIO operations and single-binary distribution.
- **React + TypeScript** frontend for a fast bilingual UI.
- **AWS SDK for Go v2** for S3-compatible APIs.

Go is used because it ships cleanly across macOS, Windows, and Linux without requiring users to install Node.js.

## Development

Requirements:

- Go 1.23+
- Node.js 20+
- npm

Install dependencies:

```bash
npm install
go mod download
```

Run the backend:

```bash
go run ./cmd/bucketdesk
```

Run the frontend during development:

```bash
npm run dev:web
```

Vite proxies `/api` to the Go server.

## Build

Build and embed the web UI, then compile the local app:

```bash
npm run build:app
```

The binary is written to:

```text
dist/bucketdesk
```

Run it:

```bash
./dist/bucketdesk
```

Open the printed local URL, usually:

```text
http://127.0.0.1:5217
```

## GitHub Repository

This project is prepared for:

```text
git@github.com:PouryaMansouri/BucketDesk.git
```

## Roadmap

- Store secrets in OS keychains instead of a local JSON file.
- Signed macOS, Windows, and Linux installers.
- Presigned URL generation.
- Prefix-level guardrails in the UI.
- Optional read-only profiles.
- Object metadata editor.
- Drag-and-drop folder uploads.

## License

BucketDesk is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE).
