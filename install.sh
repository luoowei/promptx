#!/usr/bin/env bash
set -e

REPO="luoowei/promptx"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

info() { printf "${BLUE}â†?{NC} %s\n" "$1"; }
success() { printf "${GREEN}âœ?{NC} %s\n" "$1"; }
error() { printf "${RED}âœ?{NC} %s\n" "$1" >&2; exit 1; }

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) error "Unsupported architecture: $ARCH" ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) error "Unsupported OS: $OS (Windows: use install.ps1)" ;;
esac

BINARY="px_${OS}_${ARCH}"
if [ "$OS" = "darwin" ]; then BINARY="px_darwin_${ARCH}"; fi

# Determine download URL
if [ "$VERSION" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"
fi

info "Downloading PromptX ${VERSION} for ${OS}/${ARCH}..."
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

if command -v curl &> /dev/null; then
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/px" || error "Download failed. Check your internet connection."
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$TMP_DIR/px" || error "Download failed. Check your internet connection."
else
    error "Neither curl nor wget found. Install one of them first."
fi

chmod +x "$TMP_DIR/px"

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/px" "$INSTALL_DIR/px"
else
    sudo mv "$TMP_DIR/px" "$INSTALL_DIR/px"
fi

success "PromptX installed to $INSTALL_DIR/px"

# Verify installation
if command -v px &> /dev/null; then
    success "Done! Run 'px --help' to get started."
else
    info "Installation complete. Add $INSTALL_DIR to your PATH if needed."
    info "export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo ""
echo "  Next steps:"
echo "  1. Set your API key: export OPENAI_API_KEY=\"sk-...\""
echo "  2. Try it out: px ask \"hello world\""
echo "  3. Interactive mode: px"
echo ""
echo "  â­?Star the repo: https://github.com/${REPO}"
