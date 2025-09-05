#!/bin/bash

# GitVault PhantomKit Demo Script
# This script demonstrates the basic functionality of the phantom CLI

echo "🔮 GitVault PhantomKit Demo"
echo "=============================="
echo ""

# Check if phantom command exists
if ! command -v phantom &> /dev/null; then
    echo "❌ Phantom CLI not found. Please build and install it first."
    echo "   Run: go build -o phantom cmd/phantom/phantom.go"
    echo ""
    exit 1
fi

echo "✅ Phantom CLI found: $(phantom version --json | jq -r .version)"
echo ""

# Create a demo project
echo "🚀 Creating demo project..."
phantom init demo-project --language js --description "A demo project for PhantomKit" --author "Demo User"
echo ""

# Show the generated configuration
echo "🔧 Generated configuration:"
phantom config show --config demo-project/phantom.config.js
echo ""

# Validate the configuration
echo "✅ Validating configuration:"
phantom config validate --config demo-project/phantom.config.js
echo ""

# Show project structure
echo "📁 Project structure:"
ls -la demo-project/
echo ""

# Show the generated loader
echo "📜 Generated loader file:"
cat demo-project/loader.js
echo ""

# Simulate uploading code
echo "📤 Simulating code upload..."
phantom upload demo-project/loader.js --project demo-project --script main
echo ""

# Simulate loading code
echo "📥 Simulating code load..."
phantom load main --project demo-project --runtime v8
echo ""

# Show version information
echo "ℹ️  Version information:"
phantom version --verbose
echo ""

echo "🎉 Demo completed successfully!"
echo ""
echo "📚 Next steps:"
echo "   1. Explore the generated project structure"
echo "   2. Customize the phantom.config.js file"
echo "   3. Upload your own code files"
echo "   4. Integrate with your development workflow"
echo ""
echo "🔗 Learn more: https://docs.gitvault.io"
echo "💬 Community: https://discord.gg/NsatcWJ"
