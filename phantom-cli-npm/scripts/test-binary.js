const fs = require('fs');
const path = require('path');

console.log('🧪 Testing phantom-cli-npm package...');

// Test 1: Check if all required files exist
const requiredFiles = [
  'index.js',
  'package.json',
  'README.md',
  'scripts/download-binary.js',
  'scripts/check-platform.js'
];

console.log('📁 Checking required files...');
let allFilesExist = true;

for (const file of requiredFiles) {
  const filePath = path.join(__dirname, '..', file);
  if (fs.existsSync(filePath)) {
    console.log(`  ✅ ${file}`);
  } else {
    console.log(`  ❌ ${file} - MISSING`);
    allFilesExist = false;
  }
}

// Test 2: Check if bin directory exists
const binDir = path.join(__dirname, '..', 'bin');
if (fs.existsSync(binDir)) {
  console.log('  ✅ bin/ directory exists');
} else {
  console.log('  ❌ bin/ directory missing');
  allFilesExist = false;
}

// Test 3: Validate package.json structure
console.log('📦 Validating package.json...');
try {
  const packageJson = JSON.parse(fs.readFileSync(path.join(__dirname, '..', 'package.json'), 'utf8'));
  
  const requiredFields = ['name', 'version', 'main', 'bin', 'scripts'];
  for (const field of requiredFields) {
    if (packageJson[field]) {
      console.log(`  ✅ ${field}: ${typeof packageJson[field] === 'object' ? 'object' : packageJson[field]}`);
    } else {
      console.log(`  ❌ ${field} - MISSING`);
      allFilesExist = false;
    }
  }
} catch (err) {
  console.log(`  ❌ Invalid package.json: ${err.message}`);
  allFilesExist = false;
}

// Test 4: Check if index.js is executable
const indexPath = path.join(__dirname, '..', 'index.js');
if (fs.existsSync(indexPath)) {
  const content = fs.readFileSync(indexPath, 'utf8');
  if (content.startsWith('#!/usr/bin/env node')) {
    console.log('  ✅ index.js has proper shebang');
  } else {
    console.log('  ❌ index.js missing shebang');
    allFilesExist = false;
  }
}

if (allFilesExist) {
  console.log('🎉 All tests passed! Package is ready for publishing.');
  process.exit(0);
} else {
  console.log('❌ Some tests failed. Please fix the issues before publishing.');
  process.exit(1);
}
