#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const binaryPath = path.join(__dirname, 'bin', 'phantom');

// Check if binary exists
if (!fs.existsSync(binaryPath)) {
  console.error('Phantom binary not found. Please run: npm run postinstall');
  process.exit(1);
}

// Spawn the phantom binary with all arguments
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  cwd: process.cwd()
});

child.on('close', (code) => {
  process.exit(code);
});

child.on('error', (err) => {
  console.error('Failed to start phantom:', err.message);
  process.exit(1);
});

