#!/bin/bash
# Drift Detector Installation Script
# Quick install: curl -sSL https://raw.githubusercontent.com/MeowTux/drift-detector/main/install.sh | bash

set -e

REPO="MeowTux/drift-detector"
BINARY_NAME="drift-detector"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}"
echo "╔═══════════════════════════════════════╗"
echo "║   Drift Detector Installation        ║"
echo "║   Author: MeowTux                    ║"
echo "╚═══════════════════════════════════════╝"
echo -e "${NC}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${YELLOW}Detected OS: $OS${NC}"
echo -e "${YELLOW}Detected Architecture: $ARCH${NC}"
echo ""

# Construct download URL
BINARY_SUFFIX="${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_SUFFIX="${BINARY_SUFFIX}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${BINARY_SUFFIX}"

echo -e "${YELLOW}Downloading from: $DOWNLOAD_URL${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download binary
if command -v curl &> /dev/null; then
    curl -fsSL -o "$BINARY_NAME" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -q -O "$BINARY_NAME" "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: Neither curl nor wget is installed${NC}"
    exit 1
fi

# Make executable
chmod +x "$BINARY_NAME"

# Install binary
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"

if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
fi

# Cleanup
cd -
rm -rf "$TMP_DIR"

# Verify installation
if command -v "$BINARY_NAME" &> /dev/null; then
    echo -e "${GREEN}✓ Installation successful!${NC}"
    echo ""
    "$BINARY_NAME" --version
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo "  1. Initialize configuration: ${YELLOW}drift-detector init${NC}"
    echo "  2. Edit config: ${YELLOW}config/config.yaml${NC}"
    echo "  3. Run detection: ${YELLOW}drift-detector detect${NC}"
    echo ""
    echo "For help: ${YELLOW}drift-detector --help${NC}"
else
    echo -e "${RED}Installation failed. Please check your permissions.${NC}"
    exit 1
fi
