# QuickSpin CLI - npm Deployment Guide

Simple guide for deploying QuickSpin CLI to npm using GitHub Actions.

## Overview

This setup automatically publishes your CLI to npm whenever you push a version tag. Users install it with:

```bash
npm install -g @quickspin/cli
```

## How It Works

1. **Push a tag** (e.g., `v1.0.0`)
2. **GitHub Actions** runs automatically:
   - Runs tests
   - Builds binaries for all platforms (macOS, Linux, Windows)
   - Creates GitHub Release
   - Publishes to npm
3. **Users install** via `npm install -g @quickspin/cli`
4. **npm downloads** the correct binary for their platform

## Initial Setup

### 1. Create npm Organization

```bash
# Login to npm
npm login

# Create organization (if not exists)
# Go to: https://www.npmjs.com/org/create
# Create organization: @quickspin
```

### 2. Configure GitHub Secrets

Go to: `https://github.com/quickspin/quickspin-cli/settings/secrets/actions`

Add this secret:

**NPM_TOKEN**
```bash
# Create npm token
npm login
npm token create

# Copy the token and add it to GitHub Secrets
```

### 3. Test Locally (Optional)

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Test build
goreleaser release --snapshot --clean

# Check generated files
ls -la dist/
```

## Creating a Release

### Simple Process

```bash
# 1. Ensure all changes are committed
git add .
git commit -m "Add new features"
git push

# 2. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. Done! GitHub Actions handles the rest
```

### What Happens Automatically

1. Tests run
2. Binaries built for:
   - macOS (Intel & Apple Silicon)
   - Linux (x86_64, ARM64)
   - Windows (x86_64)
3. GitHub Release created with binaries
4. npm package published

### Monitor Progress

- GitHub Actions: https://github.com/quickspin/quickspin-cli/actions
- GitHub Releases: https://github.com/quickspin/quickspin-cli/releases
- npm Package: https://www.npmjs.com/package/@quickspin/cli

## Installation for Users

Once published, users can install with:

```bash
# Install globally
npm install -g @quickspin/cli

# Verify installation
qspin version

# Use the CLI
qspin --help
```

## Project Structure

```
quick-spin-cli/
├── .github/
│   └── workflows/
│       ├── release.yml          # Release automation
│       └── build.yml            # CI/CD for PRs
├── npm/
│   ├── package.json             # npm package config
│   ├── install.js               # Downloads correct binary
│   ├── index.js                 # Binary wrapper
│   ├── uninstall.js             # Cleanup
│   └── README.md
├── cmd/
│   └── qspin/main.go            # CLI entry point
├── .goreleaser.yaml             # Build configuration
├── Makefile
├── go.mod
└── go.sum
```

## Versioning

Use semantic versioning: `vMAJOR.MINOR.PATCH`

- **Major** (`v2.0.0`): Breaking changes
- **Minor** (`v1.1.0`): New features, backward compatible
- **Patch** (`v1.0.1`): Bug fixes

Examples:
```bash
# Patch release (bug fixes)
git tag -a v1.0.1 -m "Fix connection bug"
git push origin v1.0.1

# Minor release (new features)
git tag -a v1.1.0 -m "Add new service support"
git push origin v1.1.0

# Major release (breaking changes)
git tag -a v2.0.0 -m "Redesigned CLI interface"
git push origin v2.0.0
```

## Supported Platforms

The npm package automatically installs the correct binary for:

| OS      | Architecture | Binary              |
|---------|--------------|---------------------|
| macOS   | Intel        | qspin-darwin-x86_64 |
| macOS   | Apple Silicon| qspin-darwin-arm64  |
| Linux   | x86_64       | qspin-linux-x86_64  |
| Linux   | ARM64        | qspin-linux-arm64   |
| Windows | x86_64       | qspin-windows-x86_64.exe |

## Troubleshooting

### Release Failed?

1. Check GitHub Actions logs
2. Common issues:
   - Tests failing
   - Invalid tag format (must be `vX.Y.Z`)
   - Missing NPM_TOKEN secret

### npm Publish Failed?

```bash
# Check if package name is available
npm info @quickspin/cli

# Verify you're logged in
npm whoami

# Check token is valid
npm token list
```

### Manual npm Publish

If GitHub Actions fails, you can publish manually:

```bash
# Download binaries from GitHub Release
# Place them in npm/dist/

cd npm
npm version 1.0.0 --no-git-tag-version
npm publish --access public
```

### Rollback a Release

```bash
# Delete tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Delete GitHub Release (manually)

# Unpublish from npm (within 72 hours)
npm unpublish @quickspin/cli@1.0.0
```

## Testing Before Release

```bash
# Test build locally
goreleaser release --snapshot --clean

# Test a binary
./dist/qspin-darwin-arm64/qspin version

# Test npm package locally
cd npm
npm link
qspin version
npm unlink
```

## Release Checklist

Before each release:

- [ ] All tests passing (`make test`)
- [ ] Code linted (`make lint`)
- [ ] CHANGELOG.md updated
- [ ] Version number decided
- [ ] NPM_TOKEN secret is set
- [ ] Tag format is `vX.Y.Z`

After release:

- [ ] GitHub Release created
- [ ] npm package published
- [ ] Test installation: `npm install -g @quickspin/cli`
- [ ] Verify version: `qspin version`

## GitHub Actions Workflows

### Release Workflow (`.github/workflows/release.yml`)

**Triggers:** Tag push matching `v*.*.*`

**Jobs:**
1. **goreleaser** - Build and create GitHub Release
   - Checkout code
   - Run tests
   - Build binaries for all platforms
   - Create GitHub Release with binaries

2. **npm-publish** - Publish to npm
   - Download release binaries
   - Update package.json version
   - Publish to npm registry

### Build Workflow (`.github/workflows/build.yml`)

**Triggers:** Push to main/develop, Pull requests

**Jobs:**
1. **test** - Run tests on Ubuntu, macOS, Windows
2. **lint** - Run linters
3. **build** - Build all platform binaries

## Advanced Configuration

### Custom Build Flags

Edit [.goreleaser.yaml](.goreleaser.yaml):

```yaml
builds:
  - ldflags:
      - -s -w
      - -X main.Version={{.Version}}
      - -X main.YourCustomVar=value
```

### Add More Platforms

```yaml
builds:
  - goos:
      - linux
      - darwin
      - windows
      - freebsd  # Add FreeBSD
    goarch:
      - amd64
      - arm64
      - 386      # Add 32-bit
```

### Pre-release Versions

```bash
# Beta release
git tag -a v1.0.0-beta.1 -m "Beta release"
git push origin v1.0.0-beta.1

# Install beta
npm install -g @quickspin/cli@beta
```

## Support

- **GitHub Issues:** https://github.com/quickspin/quickspin-cli/issues
- **npm Package:** https://www.npmjs.com/package/@quickspin/cli

## Summary

**Deploy to npm in 3 steps:**

1. Set up npm token in GitHub Secrets
2. Push a version tag
3. Done! npm package published automatically

**Users install with:**
```bash
npm install -g @quickspin/cli
```

That's it! 🚀
