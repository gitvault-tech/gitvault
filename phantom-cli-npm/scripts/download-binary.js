const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');

const platform = os.platform();
const arch = os.arch();

// Map platform/arch to binary names
const binaryMap = {
  'darwin-x64': 'phantom-darwin-amd64',
  'darwin-arm64': 'phantom-darwin-arm64',
  'linux-x64': 'phantom-linux-amd64',
  'linux-arm64': 'phantom-linux-arm64',
  'win32-x64': 'phantom-windows-amd64.exe'
};

const binaryName = binaryMap[`${platform}-${arch}`];
if (!binaryName) {
  console.error(`Unsupported platform: ${platform}-${arch}`);
  console.error('Supported platforms:', Object.keys(binaryMap).join(', '));
  process.exit(1);
}

// GitHub releases
const version = '1.0.0'; // This should match package.json version
const binaryUrl = `https://github.com/gitvault-tech/gitvault/releases/download/v${version}/${binaryName}`;
const binaryPath = path.join(__dirname, '..', 'bin', 'phantom');

// Ensure bin directory exists
const binDir = path.dirname(binaryPath);
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Verify binary works after download
const verifyBinary = (binaryPath) => {
  return new Promise((resolve) => {
    const { spawn } = require('child_process');
    const child = spawn(binaryPath, ['--version'], { stdio: 'pipe' });
    
    child.on('close', (code) => {
      resolve(code === 0);
    });
    
    child.on('error', () => {
      resolve(false);
    });
    
    // Timeout after 5 seconds
    setTimeout(() => {
      child.kill();
      resolve(false);
    }, 5000);
  });
};

console.log(`Downloading Phantom CLI for ${platform}-${arch}...`);
console.log(`URL: ${binaryUrl}`);

// Download binary
const file = fs.createWriteStream(binaryPath);
const request = https.get(binaryUrl, (response) => {
  if (response.statusCode !== 200) {
    console.error(`Failed to download binary: HTTP ${response.statusCode}`);
    console.error(`URL: ${binaryUrl}`);
    console.error('Make sure the release exists and the binary is available.');
    process.exit(1);
  }
  
  response.pipe(file);
  
  file.on('finish', async () => {
    file.close();
    
    try {
      // Set executable permissions
    fs.chmodSync(binaryPath, '755');
      
      // Verify the binary works
      console.log('Verifying binary...');
      const isValid = await verifyBinary(binaryPath);
      
      if (isValid) {
        console.log('✅ Phantom CLI installed successfully!');
      } else {
        console.error('❌ Downloaded binary is not working properly');
        console.error('This might be due to architecture mismatch or corrupted download.');
        process.exit(1);
      }
    } catch (err) {
      console.error('❌ Failed to set permissions or verify binary:', err.message);
      process.exit(1);
    }
  });
  
  file.on('error', (err) => {
    console.error('❌ Failed to write binary file:', err.message);
    process.exit(1);
  });
});

request.on('error', (err) => {
  console.error('❌ Failed to download Phantom CLI:', err.message);
  console.error('Please check your internet connection and try again.');
  process.exit(1);
});

// Set timeout for the request
request.setTimeout(30000, () => {
  console.error('❌ Download timeout. Please check your internet connection.');
  request.destroy();
  process.exit(1);
});