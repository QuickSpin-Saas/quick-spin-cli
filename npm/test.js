#!/usr/bin/env node

const { execSync } = require('child_process');
const path = require('path');

const binName = process.platform === 'win32' ? 'qspin.exe' : 'qspin';
const binPath = path.join(__dirname, 'bin', binName);

try {
  console.log('Testing qspin installation...');

  const output = execSync(`"${binPath}" version`, {
    encoding: 'utf8',
    stdio: 'pipe',
  });

  console.log(output);
  console.log('✅ Test passed!');
  process.exit(0);
} catch (error) {
  console.error('❌ Test failed:', error.message);
  process.exit(1);
}
