#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const https = require('https');
const zlib = require('zlib');
const tar = require('tar');

const GITHUB_REPO = 'QuickSpin-Saas/quick-spin-cli';
const packageJson = require('./package.json');
const VERSION = packageJson.version;

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  const platformMap = {
    darwin: 'darwin',
    linux: 'linux',
    win32: 'windows',
  };

  const archMap = {
    x64: 'x86_64',
    arm64: 'arm64',
  };

  const mappedPlatform = platformMap[platform];
  const mappedArch = archMap[arch];

  if (!mappedPlatform || !mappedArch) {
    throw new Error(
      `Unsupported platform: ${platform}-${arch}. Supported: darwin/linux/windows on x64/arm64`
    );
  }

  return { platform: mappedPlatform, arch: mappedArch };
}

function getBinaryName() {
  return process.platform === 'win32' ? 'qspin.exe' : 'qspin';
}

function getDownloadURL() {
  const { platform, arch } = getPlatform();
  const ext = platform === 'windows' ? 'zip' : 'tar.gz';
  return `https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/qspin-${VERSION}-${platform}-${arch}.${ext}`;
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    console.log(`Downloading ${url}...`);

    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        download(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode} ${response.statusMessage}`));
        return;
      }

      const file = fs.createWriteStream(dest);
      response.pipe(file);

      file.on('finish', () => {
        file.close();
        resolve();
      });

      file.on('error', (err) => {
        fs.unlink(dest, () => {});
        reject(err);
      });
    }).on('error', reject);
  });
}

async function extractArchive(archivePath, destDir) {
  const { platform } = getPlatform();

  if (platform === 'windows') {
    // Extract ZIP
    const AdmZip = require('adm-zip');
    const zip = new AdmZip(archivePath);
    zip.extractAllTo(destDir, true);
  } else {
    // Extract tar.gz
    await tar.x({
      file: archivePath,
      cwd: destDir,
    });
  }
}

async function install() {
  try {
    const binDir = path.join(__dirname, 'bin');
    const distDir = path.join(__dirname, 'dist');

    // Create directories
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }
    if (!fs.existsSync(distDir)) {
      fs.mkdirSync(distDir, { recursive: true });
    }

    const downloadURL = getDownloadURL();
    const { platform } = getPlatform();
    const ext = platform === 'windows' ? 'zip' : 'tar.gz';
    const archivePath = path.join(distDir, `qspin.${ext}`);

    // Download binary
    await download(downloadURL, archivePath);

    // Extract archive
    console.log('Extracting archive...');
    await extractArchive(archivePath, distDir);

    // Move binary to bin directory
    const binaryName = getBinaryName();
    const extractedBinary = path.join(distDir, binaryName);
    const targetBinary = path.join(binDir, binaryName);

    if (fs.existsSync(extractedBinary)) {
      fs.renameSync(extractedBinary, targetBinary);
    } else {
      throw new Error(`Binary not found after extraction: ${extractedBinary}`);
    }

    // Make executable on Unix systems
    if (process.platform !== 'win32') {
      fs.chmodSync(targetBinary, '755');
    }

    // Clean up
    fs.unlinkSync(archivePath);

    console.log('✅ QuickSpin CLI installed successfully!');
    console.log(`Run 'qspin --version' to verify installation.`);
    console.log(`Get started with 'qspin auth login'`);

  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.error('\nPlease try manual installation:');
    console.error(`  Visit: https://github.com/${GITHUB_REPO}/releases/tag/v${VERSION}`);
    process.exit(1);
  }
}

// Only run if called directly
if (require.main === module) {
  install();
}

module.exports = { install };
