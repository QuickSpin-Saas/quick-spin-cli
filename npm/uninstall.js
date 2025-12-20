#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

function uninstall() {
  try {
    const binDir = path.join(__dirname, 'bin');
    const distDir = path.join(__dirname, 'dist');

    // Remove bin directory
    if (fs.existsSync(binDir)) {
      fs.rmSync(binDir, { recursive: true, force: true });
      console.log('Removed bin directory');
    }

    // Remove dist directory
    if (fs.existsSync(distDir)) {
      fs.rmSync(distDir, { recursive: true, force: true });
      console.log('Removed dist directory');
    }

    console.log('✅ QuickSpin CLI uninstalled successfully!');
  } catch (error) {
    console.error('⚠️  Uninstallation warning:', error.message);
  }
}

// Only run if called directly
if (require.main === module) {
  uninstall();
}

module.exports = { uninstall };
