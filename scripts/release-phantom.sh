#!/bin/bash

# Phantom CLI Release Script
# This script builds the phantom CLI and creates a GitHub release

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get version from package.json
VERSION=$(node -p "require('./phantom-cli-npm/package.json').version")
echo -e "${BLUE}🚀 Releasing Phantom CLI v${VERSION}${NC}"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ Not in a git repository${NC}"
    exit 1
fi

# Check if we're on the main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
    echo -e "${YELLOW}⚠️  Warning: Not on main/master branch (currently on: $CURRENT_BRANCH)${NC}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if there are uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}❌ You have uncommitted changes. Please commit or stash them first.${NC}"
    exit 1
fi

# Check if version tag already exists
if git rev-parse "v${VERSION}" >/dev/null 2>&1; then
    echo -e "${RED}❌ Tag v${VERSION} already exists${NC}"
    exit 1
fi

echo -e "${BLUE}📦 Building phantom binaries...${NC}"
make phantom-release

# Check if binaries were built
if [ ! -d "dist/phantom-binaries" ]; then
    echo -e "${RED}❌ Failed to build phantom binaries${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Binaries built successfully${NC}"

# Create a temporary directory for the release
RELEASE_DIR="phantom-release-${VERSION}"
mkdir -p "$RELEASE_DIR"

# Copy binaries to release directory
cp dist/phantom-binaries/* "$RELEASE_DIR/"

# Create a checksums file
echo -e "${BLUE}📝 Creating checksums...${NC}"
cd "$RELEASE_DIR"
for file in *; do
    if [ -f "$file" ]; then
        shasum -a 256 "$file" >> "phantom-${VERSION}-checksums.txt"
    fi
done
cd ..

# Create release notes
cat > "RELEASE_NOTES.md" << EOF
# Phantom CLI v${VERSION}

## Installation

\`\`\`bash
npm install -g @gitvault/phantom
# or
npm install -g phantom-cli
# or
npm install -g phantomkit
\`\`\`

## What's New

- Initial release of Phantom CLI
- Cross-platform support (macOS, Linux, Windows)
- Secure code storage and runtime injection
- Easy project initialization and management

## Supported Platforms

- \`phantom-darwin-amd64\` - macOS Intel
- \`phantom-darwin-arm64\` - macOS Apple Silicon
- \`phantom-linux-amd64\` - Linux Intel
- \`phantom-linux-arm64\` - Linux ARM
- \`phantom-windows-amd64.exe\` - Windows Intel

## Usage

\`\`\`bash
# Initialize a new project
phantom init my-project

# Upload code to GitVault
phantom upload --key YOUR_API_KEY

# Load and execute code at runtime
phantom load my-project --key YOUR_API_KEY
\`\`\`

## Verification

All binaries are signed and can be verified using the checksums file:

\`\`\`bash
shasum -a 256 -c phantom-${VERSION}-checksums.txt
\`\`\`
EOF

echo -e "${BLUE}📝 Release notes created${NC}"

# Create git tag
echo -e "${BLUE}🏷️  Creating git tag v${VERSION}...${NC}"
git tag -a "v${VERSION}" -m "Release Phantom CLI v${VERSION}"

# Push tag to remote
echo -e "${BLUE}📤 Pushing tag to remote...${NC}"
git push origin "v${VERSION}"

# Create GitHub release using gh CLI (if available)
if command -v gh &> /dev/null; then
    echo -e "${BLUE}📦 Creating GitHub release...${NC}"
    gh release create "v${VERSION}" \
        --title "Phantom CLI v${VERSION}" \
        --notes-file "RELEASE_NOTES.md" \
        "$RELEASE_DIR"/*
    
    echo -e "${GREEN}✅ GitHub release created successfully!${NC}"
else
    echo -e "${YELLOW}⚠️  GitHub CLI (gh) not found. Please create the release manually:${NC}"
    echo -e "${YELLOW}   1. Go to https://github.com/gitvault-tech/gitvault/releases${NC}"
    echo -e "${YELLOW}   2. Click 'Create a new release'${NC}"
    echo -e "${YELLOW}   3. Select tag 'v${VERSION}'${NC}"
    echo -e "${YELLOW}   4. Upload files from: $RELEASE_DIR${NC}"
    echo -e "${YELLOW}   5. Use release notes from: RELEASE_NOTES.md${NC}"
fi

# Cleanup
echo -e "${BLUE}🧹 Cleaning up...${NC}"
rm -rf "$RELEASE_DIR"
rm -f "RELEASE_NOTES.md"

echo -e "${GREEN}🎉 Release process completed!${NC}"
echo -e "${BLUE}📋 Next steps:${NC}"
echo -e "${BLUE}   1. Test the npm package: npm install -g @gitvault/phantom${NC}"
echo -e "${BLUE}   2. Verify installation: phantom --version${NC}"
echo -e "${BLUE}   3. Update documentation if needed${NC}"
