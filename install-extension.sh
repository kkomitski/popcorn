#!/bin/bash

set -e

echo "🍿 Installing Popcorn VS Code Extension..."

# Navigate to extension directory
cd "$(dirname "$0")/extension"

# Install dependencies
echo "📦 Installing dependencies..."
npm install

# Compile TypeScript
echo "🔨 Compiling TypeScript..."
npm run compile

# Create symlink in VS Code extensions directory
EXTENSION_DIR="$HOME/.vscode/extensions/popcorn-highlight"

# Remove existing symlink/directory if it exists
if [ -L "$EXTENSION_DIR" ] || [ -d "$EXTENSION_DIR" ]; then
    echo "🗑️  Removing existing extension..."
    rm -rf "$EXTENSION_DIR"
fi

# Create symlink
echo "🔗 Creating symlink to VS Code extensions..."
ln -s "$(pwd)" "$EXTENSION_DIR"

echo "$(pwd)" "$EXTENSION_DIR"

echo "✅ Popcorn extension installed successfully!"
echo ""
echo "📝 Please restart VS Code to activate the extension."
echo "   Then open any .pop file to see syntax highlighting!"
