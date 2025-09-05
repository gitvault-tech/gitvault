const os = require('os');

const platform = os.platform();
const arch = os.arch();

const supportedPlatforms = [
  'darwin-x64', 'darwin-arm64',
  'linux-x64', 'linux-arm64',
  'win32-x64'
];

const currentPlatform = `${platform}-${arch}`;

if (!supportedPlatforms.includes(currentPlatform)) {
  console.error(`❌ Unsupported platform: ${currentPlatform}`);
  console.error('Supported platforms:', supportedPlatforms.join(', '));
  console.error('');
  console.error('If you believe this platform should be supported,');
  console.error('please open an issue at: https://github.com/gitvault-tech/gitvault/issues');
  process.exit(1);
}

console.log(`✅ Platform supported: ${currentPlatform}`);
