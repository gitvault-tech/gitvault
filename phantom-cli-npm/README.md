# @gitvault/phantom

PhantomKit CLI for secure code storage and runtime injection.

## Installation

Install using any of these package names:

```bash
# Recommended (official package)
npm install -g @gitvault/phantom

# Alternative names
npm install -g phantom-cli
npm install -g phantomkit
```

## Usage

```bash
# Initialize a new project
phantom init my-project

# Upload code to GitVault
phantom upload --key YOUR_API_KEY

# Load and execute code at runtime
phantom load my-project --key YOUR_API_KEY
```

## Getting Started

1. Get your API key from GitVault → User Settings → PhantomKit API
2. Initialize your project: `phantom init my-project`
3. Upload your code: `phantom upload --key YOUR_API_KEY`
4. Use in production: `phantom load my-project --key YOUR_API_KEY`