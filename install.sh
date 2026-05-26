#!/bin/sh
set -eu

REPO="MikeCHOKKI/lfa-cli-ai"
VERSION="${1:-latest}"
BINARY="lfa"

detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    CYGWIN*|MINGW*|MSYS*) echo "windows" ;;
    *)       echo "linux" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)             echo "amd64" ;;
  esac
}

OS=$(detect_os)
ARCH=$(detect_arch)
EXT=""
[ "$OS" = "windows" ] && EXT=".exe"

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')
fi

BINARY_NAME="lfa-${OS}-${ARCH}${EXT}"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/v${VERSION}/$BINARY_NAME"

INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "$INSTALL_DIR"

echo "Downloading LFA CLI v${VERSION} for ${OS}/${ARCH}..."
curl -fsSL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/${BINARY}${EXT}"
chmod +x "${INSTALL_DIR}/${BINARY}${EXT}"

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  case "${SHELL:-}" in
    *zsh) echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "${HOME}/.zshrc"
          echo "Added $INSTALL_DIR to ~/.zshrc" ;;
    *bash) echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "${HOME}/.bashrc"
           echo "Added $INSTALL_DIR to ~/.bashrc" ;;
  esac
fi

echo ""
echo "LFA CLI v${VERSION} installed successfully!"
echo "Run 'lfa setup' to configure your environment."
