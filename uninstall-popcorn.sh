#!/bin/bash

set -e

echo "🍿 Uninstalling Popcorn..."
echo ""

# Remove binary from /usr/local/pop/bin
if [ -f /usr/local/pop/bin/popcorn ]; then
    echo "🗑️  Removing binary from /usr/local/pop/bin..."
    sudo rm /usr/local/pop/bin/popcorn
    echo "✅ Binary removed"
else
    echo "⚠️  Binary not found in /usr/local/pop/bin"
fi

# Ask if user wants to remove PATH entry from .zshrc
echo ""
read -p "Remove PATH entry from ~/.zshrc? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️  Removing PATH entry from ~/.zshrc..."
    # Remove the popcorn installer lines from .zshrc
    sed -i.bak '/# Added by popcorn installer/,+1d' ~/.zshrc
    echo "✅ PATH entry removed (backup saved as ~/.zshrc.bak)"
    echo "   Run: source ~/.zshrc"
else
    echo "⏭️  Skipping PATH removal"
fi

echo ""
echo "👋 Popcorn uninstalled. See you next time!"