# QuickSpin CLI - Getting Started with npm Deployment

## Quick Summary

Your QuickSpin CLI is set up for automatic npm publishing. Here's all you need to know:

### Installation (for users)
```bash
npm install -g @quickspin/cli
```

### Deployment (for developers)
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# Done! npm package published automatically
```

---

## First-Time Setup

### 1. Create npm Organization

```bash
# Login to npm
npm login

# Create @quickspin organization
# Go to: https://www.npmjs.com/org/create
```

### 2. Add GitHub Secret

Go to: https://github.com/quickspin/quickspin-cli/settings/secrets/actions

**NPM_TOKEN:**
```bash
npm login
npm token create
# Copy token and add to GitHub Secrets as NPM_TOKEN
```

That's it! You're ready to release.

---

## Creating Releases

### Standard Release

```bash
# 1. Commit your changes
git add .
git commit -m "Add awesome feature"
git push

# 2. Create and push version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. Watch automation at:
# https://github.com/quickspin/quickspin-cli/actions
```

### GitHub Actions Will:
1. ✅ Run tests
2. ✅ Build binaries for macOS, Linux, Windows (Intel & ARM)
3. ✅ Create GitHub Release with binaries
4. ✅ Publish to npm as `@quickspin/cli`

### Monitor:
- **GitHub Actions**: https://github.com/quickspin/quickspin-cli/actions
- **GitHub Releases**: https://github.com/quickspin/quickspin-cli/releases
- **npm Package**: https://www.npmjs.com/package/@quickspin/cli

---

## Versioning

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

## What Users Get

When users run `npm install -g @quickspin/cli`, the package:
1. Detects their OS and architecture
2. Downloads the correct pre-built binary
3. Makes `qspin` command available globally

### Supported Platforms:
- macOS (Intel & Apple Silicon)
- Linux (x86_64, ARM64)
- Windows (x86_64)

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
├── .goreleaser.yaml      # Build config
├── Makefile              # Build commands
└── NPM_DEPLOYMENT.md     # Full guide
```

---

## Common Commands

```bash
# Build locally
make build

# Run tests
make test

# Test release (doesn't publish)
goreleaser release --snapshot --clean

# Create release
git tag -a v1.0.0 -m "Release v1.0.0" && git push origin v1.0.0
```

---

## Troubleshooting

### Release failed?
- Check GitHub Actions: https://github.com/quickspin/quickspin-cli/actions
- Ensure tests pass: `make test`
- Verify tag format: `v1.0.0` (not `1.0.0`)

### npm publish failed?
- Check NPM_TOKEN secret is set
- Verify package name `@quickspin/cli` is available
- Check you're logged in: `npm whoami`

### Rollback needed?
```bash
# Delete tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Unpublish from npm (within 72 hours)
npm unpublish @quickspin/cli@1.0.0
```

---

## Full Documentation

For complete details, see: [NPM_DEPLOYMENT.md](NPM_DEPLOYMENT.md)

---

## Example Workflow

```bash
# Make changes to code
vim cmd/qspin/main.go

# Test locally
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
qspin version  # Should show v1.1.0
```

---

## Support

- Issues: https://github.com/quickspin/quickspin-cli/issues
- npm: https://www.npmjs.com/package/@quickspin/cli

---

**That's all you need!** 🚀

The setup is complete. Just push tags to release.
