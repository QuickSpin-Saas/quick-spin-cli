#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

const binName = process.platform === 'win32' ? 'qspin.exe' : 'qspin';
const binPath = path.join(__dirname, 'bin', binName);

// Pass all arguments to the binary
const child = spawn(binPath, process.argv.slice(2), {
  stdio: 'inherit',
  shell: false,
});

child.on('exit', (code) => {
  process.exit(code || 0);
});

child.on('error', (err) => {
  console.error('Failed to execute qspin:', err.message);
  process.exit(1);
});
