# BucketDesk

![BucketDesk icon](./assets/bucketdesk.png)

**BucketDesk** is a lightweight, bilingual MinIO GUI and S3-compatible bucket manager for teams that need safe object management without exposing the full MinIO Console.

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://go.dev/)
[![Platforms](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-brightgreen.svg)](#downloads-and-installation)

BucketDesk helps users browse buckets, upload files, delete selected objects, and copy public URLs using scoped S3 credentials. Admins can keep MinIO root/admin access private and delegate only the bucket or prefix permissions users actually need.

فارسی: [راهنمای فارسی](./docs/README.fa.md)

## Keywords

MinIO GUI, MinIO client, S3 browser, S3-compatible storage manager, bucket manager, object storage UI, desktop S3 manager, MinIO desktop app, AWS S3 compatible client, bucket upload tool.

## Features

- Bilingual UI: English and Persian.
- Automatic RTL/LTR layout switching.
- Multiple MinIO/S3 connection profiles.
- Endpoint, bucket, region, CDN URL, and path-style configuration.
- Connection test with bucket access and write-permission validation.
- Folder-like browsing through S3 prefixes.
- Multi-file upload to the current prefix.
- Select and delete objects.
- Copy public object URLs.
- Local-first app: no external service, database, or telemetry.
- Portable builds for Windows, macOS, and Linux.
- Installers for Windows, macOS, and Linux.

## Why BucketDesk?

In many companies, the MinIO Console is not shared with regular users because it exposes too much operational power. BucketDesk gives those users a focused object-management interface while admins can restrict access with S3 policies.

Good use cases:

- Marketing or content teams uploading media files.
- Support teams browsing exported files.
- Developers sharing limited S3-compatible storage access.
- Internal tools where users should manage objects but not MinIO itself.

## Downloads and Installation

Go to the [GitHub Releases](https://github.com/PouryaMansouri/BucketDesk/releases) page and download the file for your OS.

| Platform | Recommended download | Portable download |
| --- | --- | --- |
| Windows | `BucketDesk_<version>_windows_amd64_setup.exe` | `bucketdesk_<version>_windows_amd64.zip` |
| macOS Apple Silicon | `BucketDesk_<version>_macos_arm64.dmg` | `bucketdesk_<version>_darwin_arm64.tar.gz` |
| macOS Intel | `BucketDesk_<version>_macos_amd64.dmg` | `bucketdesk_<version>_darwin_amd64.tar.gz` |
| Linux amd64 | `BucketDesk_<version>_linux_x86_64.AppImage` or `.deb` | `bucketdesk_<version>_linux_amd64.tar.gz` |
| Linux arm64 | `.deb` or `.tar.gz` | `bucketdesk_<version>_linux_arm64.tar.gz` |

### Windows

Download the setup `.exe`, run it, then open BucketDesk from the Start Menu or desktop shortcut.

Portable mode: download the Windows `.zip`, extract it, and run `bucketdesk.exe`.

### macOS

Download the `.dmg`, open it, and run `BucketDesk.app`.

If macOS warns that the app is from an unidentified developer, right-click the app and choose **Open**. Signed and notarized releases will remove this warning once project signing credentials are configured.

Portable mode: download the macOS `.tar.gz`, extract it, and run `bucketdesk`.

### Linux

Recommended: download the `.AppImage`, make it executable, and run it:

```bash
chmod +x BucketDesk_v0.1.0_linux_x86_64.AppImage
./BucketDesk_v0.1.0_linux_x86_64.AppImage
```

Debian/Ubuntu users can install the `.deb`:

```bash
sudo dpkg -i bucketdesk_v0.1.0_linux_amd64.deb
```

Portable mode: download the `.tar.gz`, extract it, and run `bucketdesk`.

## How to Use

1. Open BucketDesk.
2. Create a profile.
3. Enter your MinIO/S3 endpoint, access key, secret key, region, and bucket.
4. Keep **Use Path-Style Endpoint** enabled for most MinIO installations.
5. Click **Test connection**.
6. Click **Save**.
7. Browse folders/prefixes, upload files, copy URLs, or delete selected objects.

BucketDesk starts a local server on `127.0.0.1` and opens the UI in your browser. It does not send your credentials to any external service.

## Recommended Security Model

Do not use MinIO root credentials.

Create a dedicated user/access key for each team or workflow, then scope it to the required bucket and prefixes. See [IAM policy examples](./docs/IAM_POLICIES.md).

## Tech Stack

- **Go** backend for S3/MinIO operations and single-binary distribution.
- **React + TypeScript** frontend source for future UI development.
- **Embedded HTML UI** so the app can run without Node.js.
- **AWS SDK for Go v2** for S3-compatible APIs.

Go is used because it ships cleanly across macOS, Windows, and Linux without requiring users to install Node.js.

## Development

Requirements:

- Go 1.23+
- Node.js 20+ and npm, only for frontend development.

Run the backend:

```bash
go run ./cmd/bucketdesk
```

Run the frontend during development:

```bash
npm install
npm run dev:web
```

Vite proxies `/api` to the Go server.

## Build

Compile the local app:

```bash
go build -o dist/bucketdesk ./cmd/bucketdesk
```

Run it:

```bash
./dist/bucketdesk
```

## Distribution

Create all portable archives locally:

```bash
VERSION=v0.1.0 ./scripts/package-portable.sh
```

Create macOS DMG files locally:

```bash
VERSION=v0.1.0 ARCH=arm64 ./scripts/package-macos-dmg.sh
VERSION=v0.1.0 ARCH=amd64 ./scripts/package-macos-dmg.sh
```

Create Linux `.deb` packages on Linux:

```bash
VERSION=v0.1.0 ARCH=amd64 ./scripts/package-linux-deb.sh
VERSION=v0.1.0 ARCH=arm64 ./scripts/package-linux-deb.sh
```

Create a Linux AppImage:

```bash
VERSION=v0.1.0 ARCH=x86_64 ./scripts/package-linux-appimage.sh
```

Windows installer builds are handled in GitHub Actions with Inno Setup.

## Release Automation

Push a version tag to create a GitHub Release with installers and portable archives:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds:

- Windows setup `.exe`
- macOS `.dmg` for Intel and Apple Silicon
- Linux `.deb` for amd64 and arm64
- Linux `.AppImage` for x86_64
- Portable archives for Windows, macOS, and Linux

Signing and notarization are optional and controlled through GitHub Secrets. See [Signing and Notarization](./docs/SIGNING.md).

## Roadmap

BucketDesk is open to contributions. Issues, feature ideas, docs improvements, translations, packaging fixes, and platform testing reports are welcome.

Planned work:

- Store secrets in OS keychains instead of a local JSON file.
- Signed Windows installer.
- Signed and notarized macOS DMG.
- Linux AppImage polish and optional Flatpak package.
- Presigned URL generation.
- Prefix-level guardrails in the UI.
- Optional read-only profiles.
- Object metadata editor.
- Drag-and-drop folder uploads.
- Better error messages for common MinIO policy issues.
- Automated UI tests for release builds.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## Security

See [SECURITY.md](./SECURITY.md).

## License

BucketDesk is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE).
