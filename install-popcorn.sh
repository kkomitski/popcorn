#!/bin/bash

set -e

echo "🍿 Installing Popcorn..."

# Build the binary
echo "🔨 Building popcorn binary..."
make build

# Create /usr/local/popcorn/bin if it doesn't exist
sudo mkdir -p /usr/local/popcorn/bin/

# Move binary to /usr/local/popcorn/bin/
echo "📦 Moving binary to /usr/local/popcorn/bin/..."
sudo mv popcorn /usr/local/popcorn/bin/

# Check if /usr/local/popcorn/bin is in PATH for the current user
if ! echo "$PATH" | grep -q "/usr/local/popcorn/bin"; then
    echo "🔧 Adding /usr/local/popcorn/bin to PATH in ~/.zshrc..."
    echo '' >> ~/.zshrc
    echo '# Added by popcorn installer' >> ~/.zshrc
    echo 'export PATH="$PATH:/usr/local/popcorn/bin"' >> ~/.zshrc
    echo "✨ PATH updated! Please run: source ~/.zshrc"
else
    echo "✅ /usr/local/popcorn/bin already in PATH"
fi

echo ""
echo "╔═══════════════════════════════╗"
echo "║                               ║"
echo "║   🍿 Popcorn Installed! 🎉    ║"
echo "║                               ║"
echo "╚═══════════════════════════════╝"
echo ""
echo "🚀 Run 'popcorn' to start the REPL and enjoy! 🎬"
echo ""
echo "💡 If 'popcorn' command is not found, run:"
echo "   source ~/.zshrc"