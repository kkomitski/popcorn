#!/bin/bash

set -e

echo "🍿 Installing Popcorn..."

# Build the binary
echo "🔨 Building popcorn binary..."
go build -o pop main.go

# Create ~/bin if it doesn't exist
mkdir -p ~/bin

# Move binary to ~/bin
echo "📦 Moving binary to ~/bin..."
mv pop ~/bin/

# Check if ~/bin is in PATH
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    echo "🔧 Adding ~/bin to PATH in ~/.zshrc..."
    echo '' >> ~/.zshrc
    echo '# Added by popcorn installer' >> ~/.zshrc
    echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
    echo "✨ PATH updated! Please run: source ~/.zshrc"
else
    echo "✅ ~/bin already in PATH"
fi

echo ""
echo "╔═══════════════════════════════╗"
echo "║                               ║"
echo "║   🍿 Popcorn Installed! 🎉    ║"
echo "║                               ║"
echo "╚═══════════════════════════════╝"
echo ""
echo "🚀 Run 'pop' to start the REPL and enjoy! 🎬"
echo ""
echo "💡 If 'pop' command is not found, run:"
echo "   source ~/.zshrc"
