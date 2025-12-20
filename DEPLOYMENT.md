# QuickSpin CLI - npm Deployment Guide

Simple automated deployment to npm using GitHub Actions.

## Quick Start

### Installation (for users)
```bash
npm install -g @quickspin/cli
```

### Deployment (for developers)
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

That's it! GitHub Actions builds binaries and publishes to npm automatically.

---

## Setup (One-Time)

### 1. Create npm Organization

```bash
npm login
# Create @quickspin organization at: https://www.npmjs.com/org/create
```

### 2. Add GitHub Secret

Go to: `https://github.com/QuickSpin-Saas/quick-spin-cli/settings/secrets/actions`

Create secret:
- **Name:** `NPM_TOKEN`
- **Value:** Run `npm token create` and paste the token

---

## How It Works

When you push a tag like `v1.0.0`:

### GitHub Actions Automatically:

1. **Builds binaries** for all platforms:
   - macOS (Intel & Apple Silicon)
   - Linux (x86_64, ARM64)
   - Windows (x86_64)

2. **Creates GitHub Release** with:
   - Binary downloads (.tar.gz, .zip)
   - Checksums

3. **Publishes to npm** as `@quickspin/cli`

### Users Install:

```bash
npm install -g @quickspin/cli
```

The npm package:
- Detects user's OS and architecture
- Downloads correct pre-built binary
- Makes `qspin` command available globally

---

## Creating Releases

### Standard Release

```bash
# Commit changes
git add .
git commit -m "Add new feature"
git push

# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Version Scheme

Use semantic versioning: `vMAJOR.MINOR.PATCH`

```bash
# Bug fixes
git tag -a v1.0.1 -m "Fix bug"

# New features
git tag -a v1.1.0 -m "Add feature"

# Breaking changes
git tag -a v2.0.0 -m "Breaking change"
```

---

## Supported Platforms

| OS | Architecture | Binary |
|----|--------------|--------|
| macOS | Intel (x86_64) | ✅ |
| macOS | Apple Silicon (ARM64) | ✅ |
| Linux | x86_64 | ✅ |
| Linux | ARM64 | ✅ |
| Windows | x86_64 | ✅ |

---

## Monitoring

After pushing a tag, monitor:

- **GitHub Actions**: https://github.com/QuickSpin-Saas/quick-spin-cli/actions
- **GitHub Releases**: https://github.com/QuickSpin-Saas/quick-spin-cli/releases
- **npm Package**: https://www.npmjs.com/package/@quickspin/cli

---

## Troubleshooting

### Release Failed?

Check GitHub Actions logs for errors:
- Ensure tests pass locally: `make test`
- Verify tag format: `v1.0.0` (not `1.0.0`)
- Check Go builds locally: `make build`

### npm Publish Failed?

- Verify `NPM_TOKEN` secret is set in GitHub
- Check package name `@quickspin/cli` is available
- Ensure you're logged in: `npm whoami`

### Rollback

```bash
# Delete tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Delete GitHub Release (manually)

# Unpublish from npm (within 72 hours)
npm unpublish @quickspin/cli@1.0.0
```

---

## Testing Locally

```bash
# Build locally
make build

# Test binary
./bin/qspin version

# Build all platforms
make build-all
```

---

## Project Structure

```
quick-spin-cli/
├── .github/workflows/
│   ├── release.yml       # Automated releases
│   └── build.yml         # CI/CD
├── npm/
│   ├── package.json      # npm config
│   ├── install.js        # Binary downloader
│   └── index.js          # Binary wrapper
├── cmd/qspin/            # Go CLI code
├── internal/             # CLI logic
├── Makefile              # Build commands
└── DEPLOYMENT.md         # This file
```

---

## Workflow Details

### Release Workflow (`.github/workflows/release.yml`)

**Triggered by:** Tags matching `v*.*.*`

**Steps:**
1. Build binaries for all platforms
2. Create tar.gz/zip archives
3. Generate checksums
4. Create GitHub Release
5. Publish to npm

### Build Workflow (`.github/workflows/build.yml`)

**Triggered by:** Push to main/develop, Pull requests

**Steps:**
1. Run linter
2. Build all platforms
3. Upload artifacts

---

## Example Workflow

```bash
# Make changes
vim cmd/qspin/main.go

# Test
make test
make build
./bin/qspin version

# Commit
git add .
git commit -m "feat: add new command"
git push

# Release
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# Wait 2-3 minutes for GitHub Actions

# Test installation
npm install -g @quickspin/cli@latest
qspin version
```

---

## Support

- **Issues**: https://github.com/QuickSpin-Saas/quick-spin-cli/issues
- **npm**: https://www.npmjs.com/package/@quickspin/cli

---

**That's all you need!** Just push tags to release. 🚀
