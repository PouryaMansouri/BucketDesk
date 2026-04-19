# Signing and Notarization

BucketDesk can build unsigned installers immediately. For the smoothest user experience, production releases should be signed.

## macOS

Unsigned `.dmg` files may trigger Gatekeeper warnings. To sign and notarize macOS releases, configure these GitHub Secrets:

- `MACOS_CERTIFICATE_P12`: base64-encoded Developer ID Application `.p12`.
- `MACOS_CERTIFICATE_PASSWORD`: password for the `.p12`.
- `MACOS_CODESIGN_IDENTITY`: Developer ID Application identity name.
- `APPLE_ID`: Apple Developer account email.
- `APPLE_TEAM_ID`: Apple Developer Team ID.
- `APPLE_APP_SPECIFIC_PASSWORD`: app-specific password for notarization.

Create a base64 value:

```bash
base64 -i DeveloperIDApplication.p12 | pbcopy
```

The release workflow imports the certificate, signs `BucketDesk.app`, signs the `.dmg`, submits it to Apple notarization, and staples the result when all secrets are present.

## Windows

Unsigned Windows installers may trigger SmartScreen warnings. To sign Windows releases, configure these GitHub Secrets:

- `WINDOWS_CERTIFICATE_PFX`: base64-encoded code signing `.pfx`.
- `WINDOWS_CERTIFICATE_PASSWORD`: password for the `.pfx`.

Create a base64 value on macOS/Linux:

```bash
base64 -i codesign.pfx | pbcopy
```

The release workflow signs both `bucketdesk.exe` and the generated setup installer when these secrets are present.

## Linux

Linux AppImage and `.deb` packages are currently unsigned. A future release can add:

- GPG detached signatures for release artifacts.
- APT repository metadata signing.
