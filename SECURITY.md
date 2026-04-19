# Security Policy

BucketDesk is a local-first MinIO and S3-compatible bucket manager.

## Supported Versions

Security fixes target the latest released version.

## Reporting a Vulnerability

Please report vulnerabilities privately by contacting the maintainer through GitHub.

Do not include production credentials, real access keys, or private bucket names in reports.

## Credential Handling

Current behavior:

- Profiles are stored locally in the OS user config directory.
- Profile files are written with restrictive file permissions where supported.
- BucketDesk does not send credentials to external services.

Planned improvement:

- Store secrets in OS-native secure storage such as macOS Keychain, Windows Credential Manager, and libsecret.

## Recommended Usage

- Do not use MinIO root credentials.
- Create dedicated scoped access keys.
- Restrict each key to required buckets and prefixes only.
- Avoid granting `s3:DeleteObject` unless users really need delete access.
