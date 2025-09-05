#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const packageJsonPath = path.join(__dirname, '..', 'package.json');
const originalPackageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));

const packages = [
  { name: '@gitvault/phantom', description: 'PhantomKit CLI for secure code storage and runtime injection' },
  { name: 'phantom-cli', description: 'PhantomKit CLI for secure code storage and runtime injection' },
  { name: 'phantomkit', description: 'PhantomKit CLI for secure code storage and runtime injection' }
];

console.log('🚀 Publishing all Phantom CLI packages...\n');

for (const pkg of packages) {
  try {
    console.log(`📦 Publishing ${pkg.name}...`);
    
    // Update package.json
    const updatedPackageJson = {
      ...originalPackageJson,
      name: pkg.name,
      description: pkg.description
    };
    
    fs.writeFileSync(packageJsonPath, JSON.stringify(updatedPackageJson, null, 2));
    
    // Publish
    execSync('npm publish', { stdio: 'inherit', cwd: path.join(__dirname, '..') });
    
    console.log(`✅ ${pkg.name} published successfully!\n`);
    
  } catch (error) {
    console.error(`❌ Failed to publish ${pkg.name}:`, error.message);
    console.log('');
  }
}

// Restore original package.json
fs.writeFileSync(packageJsonPath, JSON.stringify(originalPackageJson, null, 2));
console.log('🎉 All packages published!');
console.log('\n📋 Published packages:');
packages.forEach(pkg => console.log(`  - ${pkg.name}`));
console.log('\n🔗 Installation commands:');
packages.forEach(pkg => console.log(`  npm install -g ${pkg.name}`));
