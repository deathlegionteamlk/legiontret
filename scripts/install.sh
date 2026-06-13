#!/bin/bash
# ─── LegionTret Install Script ─────────────────────────────────────
# By Death Legion Team
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/deathlegionteam/legiontret/main/scripts/install.sh | sh
#
# This script downloads and installs LegionTret on macOS, Linux, and WSL.

set -e

BINARY_NAME="legiontret"
INSTALL_DIR="/usr/local/bin"
GITHUB_REPO="deathlegionteam/legiontret"
VERSION="${1:-latest}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Print banner
print_banner() {
    echo ""
    echo -e "${CYAN}  ╔═══════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}  ║                                                       ║${NC}"
    echo -e "${CYAN}  ║   ██╗     ██╗███████╗██████╗  ██████╗ ███████╗       ║${NC}"
    echo -e "${CYAN}  ║   ██║     ██║██╔════╝██╔══██╗██╔═══██╗██╔════╝       ║${NC}"
    echo -e "${CYAN}  ║   ██║     ██║█████╗  ██████╔╝██║   ██║███████╗       ║${NC}"
    echo -e "${CYAN}  ║   ██║     ██║██╔══╝  ██╔══██╗██║   ██║╚════██║       ║${NC}"
    echo -e "${CYAN}  ║   ███████╗██║███████╗██████╔╝╚██████╔╝███████║       ║${NC}"
    echo -e "${CYAN}  ║   ╚══════╝╚═╝╚══════╝╚═════╝  ╚═════╝ ╚══════╝       ║${NC}"
    echo -e "${CYAN}  ║                                                       ║${NC}"
    echo -e "${CYAN}  ║         LegionTret by Death Legion Team               ║${NC}"
    echo -e "${CYAN}  ║         Run LLMs locally. Simple. Fast. Free.         ║${NC}"
    echo -e "${CYAN}  ║                                                       ║${NC}"
    echo -e "${CYAN}  ╚═══════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*) echo "darwin" ;;
        Linux*)  echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       echo "unknown" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)              echo "unknown" ;;
    esac
}

# Check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main installation
main() {
    print_banner

    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo -e "${RED}Error: Unsupported operating system or architecture.${NC}"
        echo "OS: $(uname -s), Architecture: $(uname -m)"
        exit 1
    fi

    echo -e "${BLUE}  Detected: ${OS}/${ARCH}${NC}"

    # Determine download URL
    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/legiontret-${OS}-${ARCH}"
    else
        DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/legiontret-${OS}-${ARCH}"
    fi

    echo -e "${BLUE}  Downloading LegionTret ${VERSION}...${NC}"

    # Download binary
    TMP_DIR=$(mktemp -d)
    TMP_FILE="${TMP_DIR}/legiontret"

    if command_exists curl; then
        curl -fsSL --progress-bar "$DOWNLOAD_URL" -o "$TMP_FILE"
    elif command_exists wget; then
        wget -q --show-progress "$DOWNLOAD_URL" -O "$TMP_FILE"
    else
        echo -e "${RED}Error: Neither curl nor wget is installed.${NC}"
        exit 1
    fi

    # Make executable
    chmod +x "$TMP_FILE"

    # Install
    echo -e "${BLUE}  Installing to ${INSTALL_DIR}...${NC}"

    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo -e "${YELLOW}  sudo required to install to ${INSTALL_DIR}${NC}"
        sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Clean up
    rm -rf "$TMP_DIR"

    # Verify installation
    if command_exists legiontret; then
        INSTALLED_VERSION=$(legiontret version 2>/dev/null || echo "unknown")
        echo ""
        echo -e "${GREEN}  ═══════════════════════════════════════════════════${NC}"
        echo -e "${GREEN}  LegionTret installed successfully!${NC}"
        echo -e "${GREEN}  Version: ${INSTALLED_VERSION}${NC}"
        echo -e "${GREEN}  ═══════════════════════════════════════════════════${NC}"
        echo ""
        echo -e "  ${BOLD}Get started:${NC}"
        echo -e "    ${CYAN}legiontret run gemma3${NC}       # Run Gemma 3"
        echo -e "    ${CYAN}legiontret run llama3${NC}       # Run Llama 3"
        echo -e "    ${CYAN}legiontret pull mistral${NC}     # Download Mistral"
        echo -e "    ${CYAN}legiontret list${NC}             # List models"
        echo ""
    else
        echo -e "${RED}  Installation failed. Please add ${INSTALL_DIR} to your PATH.${NC}"
        exit 1
    fi
}

main "$@"
