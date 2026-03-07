#!/usr/bin/env bash
# Jave VS Code Extension Installer for Linux/macOS
# This script creates a symlink to the vscode-jave directory

set -e

echo "Installing Jave VS Code extension..."

# Get script directory (repo root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXTENSION_SOURCE="$SCRIPT_DIR/vscode-jave"
EXTENSIONS_DIR="$HOME/.vscode/extensions"
EXTENSION_TARGET="$EXTENSIONS_DIR/jave-language-0.1.0"

# Verify source exists
if [ ! -d "$EXTENSION_SOURCE" ]; then
    echo "Error: vscode-jave directory not found at $EXTENSION_SOURCE"
    echo "Please run this script from the repository root."
    exit 1
fi

# Create extensions directory if needed
mkdir -p "$EXTENSIONS_DIR"

# Remove existing installation
if [ -e "$EXTENSION_TARGET" ]; then
    echo "Removing existing installation at $EXTENSION_TARGET"
    rm -rf "$EXTENSION_TARGET"
fi

# Create symlink
if ln -s "$EXTENSION_SOURCE" "$EXTENSION_TARGET"; then
    echo "✓ Extension installed successfully (symlinked)"
    echo ""
    echo "  Source: $EXTENSION_SOURCE"
    echo "  Target: $EXTENSION_TARGET"
    echo ""
    echo "Next steps:"
    echo "1. Reload VS Code: Ctrl+Shift+P → 'Developer: Reload Window'"
    echo "2. Open a .jave file to verify syntax highlighting"
    
    # Check if toolchain is installed
    if ! command -v javec &> /dev/null || ! command -v baggage &> /dev/null || ! command -v javevm &> /dev/null; then
        echo ""
        echo "Note: To compile and run Jave programs, install the toolchain:"
        echo "  go install github.com/asciifaceman/jave/cmd/javec@latest"
        echo "  go install github.com/asciifaceman/jave/cmd/baggage@latest"
        echo "  go install github.com/asciifaceman/jave/cmd/javevm@latest"
    fi
else
    echo "✗ Failed to create symlink"
    echo ""
    echo "Alternative: Copy the extension manually:"
    echo "  cp -r \"$EXTENSION_SOURCE\" \"$EXTENSION_TARGET\""
    exit 1
fi
